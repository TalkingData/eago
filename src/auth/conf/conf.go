package conf

import (
	"eago/common/global"
	"fmt"
	"github.com/Unknwon/goconfig"
	"time"
)

// Conf 配置
type Conf struct {
	Const *constConf

	SrvListen             string
	ApiListen             string
	GinMode               string
	MicroRegisterTtl      time.Duration
	MicroRegisterInterval time.Duration
	WorkerRegisterTtl     int64
	TokenTtl              time.Duration
	SecretKey             string

	LogLevel string
	LogPath  string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	CrowdAddress     string
	CrowdAppName     string
	CrowdAppPassword string

	MysqlAddress      string
	MysqlDbName       string
	MysqlUser         string
	MysqlPassword     string
	MysqlMaxIdleConns int
	MysqlMaxOpenConns int

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
}

func NewConfig(options ...Option) *Conf {
	opts := newOptions(options...)

	fmt.Println(fmt.Sprintf("Loading config file: %s", opts.ConfFilePathname))
	cfg, err := goconfig.LoadConfigFile(opts.ConfFilePathname)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Const: newConstConf(),

		SrvListen: cfg.MustValue("main", "srv_listen", defaultSrvListen),
		ApiListen: cfg.MustValue("main", "api_listen", defaultApiListen),
		GinMode:   cfg.MustValue("main", "gin_mode", defaultGinModel),
		MicroRegisterTtl: time.Duration(cfg.MustInt(
			"main", "register_ttl", defaultMicroRegisterTtl,
		)) * time.Second,
		MicroRegisterInterval: time.Duration(cfg.MustInt(
			"main", "register_interval", defaultMicroRegisterInterval,
		)) * time.Second,
		WorkerRegisterTtl: cfg.MustInt64("main", "worker_register_ttl", defaultWorkerRegisterTtl),
		TokenTtl: time.Duration(
			cfg.MustInt64("main", "token_ttl", defaultTokenTtl)) * time.Second,
		SecretKey: cfg.MustValue("main", "secret_key", defaultSecretKey),

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", global.DefaultConfigSeparator),
		EtcdUsername:  cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword:  cfg.MustValue("etcd", "password", defaultEtcdPassword),

		CrowdAddress:     cfg.MustValue("crowd", "address", defaultCrowdAddress),
		CrowdAppName:     cfg.MustValue("crowd", "app_name", defaultCrowdAppName),
		CrowdAppPassword: cfg.MustValue("crowd", "app_pass", defaultCrowdAppPassword),

		MysqlAddress:      cfg.MustValue("mysql", "address", defaultMysqlAddress),
		MysqlDbName:       cfg.MustValue("mysql", "db_name", defaultMysqlDbName),
		MysqlUser:         cfg.MustValue("mysql", "user", defaultMysqlUser),
		MysqlPassword:     cfg.MustValue("mysql", "password", defaultMysqlPassword),
		MysqlMaxIdleConns: cfg.MustInt("mysql", "max_idle_conns", defaultMysqlMaxIdleConns),
		MysqlMaxOpenConns: cfg.MustInt("mysql", "max_open_conns", defaultMysqlMaxOpenConns),

		RedisAddress:  cfg.MustValue("redis", "address", defaultRedisAddress),
		RedisPassword: cfg.MustValue("redis", "password", defaultRedisPassword),
		RedisDb:       cfg.MustInt("redis", "db", defaultRedisDb),

		KafkaAddresses: cfg.MustValueArray("kafka", "addresses", global.DefaultConfigSeparator),

		JaegerAddress: cfg.MustValue("tracer", "jaeger_address", defaultJaegerAddress),

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
