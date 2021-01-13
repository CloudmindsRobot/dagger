package controllers

import (
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//
// @Summary Create loki query history api
// @Description Create loki query history api
// @Accept  json
// @Produce  json
// @Param   label_json path string true "The query label value dict"
// @Param   filter_json path string true "The query filter value list"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/history/create [post]
func LokiHistoryCreate(c *gin.Context) {
	postDataByte, _ := ioutil.ReadAll(c.Request.Body)
	var postData map[string]interface{}
	err := json.Unmarshal(postDataByte, &postData)
	if err != nil {
		utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("%s", err))
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

	labelJSON := postData["label_json"].(string)
	filterJSON := postData["filter_json"].(string)

	userI, _ := c.Get("user")
	user := userI.(models.User)

	history := models.LogHistory{
		LabelJSON:  labelJSON,
		FilterJSON: filterJSON,
		CreateAt:   time.Now(),
		UserID:     user.ID,
	}
	databases.DB.Save(&history)

	c.JSON(201, nil)
	return
}

//
// @Summary Get loki query history labels
// @Description Get loki query history labels
// @Accept  json
// @Produce  json
// @Param   page_size path int true "Every page count"
// @Param   page path int true "Page index"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/history [get]
func LokiHistoryList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var total int64
	var historys []models.LogHistory

	userI, _ := c.Get("user")
	user := userI.(models.User)
	countDB := databases.DB.Order("id desc").Model(&models.LogHistory{}).Where("user_id = ?", user.ID)
	dataDB := databases.DB.Order("id desc").Offset((page-1)*pageSize).Limit(pageSize).Preload("User").Where("user_id = ?", user.ID)

	countDB.Count(&total)
	dataDB.Find(&historys)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "total": total, "data": historys, "page": page, "page_size": pageSize})
	return
}

//
// @Summary Delete loki query history labels
// @Description Delete loki query history labels
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/history/delete/:id [delete]
func LokiHistoryDelete(c *gin.Context) {
	historyID, _ := strconv.Atoi(c.Param("id"))

	databases.DB.Delete(&models.LogHistory{}, historyID)

	c.JSON(204, nil)
	return
}
