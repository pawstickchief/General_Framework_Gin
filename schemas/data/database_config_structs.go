package data

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	ETCD  ETCDConfig  `mapstructure:"etcd"`
}

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxConns     int    `mapstructure:"max_conns"`
}

// ETCDConfig ETCD 数据库配置
type ETCDConfig struct {
	Endpoints   []string `mapstructure:"endpoints"`
	DialTimeout int      `mapstructure:"dial_timeout"`
	CACert      string   `mapstructure:"ca_cert"`
	CertFile    string   `mapstructure:"cert_file"`
	KeyFile     string   `mapstructure:"key_file"`
	ServerName  string   `mapstructure:"server_name"`
	EtcdName    string   `mapstructure:"etcdname"`
	Password    string   `mapstructure:"password"`
}
