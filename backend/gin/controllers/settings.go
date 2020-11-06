package controllers

import (
	"dagger/backend/gin/runtime"
	"fmt"

	"github.com/gin-gonic/gin"
)

//
// @Summary LoadSettings
// @Description Load Settings for UI
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"{}"
// @Router /api/v1/loki/settings/load/ [get]
func LoadSettings(c *gin.Context) {
	allowSignUp, err := runtime.Cfg.Bool("users", "allow_sign_up")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": fmt.Sprintf("配置加载失败：%s", err.Error())})
		return
	}

	c.AbortWithStatusJSON(200, gin.H{"success": true, "data": map[string]interface{}{"allowSignUp": allowSignUp}})
	return
}
