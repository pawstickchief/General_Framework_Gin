package policy

import (
	"General_Framework_Gin/controllers"
	"General_Framework_Gin/database/casbin"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterPolicyRoutes(authGroup *gin.RouterGroup) {
	policyGroup := authGroup.Group("/policies")
	{
		// 添加策略
		policyGroup.POST("/edit", controllers.HandlePolicy)
		policyGroup.POST("/", listPolicy)
	}
	roleGroup := authGroup.Group("/role")
	{
		// 添加策略
		roleGroup.POST("/edit", controllers.HandlePolicy)
		roleGroup.POST("/", listPolicy)
	}
}
func listPolicy(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}

	// 查询该角色的所有策略（文件访问权限）
	policies, err := casbin.Enforcer.GetFilteredPolicy(0, role) // 0 表示查询角色（sub）
	if len(policies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "该角色无文件访问权限", "policies": []string{}})
		return
	}
	if err != nil {
		zap.L().Info("权限查询失败:", zap.Error(err))
	}

	// 格式化输出
	formattedPolicies := []map[string]string{}
	for _, policy := range policies {
		if len(policy) >= 3 { // 确保策略格式正确
			formattedPolicies = append(formattedPolicies, map[string]string{
				"role":       policy[0], // 角色
				"resource":   policy[1], // 访问的资源（文件路径）
				"permission": policy[2], // 访问权限（read/write等）
			})
		}
	}

	// 返回查询结果
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "策略获取成功",
		"policies": formattedPolicies,
	})
}
