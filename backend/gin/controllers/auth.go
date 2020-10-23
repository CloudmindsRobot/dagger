package controllers

import (
	"dagger/backend/gin/databases"
	"dagger/backend/gin/models"
	"dagger/backend/gin/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}

	username := postData["username"].(string)
	password := postData["password"].(string)

	var user models.User
	result := databases.DB.Model(&models.User{}).Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": "用户名或密码错误"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": "用户名或密码错误"})
		return
	} else {
		token, err := utils.GenerateToken(user.ID, user.Username, time.Hour*24)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{"success": false, "message": "token认证错误"})
			return
		}
		c.AbortWithStatusJSON(200, gin.H{"success": true, "token": token})
		return
	}
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
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}

	username := postData["username"].(string)
	password := postData["password"].(string)
	email := postData["email"].(string)

	var user models.User
	result := databases.DB.Model(&models.User{}).Where("username = ?", username).First(&user)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": "	用户已经存在"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}

	user = models.User{
		Username:    username,
		Password:    string(hash),
		Email:       email,
		IsActive:    true,
		IsSuperuser: false,
		CreateAt:    time.Now(),
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
	user := sessions.Default(c).Get("user").(models.User)

	c.AbortWithStatusJSON(200, gin.H{"success": true, "user": map[string]string{"username": user.Username}})
	return
}
