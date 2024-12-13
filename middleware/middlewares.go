package middleware

import (
	"General_Framework_Gin/config"
	"General_Framework_Gin/database/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

// CORS 中间件，支持跨域
func CORS(clientURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", clientURL)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Logger 中间件，记录请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next() // 处理请求

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		log.Printf("%s | %3d | %13v | %15s | %s",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			c.ClientIP(),
			c.Request.Method+" "+c.Request.RequestURI,
		)
	}
}

// AuthRequired 中间件验证
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		// 去掉 "Bearer " 前缀
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 解析令牌
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		// 检查解析结果
		if err != nil || !token.Valid {
			log.Printf("令牌验证失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效令牌"})
			c.Abort()
			return
		}

		// 提取用户名与角色信息
		username := (*claims)["username"].(string)
		role := (*claims)["role"].(string)

		// 将信息存储到上下文中
		c.Set("username", username)
		c.Set("role", role)

		// 权限检查示例
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// Recovery 中间件，捕获异常，防止崩溃
func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter)
}
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetString("role") // 用户角色
		obj := c.Request.URL.Path  // 请求路径
		act := c.Request.Method    // 请求方法

		allowed, err := casbin.Enforcer.Enforce(sub, obj, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限验证失败"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
			c.Abort()
			return
		}

		c.Next()
	}
}
