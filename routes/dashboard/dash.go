package dashboard

import (
	"General_Framework_Gin/database/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterDashRoutes(authGroup *gin.RouterGroup) {
	userGroup := authGroup.Group("/dashboard")
	{
		userGroup.GET("/", listCount)

	}
	dashGroup := authGroup.Group("/dashboards")
	{
		dashGroup.GET("/", listCount)

	}
}

// 示例控制器方法
func listCount(ctx *gin.Context) {
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
