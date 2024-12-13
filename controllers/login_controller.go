package controllers

import (
	"General_Framework_Gin/config"

	"General_Framework_Gin/database/mysql"
	"General_Framework_Gin/schemas/business"
	"General_Framework_Gin/schemas/request"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// LoginUserVerif 用户登录验证
func LoginUserVerif(c *gin.Context) {
	var loginReq request.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 查询数据库中的用户信息
	user, err := mysql.GetUserByUsername(&loginReq.Username)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if user.Username == "" {
		ResponseError(c, CodeUserNotExist)
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 创建 JWT Token
	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "生成令牌失败")
		return
	}

	// 设置角色到上下文
	c.Set("role", user.Role)

	// 返回成功响应
	ResponseSuccess(c, business.Token{Token: tokenString})
}
func UpdateUserPassword(c *gin.Context) {
	var loginReq request.LoginUpdatePasswordRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 查询数据库中的用户信息
	user, err := mysql.UpdateUserByUsernameAndEmail(loginReq.Username, loginReq.Email, loginReq.Password)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if user != true {
		ResponseError(c, CodeUserNotExist)
		return
	}

	// 返回成功响应
	ResponseSuccess(c, user)
}
