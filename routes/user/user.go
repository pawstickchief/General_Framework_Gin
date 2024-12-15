package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理相关路由
func RegisterUserRoutes(authGroup *gin.RouterGroup) {
	userGroup := authGroup.Group("/users")
	{
		userGroup.GET("/", listUsers)
		userGroup.POST("/", createUser)
		userGroup.PUT("/:id", updateUser)
		userGroup.DELETE("/:id", deleteUser)

	}
}

// 示例控制器方法
func listUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "列出所有用户"})
}

func createUser(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "创建用户成功"})
}

func updateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "更新用户成功"})
}

func deleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "删除用户成功"})
}
