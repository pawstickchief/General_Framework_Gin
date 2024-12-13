package schemas

import "General_Framework_Gin/schemas/data"

// Config 应用配置结构
type Config struct {
	AppName    string              `mapstructure:"app_name"`
	JWTSecret  string              `mapstructure:"jwt_secret"`
	Server     ServerConfig        `mapstructure:"server"`
	Log        LogConfig           `mapstructure:"log"`
	Database   data.DatabaseConfig `mapstructure:"database"`
	FileConfig FileStorageConfig   `mapstructure:"file_storage"`
	Update     UpdateConfig        `mapstructure:"update"`
}

// UpdateConfig 更新配置
type UpdateConfig struct {
	ServerURL string `mapstructure:"server_url"`
	Platform  string `mapstructure:"platform"`
	AppName   string `mapstructure:"app_name"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Address      string `mapstructure:"address"`
	Port         int    `mapstructure:"port"`
	RedirectPort int    `mapstructure:"redirect_port"`
	ClientURL    string `mapstructure:"client_url"`
	CertFile     string `mapstructure:"cert_file"`
	KeyFile      string `mapstructure:"key_file"`
}

// LogConfig 日志配置
type LogConfig struct {
	Mode       string `mapstructure:"mode"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type FileStorageConfig struct {
	UploadDir    string   `mapstructure:"upload_dir"`
	MaxFileSize  int64    `mapstructure:"max_file_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	EnableResume bool     `mapstructure:"enable_resume"`
	ChunkSize    int      `mapstructure:"chunk_size"`
}
