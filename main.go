package main

import (
	"General_Framework_Gin/config"
	"General_Framework_Gin/database/etcd"
	"General_Framework_Gin/database/mysql"
	"General_Framework_Gin/logger"
	"General_Framework_Gin/routes"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var appConfigPath string
	flag.StringVar(&appConfigPath, "c", "config.yaml", "配置文件路径")
	help := flag.Bool("help", false, "显示帮助信息")
	flag.Parse()
	// 加载配置文件
	if err := config.LoadConfig(appConfigPath); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}
	log.Println("配置文件加载成功")
	if *help {
		fmt.Println("使用示例: ./yourapp -config=./config.yaml")
		os.Exit(0)
	}

	// 初始化日志文件
	logger.InitLogger()
	logger.Log.Info("日志系统初始化成功")

	// 初始化数据库连接
	mysql.Init()
	logger.Log.Info("数据库连接成功")

	// 初始化 Redis 和 ETCD (可选)
	etcd.Init()
	logger.Log.Info("缓存和配置数据库初始化成功")

	// 初始化 Gin 路由
	r := routes.SetupRouter(config.AppConfig)
	logger.Log.Info("Gin 路由初始化成功")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppConfig.Server.RedirectPort), // HTTPS 端口
		Handler: r,
	}

	// 启动 HTTPS 服务
	go func() {
		if err := srv.ListenAndServeTLS("ssl/file/cert.pem", "ssl/file/privkey.pem"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("HTTPS listen: %s\n", zap.Error(err))
		}
	}()

	// 优雅关闭服务器
	gracefulShutdown()
}

// 优雅关闭程序
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Log.Info("正在关闭服务器...")
	mysql.Close()
	etcd.Close()
	logger.Log.Info("数据库连接已关闭")
	logger.Log.Info("服务器已停止")
	time.Sleep(2 * time.Second)
	os.Exit(0)
}
