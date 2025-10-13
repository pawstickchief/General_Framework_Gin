package logger

import (
	"General_Framework_Gin/schemas"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func Init(cfg *schemas.LogConfig, mode string) (err error) {
	writeSyncer := GetLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)
	encoder := GetEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}
	var core zapcore.Core
	if mode == "dev" {
		//开发模式,日志输出到终端
		consoleenbcoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleenbcoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg := zap.New(core, zap.AddCaller())
	//替换zap库中全局的logger
	zap.ReplaceGlobals(lg)
	return
}
func GetLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	//日志分割根据日志文件大小
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		//Filename:   viper.GetString("log.filename"),
		//MaxSize:    viper.GetInt("log.max_size"),
		//MaxBackups: viper.GetInt("log.max_backups"),
		//MaxAge:     viper.GetInt("log.max_age"),
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
func GetEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//zap预定义的日志格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}
