package config

import (
	"github.com/Unknwon/goconfig"
	"os"
	"time"
)

const (
	// 配置文件路径
	configFilePathName = "./config/eago_auth.conf"

	defaultServiceName      = "eago-auth"
	defaultRpcServiceName   = "eago.srv.auth"
	defaultApiV1ServiceName = "eago.api.v1.auth"
	defaultSecretKey        = "eago_default_secret_key"
	// Token生命周期默认配置，秒
	defaultTokenTtl = 3600

	// Etcd默认配置
	defaultEtcdAddress  = "127.0.0.1:2379"
	defaultEtcdUsername = "root"
	defaultEtcdPassword = "root"

	// Crowd默认配置
	defaultCrowdAddress     = "127.0.0.1"
	defaultCrowdAppName     = "eago"
	defaultCrowdAppPassword = "eago"

	// Mysql默认配置
	defaultMysqlAddress            = "127.0.0.1:3306"
	defaultMysqlDbName             = "eago_auth"
	defaultMysqlUser               = "root"
	defaultMysqlPassword           = "root"
	defaultMysqlMaxOpenConnections = 20
	defaultMysqlMaxIdleConnections = 5

	// Redis默认配置
	defaultRedisAddress  = "127.0.0.1:6379"
	defaultRedisPassword = ""
	defaultRedisDb       = 1

	defaultLogPath = "./logs"
)

var Config *Conf

// Conf 配置
type Conf struct {
	ServiceName      string
	RpcServiceName   string
	ApiV1ServiceName string
	SecretKey        string
	TokenTtl         time.Duration

	EtcdAddress  string
	EtcdUsername string
	EtcdPassword string

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

// 验证配置文件
func (c *Conf) validateConfig() error {
	return nil
}

func InitConfig() error {
	cfg, err := goconfig.LoadConfigFile(configFilePathName)
	if err != nil {
		return err
	}

	Config = &Conf{
		ServiceName:      cfg.MustValue("main", "service_name", defaultServiceName),
		RpcServiceName:   cfg.MustValue("main", "rpc_service_name", defaultRpcServiceName),
		ApiV1ServiceName: cfg.MustValue("main", "api_v1_service_name", defaultApiV1ServiceName),
		SecretKey:        cfg.MustValue("main", "secret_key", defaultSecretKey),
		TokenTtl:         time.Duration(cfg.MustInt64("main", "token_ttl", defaultTokenTtl)) * time.Second,

		EtcdAddress:  cfg.MustValue("etcd", "address", defaultEtcdAddress),
		EtcdUsername: cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword: cfg.MustValue("etcd", "password", defaultEtcdPassword),

		CrowdAddress:     cfg.MustValue("crowd", "crowd_url", defaultCrowdAddress),
		CrowdAppName:     cfg.MustValue("crowd", "app_name", defaultCrowdAppName),
		CrowdAppPassword: cfg.MustValue("crowd", "app_pass", defaultCrowdAppPassword),

		MysqlAddress:            cfg.MustValue("mysql", "address", defaultMysqlAddress),
		MysqlDbName:             cfg.MustValue("mysql", "db_name", defaultMysqlDbName),
		MysqlUser:               cfg.MustValue("mysql", "user", defaultMysqlUser),
		MysqlPassword:           cfg.MustValue("mysql", "password", defaultMysqlPassword),
		MysqlMaxOpenConnections: cfg.MustInt("mysql", "max_open_connections", defaultMysqlMaxOpenConnections),
		MysqlMaxIdleConnections: cfg.MustInt("mysql", "max_idle_connections", defaultMysqlMaxIdleConnections),

		RedisAddress:  cfg.MustValue("redis", "address", defaultRedisAddress),
		RedisPassword: cfg.MustValue("redis", "password", defaultRedisPassword),
		RedisDb:       cfg.MustInt("redis", "db", defaultRedisDb),

		LogPath: cfg.MustValue("log", "path", defaultLogPath),
	}

	os.Mkdir(Config.LogPath, 0755)

	return Config.validateConfig()
}
