package conf

import (
	"github.com/Unknwon/goconfig"
	"os"
	"strings"
	"time"
)

const (
	_DEFAULT_SECRET_KEY = "eago_default_secret_key"
	// Token生命周期默认配置，秒
	_DEFAULT_TOKEN_TTL       = 900
	_DEFAULT_ADMIN_ROLE_NAME = "auth_admin"

	// Etcd默认配置
	_DEFAULT_ETCD_ADDRESSES = "127.0.0.1:2379,127.0.0.1:2379,127.0.0.1:2379"
	_DEFAULT_ETCD_USERNAME  = "root"
	_DEFAULT_ETCD_PASSWORD  = "root"

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

	_DEFAULT_LOG_LEVEL = "debug"
	_DEFAULT_LOG_PATH  = "./logs"
)

var Config *Conf

// Conf 配置
type Conf struct {
	LogLevel  string
	SecretKey string
	TokenTtl  time.Duration

	AdminRoleName string

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

	LogPath string
}

// init 初始化配置文件
func init() {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATHNAME)
	if err != nil {
		panic(err)
	}

	Config = &Conf{
		SecretKey:     cfg.MustValue("main", "secret_key", _DEFAULT_SECRET_KEY),
		TokenTtl:      time.Duration(cfg.MustInt64("main", "token_ttl", _DEFAULT_TOKEN_TTL)) * time.Second,
		AdminRoleName: cfg.MustValue("main", "admin_role_name", _DEFAULT_ADMIN_ROLE_NAME),

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
	}

	_ = os.Mkdir(Config.LogPath, 0755)
}
