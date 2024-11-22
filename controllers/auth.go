package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yuqueppbackend/dao"
	"yuqueppbackend/models"
)

var userDao = dao.NewUserDAO()

func LoginPage(c *gin.Context) {
	id, b64s, err := GenerateCaptcha()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"captchaId": id, "captcha": b64s})
}

// 用户注册
func Register(c *gin.Context) {
	var registerData struct {
		Email     string `json:"email" binding:"required"`
		Nickname  string `json:"nickname" binding:"required"`
		Password  string `json:"password" binding:"required"`
		CaptchaId string `json:"captchaId" binding:"required"`
		Value     string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		Email:    registerData.Email,
		Nickname: registerData.Nickname,
		Password: registerData.Password,
	}
	// 注册用户
	if err := userDao.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 用户登录
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(user); err != nil {
	}
	var loginData struct {
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
		CaptchaId string `json:"captchaId" binding:"required"`
		Value     string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证验证码
	if !store.Verify(loginData.CaptchaId, loginData.Value, true) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid captcha"})
		return
	}

	// 验证用户名和密码
	user = models.User{Email: loginData.Email, Password: loginData.Password}
	pass, err := userDao.CheckPassword(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		if pass {
			c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		}
	}
}
