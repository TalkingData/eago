package conf

import (
	"github.com/Unknwon/goconfig"
	"strings"
	"time"
)

const (
	_DEFAULT_SECRET_KEY = "eago_default_secret_key"
	// Token生命周期默认配置，秒
	_DEFAULT_TOKEN_TTL = 900
	_DEFAULT_GIN_MODEL = "release"

	// Etcd默认配置
	_DEFAULT_ETCD_ADDRESSES = "127.0.0.1:2379,127.0.0.1:2379,127.0.0.1:2379"
	_DEFAULT_ETCD_USERNAME  = ""
	_DEFAULT_ETCD_PASSWORD  = ""

	// IAM默认配置
	_DEFAULT_IAM_ADDRESS = "https://127.0.0.1/auth"

	// Crowd默认配置
	_DEFAULT_CROWD_ADDRESS      = "http://127.0.0.1/crowd"
	_DEFAULT_CROWD_APP_NAME     = "eago"
	_DEFAULT_CROWD_APP_PASSWORD = "eago"

	// Mysql默认配置
	_DEFAULT_MYSQL_ADDRESS              = "127.0.0.1:3306"
	_DEFAULT_MYSQL_DB_NAME              = "eago_auth"
	_DEFAULT_MYSQL_USER                 = "root"
	_DEFAULT_MYSQL_PASSWORD             = "root"
	_DEFAULT_MYSQL_MAX_OPEN_CONNECTIONS = 20
	_DEFAULT_MYSQL_MAX_IDLE_CONNECTIONS = 5

	// Redis默认配置
	_DEFAULT_REDIS_ADDRESS  = "127.0.0.1:6379"
	_DEFAULT_REDIS_PASSWORD = ""
	_DEFAULT_REDIS_DB       = 1

	// Kafka默认配置
	_DEFAULT_KAFKA_ADDRESSES = "127.0.0.1:9092,127.0.0.1:9092,127.0.0.1:9092"

	// Jaeger默认配置
	_DEFAULT_JAEGER_ADDRESS = "127.0.0.1:5775"

	_DEFAULT_LOG_LEVEL = "debug"
	_DEFAULT_LOG_PATH  = "./logs"
)

// conf 配置
type conf struct {
	LogLevel  string
	SecretKey string
	TokenTtl  time.Duration
	GinMode   string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	IamAddress string

	CrowdAddress     string
	CrowdAppName     string
	CrowdAppPassword string

	MysqlAddress            string
	MysqlDbName             string
	MysqlUser               string
	MysqlPassword           string
	MysqlMaxOpenConnections int
	MysqlMaxIdleConnections int

	RedisAddress  string
	RedisPassword string
	RedisDb       int

	KafkaAddresses []string

	JaegerAddress string

	// 企业微信配置
	WeworkAgentId    string
	WeworkCorpId     string
	WeworkCorpSecret string

	// Eagle配置
	EagleAddress  string
	EagleDbName   string
	EagleUser     string
	EaglePassword string

	LogPath string
}

// newLocalConf 载入本地配置文件
func newLocalConf() *conf {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATHNAME)
	if err != nil {
		panic(err)
	}

	return &conf{
		SecretKey: cfg.MustValue("main", "secret_key", _DEFAULT_SECRET_KEY),
		TokenTtl:  time.Duration(cfg.MustInt64("main", "token_ttl", _DEFAULT_TOKEN_TTL)) * time.Second,
		GinMode:   cfg.MustValue("main", "gin_mode", _DEFAULT_GIN_MODEL),

		LogLevel: cfg.MustValue("log", "level", _DEFAULT_LOG_LEVEL),
		LogPath:  cfg.MustValue("log", "path", _DEFAULT_LOG_PATH),

		EtcdAddresses: strings.Split(cfg.MustValue("etcd", "addresses", _DEFAULT_ETCD_ADDRESSES), ","),
		EtcdUsername:  cfg.MustValue("etcd", "username", _DEFAULT_ETCD_USERNAME),
		EtcdPassword:  cfg.MustValue("etcd", "password", _DEFAULT_ETCD_PASSWORD),

		IamAddress: cfg.MustValue("iam", "address", _DEFAULT_IAM_ADDRESS),

		CrowdAddress:     cfg.MustValue("crowd", "address", _DEFAULT_CROWD_ADDRESS),
		CrowdAppName:     cfg.MustValue("crowd", "app_name", _DEFAULT_CROWD_APP_NAME),
		CrowdAppPassword: cfg.MustValue("crowd", "app_pass", _DEFAULT_CROWD_APP_PASSWORD),

		MysqlAddress:            cfg.MustValue("mysql", "address", _DEFAULT_MYSQL_ADDRESS),
		MysqlDbName:             cfg.MustValue("mysql", "db_name", _DEFAULT_MYSQL_DB_NAME),
		MysqlUser:               cfg.MustValue("mysql", "user", _DEFAULT_MYSQL_USER),
		MysqlPassword:           cfg.MustValue("mysql", "password", _DEFAULT_MYSQL_PASSWORD),
		MysqlMaxOpenConnections: cfg.MustInt("mysql", "max_open_connections", _DEFAULT_MYSQL_MAX_OPEN_CONNECTIONS),
		MysqlMaxIdleConnections: cfg.MustInt("mysql", "max_idle_connections", _DEFAULT_MYSQL_MAX_IDLE_CONNECTIONS),

		RedisAddress:  cfg.MustValue("redis", "address", _DEFAULT_REDIS_ADDRESS),
		RedisPassword: cfg.MustValue("redis", "password", _DEFAULT_REDIS_PASSWORD),
		RedisDb:       cfg.MustInt("redis", "db", _DEFAULT_REDIS_DB),

		KafkaAddresses: strings.Split(cfg.MustValue("kafka", "addresses", _DEFAULT_KAFKA_ADDRESSES), ","),

		JaegerAddress: cfg.MustValue("tracer", "jaeger_address", _DEFAULT_JAEGER_ADDRESS),

		// 企业微信配置
		WeworkAgentId:    cfg.MustValue("wework", "agent_id"),
		WeworkCorpId:     cfg.MustValue("wework", "corp_id"),
		WeworkCorpSecret: cfg.MustValue("wework", "corp_secret"),

		// Eagle配置
		EagleAddress:  cfg.MustValue("eagle", "address"),
		EagleDbName:   cfg.MustValue("eagle", "db_name"),
		EagleUser:     cfg.MustValue("eagle", "user"),
		EaglePassword: cfg.MustValue("eagle", "password"),
	}
}