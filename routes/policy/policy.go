package policy

import (
	"General_Framework_Gin/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterPolicyRoutes(authGroup *gin.RouterGroup) {
	policyGroup := authGroup.Group("/policies")
	{
		// 添加策略
		policyGroup.POST("/", controllers.HandlePolicy)
	}
}
