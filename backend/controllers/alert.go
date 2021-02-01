package controllers

import (
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/runtime"
	"dagger/backend/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//
// @Summary Get loki rule list
// @Description Get loki rule list
// @Accept  json
// @Produce  json
// @Param   page_size path int true "Every page count"
// @Param   page path int true "Page index"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/rule [get]
func LokiRuleList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var total int64
	var rules []models.LogRule

	userI, _ := c.Get("user")
	user := userI.(models.User)
	var logUsers []models.LogUser
	databases.DB.Model(&models.LogUser{}).Preload("LogUserGroup").Preload("LogUserGroup.Groups").Where("user_id = ?", user.ID).Find(&logUsers)
	ruleids := []int{}
	for _, user := range logUsers {
		for _, rule := range user.LogUserGroup.Groups {
			ruleids = append(ruleids, rule.RuleID)
		}
	}

	countDB := databases.DB.Order("id desc").Model(&models.LogRule{}).Where("user_id = ? or id in (?)", user.ID, ruleids)
	dataDB := databases.DB.Preload("Labels").Preload("User").Preload("Groups").Preload("Groups.LogUserGroup").Preload("Groups.LogUserGroup.Users").Preload("Groups.LogUserGroup.Users.User").Order("id desc").Where("user_id = ? or id in (?)", user.ID, ruleids).Offset((page - 1) * pageSize).Limit(pageSize)

	countDB.Count(&total)
	dataDB.Find(&rules)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "total": total, "data": rules, "page": page, "page_size": pageSize})
	return
}

func LokiRuleDelete(c *gin.Context) {

	var ruleOnline models.LogRule
	databases.DB.Preload("User").Model(&models.LogRule{}).Where("id = ?", c.Param("id")).First(&ruleOnline)

	tx := databases.DB.Begin()
	res, err := utils.DeleteRuleGroup(ruleOnline.User.Username, fmt.Sprintf("%d", ruleOnline.ID))
	if !res {
		tx.Rollback()
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	tx.Delete(&models.LogRule{ID: id})
	tx.Commit()

	utils.CacheRule()

	flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
	if flush2Alertmanager {
		err := utils.DynamicAlertmanagerConf()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(204, nil)
	return
}

//
// @Summary Get loki group list
// @Description Get loki group list
// @Accept  json
// @Produce  json
// @Param   page_size path int true "Every page count"
// @Param   page path int true "Page index"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/group [get]
func LokiUserGroupList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var total int64
	var groups []models.LogUserGroup

	countDB := databases.DB.Order("id desc").Model(&models.LogUserGroup{})
	dataDB := databases.DB.Preload("User").Preload("Users").Preload("Users.User").Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize)
	countDB.Count(&total)
	dataDB.Find(&groups)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "total": total, "data": groups, "page": page, "page_size": pageSize})
	return
}

