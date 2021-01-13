package middlewares

import (
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func JWTCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTToken := c.Request.Header.Get("Authorization")
		JWTToken = strings.ReplaceAll(JWTToken, "JWT ", "")
		JWTStruct := strings.Split(JWTToken, ".")
		if len(JWTStruct) != 3 {
			utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("jwt parse error, %s", JWTToken))
			c.AbortWithStatusJSON(400, map[string]interface{}{"msg": "jwt parse error"})
		}
		bytes, err := base64.RawStdEncoding.DecodeString(JWTStruct[1])
		if err != nil {
			utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("base64 decode jwt error, %s", err))
			c.AbortWithStatusJSON(400, map[string]interface{}{"msg": "base64 decode jwt error"})
			return
		} else {
			var userInfo map[string]interface{}
			err := json.Unmarshal(bytes, &userInfo)
			if err != nil {
				utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("json unmarshal jwt info error, %s", err))
				c.AbortWithStatusJSON(400, map[string]interface{}{"msg": "json unmarshal jwt info error"})
				return
			}
			if userInfo["exp"].(float64) > float64(time.Now().Unix()) {
				var user models.User
				result := databases.DB.Where("username = ?", userInfo["username"].(string)).First(&user)
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					c.AbortWithStatusJSON(401, map[string]interface{}{"msg": "no user"})
					return
				} else {
					c.Set("user", user)
				}
				c.Next()
			} else {
				utils.Log4Zap(zap.ErrorLevel).Error(fmt.Sprintf("jwt exprired, %s", JWTToken))
				c.AbortWithStatusJSON(401, map[string]interface{}{"msg": "jwt exprired"})
				return
			}
		}
	}
}
