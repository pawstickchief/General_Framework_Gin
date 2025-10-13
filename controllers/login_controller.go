package controllers

import (
	"General_Framework_Gin/config"
	"General_Framework_Gin/schemas/business"
	"log"

	"General_Framework_Gin/database/mysql"
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

	// **确保 `rememberMe` 具有默认值**
	rememberMe := false

	if loginReq.ReMemberMe != nil {
		rememberMe = *loginReq.ReMemberMe
	}
	// **根据 `rememberMe` 选择 Token 过期时间**
	expirationTime := time.Now().Add(24 * time.Hour) // 默认 1 天
	if rememberMe {                                  // 记住登录 7 天
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
	}

	// 创建 JWT Token
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

	ip := c.ClientIP()

	// 写入登录日志
	logEntry := business.UserLoginLog{
		UserID:      uint(user.ID),
		Username:    user.Username,
		Role:        user.Role,
		LoginTime:   time.Now(),
		IPAddress:   ip,
		TokenExpire: expirationTime,
		Remark:      "用户登录成功",
	}

	if err := mysql.DB.Create(&logEntry).Error; err != nil {
		log.Printf("写入登录日志失败: %v\n", err)
		// 不影响主流程，可忽略返回
	}
	// **返回 Token 和过期时间**
	ResponseSuccess(c, gin.H{
		"token":     tokenString,
		"expiresIn": int(expirationTime.Sub(time.Now()).Seconds()), // 以秒为单位返回
	})
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
