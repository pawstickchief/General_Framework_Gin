package logger

import (
	"General_Framework_Gin/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

var Log *zap.Logger

// InitLogger 初始化日志系统
func InitLogger() {
	if config.AppConfig.Log.Mode == "" {
		log.Fatalf("日志模式未定义")
	}

	mode := config.AppConfig.Log.Mode
	if mode == "dev" {
		initDevLogger()
	} else {
		initBuildLogger()
	}
}

// 初始化开发模式日志（终端输出）
func initDevLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic("无法初始化开发日志: " + err.Error())
	}
	Log = logger
}

// 初始化构建模式日志（文件输出，带日志分割）
func initBuildLogger() {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.AppConfig.Log.Filename,
		MaxSize:    config.AppConfig.Log.MaxSize,    // 每个日志文件最大尺寸
		MaxBackups: config.AppConfig.Log.MaxBackups, // 保留最近日志文件数量
		MaxAge:     config.AppConfig.Log.MaxAge,     // 日志文件最多保存天数
		Compress:   config.AppConfig.Log.Compress,   // 是否压缩日志文件
	})

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	Log = zap.New(core, zap.AddCaller())
}

// CloseLogger 关闭日志文件
func CloseLogger() {
	if Log != nil {
		Log.Sync()
	}
}