func LokiUserGroupJoin(c *gin.Context) {
	var u models.LogUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	userI, _ := c.Get("user")
	user := userI.(models.User)
	u.UserID = user.ID

	var dbCount int64
	databases.DB.Model(&models.LogUser{}).Where("group_id = ? and user_id = ?", u.LogUserGroupID, u.UserID).Count(&dbCount)
	if dbCount == 0 {
		databases.DB.Create(&u)
	}

	utils.CacheRule()

	flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
	if flush2Alertmanager {
		err := utils.DynamicAlertmanagerConf()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(201, user)
	return
}

func LokiUserGroupLeave(c *gin.Context) {
	var u models.LogUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	userI, _ := c.Get("user")
	user := userI.(models.User)
	u.UserID = user.ID

	databases.DB.Delete(&models.LogUser{}, "group_id = ? and user_id = ?", u.LogUserGroupID, u.UserID)

	utils.CacheRule()

	flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
	if flush2Alertmanager {
		err := utils.DynamicAlertmanagerConf()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(204, nil)
	return
}

//
// @Summary Create loki alert group api
// @Description Create loki alert group api
// @Accept  json
// @Produce  json
// @Param   group_name path string true "The group name"
// @Param   user_id path string true "The create user id"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/group/create [post]
func LokiUserGroupCreate(c *gin.Context) {
	var group models.LogUserGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	userI, _ := c.Get("user")
	user := userI.(models.User)
	group.CreateAt = time.Now().UTC()
	group.UserID = user.ID

	databases.DB.Create(&group)

	c.JSON(201, nil)
	return
}

func LokiUserGroupUpdate(c *gin.Context) {
	var group models.LogUserGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	databases.DB.Save(&group)

	utils.CacheRule()

	flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
	if flush2Alertmanager {
		err := utils.DynamicAlertmanagerConf()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(201, nil)
	return
}

func LokiUserGroupDelete(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	databases.DB.Delete(&models.LogUserGroup{ID: id})

	utils.CacheRule()

	flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
	if flush2Alertmanager {
		err := utils.DynamicAlertmanagerConf()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(204, nil)
	return
}

func LokiEventDetailList(c *gin.Context) {
	var details []models.LogEventDetail

	databases.DB.Order("id desc").Preload("LogRule").Where("event_id = ?", c.Param("id")).Find(&details)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "data": details})
	return
}

//
// @Summary Get loki event list
// @Description Get loki event list
// @Accept  json
// @Produce  json
// @Param   page_size path int true "Every page count"
// @Param   page path int true "Page index"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/event [get]
func LokiEventList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var total int64
	var events []models.LogEvent

	userI, _ := c.Get("user")
	user := userI.(models.User)
	var logUsers []models.LogUser
	databases.DB.Model(&models.LogUser{}).Preload("LogUserGroup").Preload("LogUserGroup.Groups").Where("user_id = ?", user.ID).Find(&logUsers)
	ruleids := []int{}
	for _, user := range logUsers {
		for _, rule := range user.LogUserGroup.Groups {
			ruleids = append(ruleids, rule.RuleID)
		}
	}

	countDB := databases.DB.Order("id desc").Model(&models.LogEvent{}).Where("rule_id in (?) or rule_id in (?)", databases.DB.Table("log_rule").Select("id").Where("user_id = ?", user.ID), ruleids)
	dataDB := databases.DB.Order("id desc").Preload("LogRule").Where("rule_id in (?) or rule_id in (?)", databases.DB.Table("log_rule").Select("id").Where("user_id = ?", user.ID), ruleids).Offset((page - 1) * pageSize).Limit(pageSize)

	status := c.DefaultQuery("status", "")
	if status != "" {
		countDB.Where("status = ?", status)
		dataDB.Where("status = ?", status)
	}

	search := c.DefaultQuery("search", "")
	if search != "" {
		countDB.Joins("inner join log_rule on log_rule.id = log_event.rule_id").Where("log_rule.name like ?", fmt.Sprintf("%%%s%%", search))
		dataDB.Joins("inner join log_rule on log_rule.id = log_event.rule_id").Where("log_rule.name like ?", fmt.Sprintf("%%%s%%", search))
	}

	countDB.Count(&total)
	dataDB.Find(&events)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "total": total, "data": events, "page": page, "page_size": pageSize})
	return
}

func LokiEventArchive(c *gin.Context) {
	var eventids []int
	if err := c.ShouldBindJSON(&eventids); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	databases.DB.Model(&models.LogEvent{}).Where("id in (?)", eventids).Updates(map[string]interface{}{"status": "resolved"})

	c.JSON(201, nil)
	return
}

func LokiEventCollect(c *gin.Context) {
	postDataByte, _ := ioutil.ReadAll(c.Request.Body)
	utils.Log4Zap(zap.InfoLevel).Info(string(postDataByte))
	var alerts []interface{}
	err := json.Unmarshal(postDataByte, &alerts)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	details := []models.LogEventDetail{}
	detailids := []int{}
	rule := models.LogRule{}

	for _, a := range alerts {
		alert := a.(map[string]interface{})

		// 移除超时告警
		endsAt := utils.String2Time(alert["endsAt"].(string), "UTC")
		if time.Now().UTC().Sub(endsAt).Seconds() > 0 && endsAt.Unix() > 0 {
			continue
		}

		// 抛弃不存在告警规则的告警
		if _, ok := alert["annotations"].(map[string]interface{})["key"]; !ok {
			continue
		}

		startsAt := utils.String2Time(alert["startsAt"].(string), "UTC")
		key := alert["annotations"].(map[string]interface{})["key"].(string)
		ruleMapI, _ := databases.GC.Get("rule-map")
		ruleMap := ruleMapI.(map[string]models.LogRule)

		// 抛弃找不到规则的告警
		if ruleMap[key].ID == 0 {
			continue
		}

		if rule.ID == 0 {
			rule = ruleMap[key]
		}

		labels, _ := json.Marshal(alert["labels"])
		detail := models.LogEventDetail{
			Summary:     alert["annotations"].(map[string]interface{})["summary"].(string),
			Description: alert["annotations"].(map[string]interface{})["description"].(string),
			StartsAt:    startsAt,
			RuleID:      &(rule.ID),
			Labels:      string(labels),
			Level:       utils.LogLevel(string(labels)),
		}

		databases.DB.Create(&detail)
		details = append(details, detail)
		detailids = append(detailids, detail.ID)

	}

	// 发送告警
	if len(details) > 0 {
		// 汇总
		event := models.LogEvent{
			CreateAt: time.Now().UTC(),
			Status:   "firing",
			RuleID:   &(rule.ID),
			Count:    len(details),
		}

		databases.DB.Create(&event)
		databases.DB.Model(&models.LogEventDetail{}).Where("id in (?)", detailids).Updates(map[string]interface{}{"event_id": event.ID})

		flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
		if flush2Alertmanager {
			var data interface{}
			json.Unmarshal(postDataByte, &data)
			utils.Push2Alertmanager(data)
		}
	}

	c.JSON(200, nil)
	return
}

