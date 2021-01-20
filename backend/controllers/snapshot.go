package controllers

import (
	"bufio"
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//
// @Summary Get loki query result snapshot list
// @Description Get loki query result snapshot list
// @Accept  json
// @Produce  json
// @Param   page_size path int true "Every page count"
// @Param   page path int true "Page index"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/snapshot [get]
func LokiSnapshotList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var total int64
	var snapshots []models.LogSnapshot

	userI, _ := c.Get("user")
	user := userI.(models.User)
	countDB := databases.DB.Order("id desc").Model(&models.LogSnapshot{}).Where("user_id = ?", user.ID)
	dataDB := databases.DB.Order("id desc").Offset((page-1)*pageSize).Limit(pageSize).Preload("User").Where("user_id = ?", user.ID)

	countDB.Count(&total)
	dataDB.Find(&snapshots)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "total": total, "data": snapshots, "page": page, "page_size": pageSize})
	return
}

//
// @Summary Delete loki query result snapshot file
// @Description Delete loki query result snapshot file
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/snapshot/delete/:id [delete]
func LokiSnapshotDelete(c *gin.Context) {
	snapshotID, _ := strconv.Atoi(c.Param("id"))

	var snapshot models.LogSnapshot
	result := databases.DB.Model(&models.LogSnapshot{}).Where("id = ?", snapshotID).First(&snapshot)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "快照文件不存在"})
		return
	}

	dir, _ := os.Getwd()
	filepath := fmt.Sprintf("%s/%s", dir, snapshot.Dir)
	if utils.FileExists(filepath) {
		cmd := fmt.Sprintf("rm -rf %s", filepath)
		_, err := exec.Command("bash", "-c", cmd).CombinedOutput()
		if err != nil {
			utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("rm file error, %s", err))
			c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
			return
		}
	}

	databases.DB.Delete(&models.LogSnapshot{}, snapshotID)

	c.JSON(204, nil)
	return
}

//
// @Summary Create loki query result snapshot result
// @Description Create loki query result snapshot result
// @Accept  json
// @Produce  json
// @Param   name path string true "Snapshot filename"
// @Param   tmp_file path string true "Snapshot result temp file"
// @Param   start_time path string true "Snapshot query result start time"
// @Param   end_time path string true "Snapshot query result end time"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/snapshot/create [post]
func LokiSnapshotCreate(c *gin.Context) {
	postDataByte, _ := ioutil.ReadAll(c.Request.Body)
	var postData map[string]interface{}
	err := json.Unmarshal(postDataByte, &postData)
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("%s", err))
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

	dir, _ := os.Getwd()

	snapshotDir := fmt.Sprintf("%s/static/snapshot/%s", dir, time.Now().UTC().Format("20060102"))
	cmd := fmt.Sprintf("mkdir -p %s", snapshotDir)
	_, err = exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("mkdir error, %s", err))
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "创建快照文件目录失败"})
		return
	}

	tmpfilename := postData["tmp_filename"].(string)
	filename := postData["name"].(string)
	absolutePath := fmt.Sprintf("%s/static/export/%s", dir, tmpfilename)

	if utils.FileExists(absolutePath) {
		targetPath := fmt.Sprintf("%s/%s", snapshotDir, filename)
		if !utils.FileExists(targetPath) {
			cmd = fmt.Sprintf("mv %s %s", absolutePath, targetPath)
			_, err = exec.Command("bash", "-c", cmd).CombinedOutput()
			if err != nil {
				utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("mv file error, %s", err))
				c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "重命名结果临时文件失败"})
				return
			}
		} else {
			utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("mv file error, %s", err))
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "快照文件已存在"})
			return
		}

		cmd = fmt.Sprintf("wc -l %s", targetPath)
		out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
		if err != nil {
			utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("mv file error, %s", err))
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "统计快照文件内容失败"})
			return
		}
		reg, _ := regexp.Compile(`(\d+)`)
		count, _ := strconv.Atoi(reg.FindString(string(out)))

		t, _ := time.LoadLocation("Local")
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", postData["start_time"].(string), t)
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", postData["end_time"].(string), t)

		userI, _ := c.Get("user")
		user := userI.(models.User)
		snapshot := models.LogSnapshot{
			Name:        filename,
			CreateAt:    time.Now().UTC(),
			StartTime:   startTime.Add(time.Hour * -8),
			EndTime:     endTime.Add(time.Hour * -8),
			DownloadURL: fmt.Sprintf("/api/v1/loki/static/snapshot/%s", fmt.Sprintf("%s/%s", time.Now().UTC().Format("20060102"), filename)),
			Count:       count,
			User:        user,
			Dir:         fmt.Sprintf("static/snapshot/%s/%s", time.Now().UTC().Format("20060102"), filename),
		}

		databases.DB.Create(&snapshot)

		c.JSON(201, nil)
		return
	}

	c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "结果临时文件无法找到"})
	return
}

//
// @Summary Get loki query result snapshot detail
// @Description Get loki query result snapshot detail
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/snapshot/detail/:id [get]
func LokiSnapshotDetail(c *gin.Context) {
	snapshotID, _ := strconv.Atoi(c.Param("id"))

	var snapshot models.LogSnapshot
	result := databases.DB.Model(&models.LogSnapshot{}).Where("id = ?", snapshotID).First(&snapshot)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "快照文件不存在"})
		return
	}
	dir, _ := os.Getwd()
	filepath := fmt.Sprintf("%s/%s", dir, snapshot.Dir)

	if !utils.FileExists(filepath) {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "快照文件不存在"})
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("read file error, %s", err))
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "读取快照文件失败"})
		return
	}
	defer file.Close()

	queryResults := []interface{}{}
	reader := bufio.NewReader(file)
	index := 0
	for {
		index++
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		item := make(map[string]interface{})
		logLevel := utils.LogLevel(string(line))
		item["level"] = logLevel
		item["id"] = index
		item["message"] = string(line)
		queryResults = append(queryResults, item)
	}

	c.AbortWithStatusJSON(200, gin.H{"success": true, "data": queryResults})
}
