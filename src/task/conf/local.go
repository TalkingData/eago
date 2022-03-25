package conf

import (
	"github.com/Unknwon/goconfig"
	"strings"
	"time"
)

const (
	_DEFAULT_SRV_LISTEN              = "127.0.0.1:0"
	_DEFAULT_API_LISTEN              = "127.0.0.1:0"
	_DEFAULT_GIN_MODEL               = "release"
	_DEFAULT_MICRO_REGISTER_TTL      = 10
	_DEFAULT_MICRO_REGISTER_INTERVAL = 3
	_DEFAULT_SCHEDULER_REGISTER_TTL  = 10

	_DEFAULT_LOG_LEVEL = "debug"
	_DEFAULT_LOG_PATH  = "./logs"

	// Etcd默认配置
	_DEFAULT_ETCD_ADDRESSES = "127.0.0.1:2379,127.0.0.1:2379,127.0.0.1:2379"
	_DEFAULT_ETCD_USERNAME  = ""
	_DEFAULT_ETCD_PASSWORD  = ""

	// Mysql默认配置
	_DEFAULT_MYSQL_ADDRESS              = "127.0.0.1:3306"
	_DEFAULT_MYSQL_DB_NAME              = "eago_task"
	_DEFAULT_MYSQL_USER                 = "root"
	_DEFAULT_MYSQL_PASSWORD             = "root"
	_DEFAULT_MYSQL_MAX_OPEN_CONNECTIONS = 20
	_DEFAULT_MYSQL_MAX_IDLE_CONNECTIONS = 5

	// Redis默认配置
	_DEFAULT_REDIS_ADDRESS  = "127.0.0.1:6379"
	_DEFAULT_REDIS_PASSWORD = ""
	_DEFAULT_REDIS_DB       = 1

	// Jaeger默认配置
	_DEFAULT_JAEGER_ADDRESS = "127.0.0.1:5775"
)

// conf 配置
type conf struct {
	SrvListen             string
	ApiListen             string
	GinMode               string
	MicroRegisterTtl      time.Duration
	MicroRegisterInterval time.Duration
	SchedulerRegisterTtl  int64

	LogLevel string
	LogPath  string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	MysqlAddress            string
	MysqlDbName             string
	MysqlUser               string
	MysqlPassword           string
	MysqlMaxOpenConnections int
	MysqlMaxIdleConnections int

	RedisAddress  string
	RedisPassword string
	RedisDb       int

	JaegerAddress string
}

// 验证配置文件
func (c *conf) validateConfig() error {
	return nil
}

// newLocalConf 载入本地配置文件
func newLocalConf() *conf {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATHNAME)
	if err != nil {
		panic(err)
	}

	return &conf{
		SrvListen:             cfg.MustValue("main", "srv_listen", _DEFAULT_SRV_LISTEN),
		ApiListen:             cfg.MustValue("main", "api_listen", _DEFAULT_API_LISTEN),
		GinMode:               cfg.MustValue("main", "gin_mode", _DEFAULT_GIN_MODEL),
		MicroRegisterTtl:      time.Duration(cfg.MustInt("main", "micro_register_ttl", _DEFAULT_MICRO_REGISTER_TTL)) * time.Second,
		MicroRegisterInterval: time.Duration(cfg.MustInt("main", "micro_register_interval", _DEFAULT_MICRO_REGISTER_INTERVAL)) * time.Second,
		SchedulerRegisterTtl:  cfg.MustInt64("main", "scheduler_register_ttl", _DEFAULT_SCHEDULER_REGISTER_TTL),

		LogLevel: cfg.MustValue("log", "level", _DEFAULT_LOG_LEVEL),
		LogPath:  cfg.MustValue("log", "path", _DEFAULT_LOG_PATH),

		EtcdAddresses: strings.Split(cfg.MustValue("etcd", "addresses", _DEFAULT_ETCD_ADDRESSES), ","),
		EtcdUsername:  cfg.MustValue("etcd", "username", _DEFAULT_ETCD_USERNAME),
		EtcdPassword:  cfg.MustValue("etcd", "password", _DEFAULT_ETCD_PASSWORD),

		MysqlAddress:            cfg.MustValue("mysql", "address", _DEFAULT_MYSQL_ADDRESS),
		MysqlDbName:             cfg.MustValue("mysql", "db_name", _DEFAULT_MYSQL_DB_NAME),
		MysqlUser:               cfg.MustValue("mysql", "user", _DEFAULT_MYSQL_USER),
		MysqlPassword:           cfg.MustValue("mysql", "password", _DEFAULT_MYSQL_PASSWORD),
		MysqlMaxOpenConnections: cfg.MustInt("mysql", "max_open_connections", _DEFAULT_MYSQL_MAX_OPEN_CONNECTIONS),
		MysqlMaxIdleConnections: cfg.MustInt("mysql", "max_idle_connections", _DEFAULT_MYSQL_MAX_IDLE_CONNECTIONS),

		RedisAddress:  cfg.MustValue("redis", "address", _DEFAULT_REDIS_ADDRESS),
		RedisPassword: cfg.MustValue("redis", "password", _DEFAULT_REDIS_PASSWORD),
		RedisDb:       cfg.MustInt("redis", "db", _DEFAULT_REDIS_DB),

		JaegerAddress: cfg.MustValue("tracer", "jaeger_address", _DEFAULT_JAEGER_ADDRESS),
	}
}
