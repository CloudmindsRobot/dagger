package controllers

import (
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LokiUserList(c *gin.Context) {
	pageSizeDefault, _ := runtime.Cfg.GetValue("utils", "page_size")

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", pageSizeDefault))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var total int64
	var users []models.User

	countDB := databases.DB.Order("id desc").Model(&models.LogHistory{})
	dataDB := databases.DB.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize)

	countDB.Count(&total)
	dataDB.Find(&users)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "total": total, "data": users, "page": page, "page_size": pageSize})
	return
}
