package menus

import (
	"General_Framework_Gin/database/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterMenusRoutes 注册用户管理相关路由
func RegisterMenusRoutes(authGroup *gin.RouterGroup) {
	userGroup := authGroup.Group("/menus")
	{
		userGroup.GET("/", listMenus)
		userGroup.POST("/", createMenus)
		userGroup.PUT("/:id", updateMenus)
		userGroup.DELETE("/:id", deleteMenus)

	}
	menuGroup := authGroup.Group("/menu")
	{
		menuGroup.GET("/", listMenus)
		menuGroup.POST("/", createMenus)
		menuGroup.PUT("/:id", updateMenus)
		menuGroup.DELETE("/:id", deleteMenus)
	}
}

// 示例控制器方法
func listMenus(ctx *gin.Context) {
	// 从请求参数中获取角色
	role := ctx.GetString("role") // 获取 URL 参数 ?role=xxx
	if role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}

	// 根据角色查询菜单
	menus, err := mysql.GetMenusByRole(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询菜单失败", "details": err.Error()})
		return
	}

	// 返回菜单数据
	ctx.JSON(http.StatusOK, gin.H{
		"message": "成功获取菜单",
		"menus":   menus,
	})
}

func createMenus(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "创建菜单成功"})
}

func updateMenus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "更新菜单成功"})
}

func deleteMenus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "删除菜单成功"})
}
