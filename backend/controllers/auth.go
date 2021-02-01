package controllers

import (
	"dagger/backend/databases"
	"dagger/backend/models"
	"dagger/backend/runtime"
	"dagger/backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	User *models.User
)

//
// @Summary Login
// @Description Login
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"{}"
// @Router /api/v1/loki/auth/login/ [post]
func Login(c *gin.Context) {
	postDataByte, _ := ioutil.ReadAll(c.Request.Body)
	var postData map[string]interface{}
	err := json.Unmarshal(postDataByte, &postData)
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("%s", err))
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

	username := postData["username"].(string)
	password := postData["password"].(string)

	var user *models.User

	ldapEnabled, _ := runtime.Cfg.Bool("ldap", "enabled")
	if ldapEnabled {
		var result bool
		result, user = utils.LdapCheck(username, password)
		if !result {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "用户名或密码错误"})
			return
		}
	} else {
		var u models.User
		user = &u
		result := databases.DB.Model(&models.User{}).Where("username = ?", username).First(&u)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "不存在的用户"})
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "用户名或密码错误"})
			return
		}
	}

	user.LastLoginAt = time.Now().UTC()
	databases.DB.Save(&user)

	token, err := utils.GenerateToken(user.ID, user.Username, time.Hour*24*7)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "token认证错误"})
		return
	}
	c.AbortWithStatusJSON(200, gin.H{"success": true, "token": token})
	return
}

//
// @Summary Register
// @Description Register
// @Accept  json
// @Produce  json
// @Success 201 {string} string	"{}"
// @Router /api/v1/loki/auth/register/ [post]
func Register(c *gin.Context) {
	postDataByte, _ := ioutil.ReadAll(c.Request.Body)
	var postData map[string]interface{}
	err := json.Unmarshal(postDataByte, &postData)
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("%s", err))
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

	allowSignUp, _ := runtime.Cfg.Bool("users", "allow_sign_up")
	if !allowSignUp {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "不需要进行用户注册，请查看配置文件"})
		return
	}

	ldapEnabled, _ := runtime.Cfg.Bool("ldap", "enabled")
	if ldapEnabled {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "已启用Ldap, 请使用ldap账户登陆，无需注册"})
		return
	}

	username := postData["username"].(string)
	password := postData["password"].(string)
	email := postData["email"].(string)

	var user models.User
	result := databases.DB.Model(&models.User{}).Where("username = ?", username).First(&user)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "用户已经存在"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log4Zap(zap.WarnLevel).Warn(fmt.Sprintf("%s", err))
		c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "请查看服务器日志"})
		return
	}

	user = models.User{
		Username:    username,
		Password:    string(hash),
		Email:       email,
		IsActive:    true,
		IsSuperuser: false,
		CreateAt:    time.Now().UTC(),
		LastLoginAt: time.Now().UTC(),
	}

	databases.DB.Create(&user)

	c.AbortWithStatusJSON(201, gin.H{"success": true})
	return
}

//
// @Summary Get userinfo by jwt token
// @Description Get userinfo by jwt token
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"{}"
// @Router /api/v1/loki/auth/userinfo/ [get]
func GetUserInfo(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(models.User)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "user": map[string]string{"username": user.Username}})
	return
}
