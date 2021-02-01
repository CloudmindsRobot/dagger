package controllers

import (
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

//
// @Summary Get loki user list
// @Description Get loki user list
// @Accept  json
// @Produce  json
// @Param   page_size path int true "Every page count"
// @Param   page path int true "Page index"
// @Success 200 {string} string	"[]"
// @Router /api/v1/loki/user [get]
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
