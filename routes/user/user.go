package user

import (
	"General_Framework_Gin/controllers"
	"General_Framework_Gin/database/mysql"
	"General_Framework_Gin/schemas/business"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理相关路由
func RegisterUserRoutes(authGroup *gin.RouterGroup) {
	userGroup := authGroup.Group("/users")
	{
		userGroup.POST("/", listUsers)
		userGroup.PUT("/update", updateUser)
		userGroup.POST("/add", addUser)
		userGroup.POST("/delete", deleteUser)

	}
}

// 示例控制器方法
func listUsers(ctx *gin.Context) {
	role := ctx.GetString("role") // 获取 URL 参数 ?role=xxx
	if role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}

	var params business.PaginationParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		// 使用自定义错误响应
		controllers.ResponseErrorWithMsg(ctx, controllers.CodeInvalidParam, "分页参数验证失败: "+err.Error())
		return
	}

	// 调用数据库方法获取用户列表
	usersList, total, err := mysql.GetUsers(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户列表失败", "details": err.Error()})
		return
	}

	// 返回用户数据和总数
	ctx.JSON(http.StatusOK, gin.H{
		"message": "成功获取用户列表",
		"users":   usersList,
		"total":   total,
		"page":    params.Page,
		"limit":   params.Limit,
	})
}

// 添加新用户
func addUser(ctx *gin.Context) {
	role := ctx.GetString("role") // 获取 URL 参数 ?role=xxx
	if role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}
	var params business.User
	if err := ctx.ShouldBindJSON(&params); err != nil {
		// 使用自定义错误响应
		controllers.ResponseErrorWithMsg(ctx, controllers.CodeInvalidParam, "分页参数验证失败: "+err.Error())
		return
	}
	err := mysql.AddUser(params.Username, params.Password, params.Email, params.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "添加用户失败",
			"users":   params.Username,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "成功添加用户",
		"users":   params.Username,
	})
}

func updateUser(ctx *gin.Context) {
	role := ctx.GetString("role") // 获取 URL 参数 ?role=xxx
	if role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}
	var params business.User
	if err := ctx.ShouldBindJSON(&params); err != nil {
		// 使用自定义错误响应
		controllers.ResponseErrorWithMsg(ctx, controllers.CodeInvalidParam, "参数验证失败: "+err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新用户成功"})
}

func deleteUser(ctx *gin.Context) {
	role := ctx.GetString("role") // 获取 URL 参数 ?role=xxx
	if role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}
	var params business.User
	if err := ctx.ShouldBindJSON(&params); err != nil {
		// 使用自定义错误响应
		controllers.ResponseErrorWithMsg(ctx, controllers.CodeInvalidParam, "参数验证失败: "+err.Error())
		return
	}

	err := mysql.DeleteUserByID(params.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "删除用户失败",
			"users": params.Username})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除用户成功",
		"users": params.Username})
}