func LokiRuleCreate(c *gin.Context) {
	var rule models.LogRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	tx := databases.DB.Begin()
	userI, _ := c.Get("user")
	user := userI.(models.User)
	rule.CreateAt = time.Now().UTC()
	rule.UpdateAt = time.Now().UTC()
	rule.UserID = user.ID
	rule.Summary = rule.Description
	labelstr := utils.StructLables(rule.Labels, rule.Name)
	rule.Level = utils.LogLevel(labelstr)
	rule.Key = utils.Md5(labelstr)

	res := tx.Create(&rule)
	if res.Error != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": res.Error.Error()})
		return
	}

	content, err := utils.GenerateYAML(rule)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "生成yaml文件失败"})
		return
	}

	result, err := utils.CreateOrUpdateRuleGroup(user.Username, content)
	if !result {
		tx.Rollback()
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	tx.Commit()

	utils.CacheRule()

	flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
	if flush2Alertmanager {
		err := utils.DynamicAlertmanagerConf()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(201, nil)
	return
}

func LokiRuleUpdate(c *gin.Context) {
	var rule models.LogRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	rule.UpdateAt = time.Now().UTC()
	rule.Summary = rule.Description
	labelstr := utils.StructLables(rule.Labels, rule.Name)
	rule.Level = utils.LogLevel(labelstr)
	rule.Key = utils.Md5(labelstr)

	tx := databases.DB.Begin()

	var ruleOnline models.LogRule
	databases.DB.Preload("User").Model(&models.LogRule{}).Where("id = ?", rule.ID).First(&ruleOnline)

	res := tx.Save(&rule)
	if res.Error != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": res.Error.Error()})
		return
	}

	content, err := utils.GenerateYAML(rule)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "生成yaml文件失败"})
		return
	}

	result, err := utils.CreateOrUpdateRuleGroup(ruleOnline.User.Username, content)
	if !result {
		tx.Rollback()
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	tx.Commit()

	utils.CacheRule()

	flush2Alertmanager, _ := runtime.Cfg.Bool("alertmanager", "enabled")
	if flush2Alertmanager {
		err := utils.DynamicAlertmanagerConf()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(201, gin.H{"success": true, "data": rule})
	return
}

func LokiRuleDownload(c *gin.Context) {

	dir, _ := os.Getwd()
	exportDir := fmt.Sprintf("%s/static/rules", dir)
	cmd := fmt.Sprintf("mkdir -p %s", exportDir)
	_, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("mkdir error, %s", err))
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": "创建规则下载目录失败"})
		return
	}

	userI, _ := c.Get("user")
	user := userI.(models.User)

	filename := fmt.Sprintf("%s.yaml", user.Username)
	absolutePath := fmt.Sprintf("%s/static/rules/%s", dir, filename)
	file, err := os.Create(absolutePath)
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("open loki rule file error, %s", err))
	}
	defer file.Close()

	file.WriteString("\xEF\xBB\xBF")

	content := utils.LoadRules(user.Username)
	if content == "" {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "无所属的规则文件"})
		return
	}
	file.WriteString(content)

	c.JSON(200, gin.H{"success": true, "download": filename})
}
