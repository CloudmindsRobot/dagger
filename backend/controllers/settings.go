package controllers

import (
	"dagger/backend/runtime"
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
	allowSignUp, errSign := runtime.Cfg.Bool("users", "allow_sign_up")
	alertEnabled, errAlert := runtime.Cfg.Bool("global", "alert_enabled")
	if errSign != nil || errAlert != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": fmt.Sprintf("配置加载失败")})
		return
	}

	c.AbortWithStatusJSON(200, gin.H{"success": true, "data": map[string]interface{}{"allowSignUp": allowSignUp, "alertEnabled": alertEnabled}})
	return
}
