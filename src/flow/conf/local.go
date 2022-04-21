package conf

import (
	"github.com/Unknwon/goconfig"
	"time"
)

const (
	_DEFAULT_API_LISTEN              = "127.0.0.1:0"
	_DEFAULT_GIN_MODEL               = "release"
	_DEFAULT_MICRO_REGISTER_TTL      = 10
	_DEFAULT_MICRO_REGISTER_INTERVAL = 3

	_DEFAULT_LOG_LEVEL = "debug"
	_DEFAULT_LOG_PATH  = "./logs"

	_DEFAULT_NOTIFY_TITLE    = "[CMDB-Eago]Flow notify"
	_DEFAULT_NOTIFY_BASE_URL = "https://eago.tendcloud.com"

	// Etcd默认配置
	_DEFAULT_ETCD_ADDRESSES = "127.0.0.1:2379,127.0.0.1:2379,127.0.0.1:2379"
	_DEFAULT_ETCD_USERNAME  = ""
	_DEFAULT_ETCD_PASSWORD  = ""

	// Mysql默认配置
	_DEFAULT_MYSQL_ADDRESS              = "127.0.0.1:3306"
	_DEFAULT_MYSQL_DB_NAME              = "eago_flow"
	_DEFAULT_MYSQL_USER                 = "root"
	_DEFAULT_MYSQL_PASSWORD             = "root"
	_DEFAULT_MYSQL_MAX_OPEN_CONNECTIONS = 20
	_DEFAULT_MYSQL_MAX_IDLE_CONNECTIONS = 5

	// Jaeger默认配置
	_DEFAULT_JAEGER_ADDRESS = "127.0.0.1:5775"
)

// conf 配置
type conf struct {
	ApiListen             string
	GinMode               string
	MicroRegisterTtl      time.Duration
	MicroRegisterInterval time.Duration

	LogLevel string
	LogPath  string

	NotifyTitle   string
	NotifyBaseUrl string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	MysqlAddress            string
	MysqlDbName             string
	MysqlUser               string
	MysqlPassword           string
	MysqlMaxOpenConnections int
	MysqlMaxIdleConnections int

	JaegerAddress string
}

// newLocalConf 载入本地配置文件
func newLocalConf() *conf {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATHNAME)
	if err != nil {
		panic(err)
	}

	return &conf{
		ApiListen:             cfg.MustValue("main", "api_listen", _DEFAULT_API_LISTEN),
		GinMode:               cfg.MustValue("main", "gin_mode", _DEFAULT_GIN_MODEL),
		MicroRegisterTtl:      time.Duration(cfg.MustInt("main", "micro_register_ttl", _DEFAULT_MICRO_REGISTER_TTL)) * time.Second,
		MicroRegisterInterval: time.Duration(cfg.MustInt("main", "micro_register_interval", _DEFAULT_MICRO_REGISTER_INTERVAL)) * time.Second,

		LogLevel: cfg.MustValue("log", "level", _DEFAULT_LOG_LEVEL),
		LogPath:  cfg.MustValue("log", "path", _DEFAULT_LOG_PATH),

		NotifyTitle:   cfg.MustValue("notify", "title", _DEFAULT_NOTIFY_TITLE),
		NotifyBaseUrl: cfg.MustValue("notify", "base_url", _DEFAULT_NOTIFY_BASE_URL),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", _DEFAULT_CONFIG_SEPARATOR),
		EtcdUsername:  cfg.MustValue("etcd", "username", _DEFAULT_ETCD_USERNAME),
		EtcdPassword:  cfg.MustValue("etcd", "password", _DEFAULT_ETCD_PASSWORD),

		MysqlAddress:            cfg.MustValue("mysql", "address", _DEFAULT_MYSQL_ADDRESS),
		MysqlDbName:             cfg.MustValue("mysql", "db_name", _DEFAULT_MYSQL_DB_NAME),
		MysqlUser:               cfg.MustValue("mysql", "user", _DEFAULT_MYSQL_USER),
		MysqlPassword:           cfg.MustValue("mysql", "password", _DEFAULT_MYSQL_PASSWORD),
		MysqlMaxOpenConnections: cfg.MustInt("mysql", "max_open_connections", _DEFAULT_MYSQL_MAX_OPEN_CONNECTIONS),
		MysqlMaxIdleConnections: cfg.MustInt("mysql", "max_idle_connections", _DEFAULT_MYSQL_MAX_IDLE_CONNECTIONS),

		JaegerAddress: cfg.MustValue("tracer", "jaeger_address", _DEFAULT_JAEGER_ADDRESS),
	}
}
