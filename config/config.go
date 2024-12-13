package config

import (
	"log"
	"sync"

	"General_Framework_Gin/schemas"
	"github.com/spf13/viper"
)

// AppConfig 全局配置实例
var AppConfig *schemas.Config
var once sync.Once

// LoadConfig 加载配置文件
func LoadConfig(configFile string) error {
	once.Do(func() {
		viper.SetConfigFile(configFile)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("无法加载配置文件: %v", err)
		}

		if err := viper.Unmarshal(&AppConfig); err != nil {
			log.Fatalf("无法解析配置文件: %v", err)
		}
	})
	return nil
}
