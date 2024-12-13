package routes

import (
	"General_Framework_Gin/controllers"
	"General_Framework_Gin/middleware"
	"General_Framework_Gin/routes/user"
	"General_Framework_Gin/schemas"
	"General_Framework_Gin/services/base"

	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupRouter 配置 Gin 路由和中间件
func SetupRouter(appConfig *schemas.Config) *gin.Engine {
	if appConfig.Log.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 引擎实例
	r := gin.New()

	// 加载全局中间件
	r.Use(middleware.CORS(appConfig.Server.ClientURL))
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.MaxMultipartMemory = appConfig.FileConfig.MaxFileSize << 20

	// 登录路由
	r.POST("/login", controllers.LoginUserVerif)
	r.POST("/user/update", controllers.UpdateUserPassword)
	// 受保护的业务路由组，需认证
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthRequired())
	authGroup.Use(middleware.CasbinMiddleware())

	// 注册文件管理路由
	registerFileRoutes(authGroup, appConfig)

	// 注册用户管理路由 (示例子路由)
	user.RegisterUserRoutes(authGroup)

	// 404 路由处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404 页面不存在",
		})
	})

	return r
}

// registerFileRoutes 注册文件管理路由
func registerFileRoutes(authGroup *gin.RouterGroup, appConfig *schemas.Config) {
	authGroup.POST("/upload", func(ctx *gin.Context) {
		if err := base.HandleFileUpload(ctx, appConfig.FileConfig.UploadDir); err != nil {
			ctx.String(http.StatusInternalServerError, "上传失败: %v", err.Error())
		}
	})

	authGroup.GET("/download", func(ctx *gin.Context) {
		base.HandleFileDownload(ctx, appConfig.FileConfig.UploadDir)
	})
}
