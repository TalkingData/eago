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
	SchedulerRegisterTtl  int64
	SrvTokenTtlSecs       time.Duration

	LogLevel string
	LogPath  string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	MysqlAddress      string
	MysqlDbName       string
	MysqlUser         string
	MysqlPassword     string
	MysqlMaxIdleConns int
	MysqlMaxOpenConns int

	RedisAddress  string
	RedisPassword string
	RedisDb       int

	JaegerAddress string
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
		SchedulerRegisterTtl: cfg.MustInt64("main", "scheduler_register_ttl", defaultSchedulerRegisterTtl),
		SrvTokenTtlSecs: time.Duration(cfg.MustInt(
			"main", "srv_token_ttl_secs", defaultSrvTokenTtlSecs,
		)) * time.Second,

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", global.DefaultConfigSeparator),
		EtcdUsername:  cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword:  cfg.MustValue("etcd", "password", defaultEtcdPassword),

		MysqlAddress:      cfg.MustValue("mysql", "address", defaultMysqlAddress),
		MysqlDbName:       cfg.MustValue("mysql", "db_name", defaultMysqlDbName),
		MysqlUser:         cfg.MustValue("mysql", "user", defaultMysqlUser),
		MysqlPassword:     cfg.MustValue("mysql", "password", defaultMysqlPassword),
		MysqlMaxIdleConns: cfg.MustInt("mysql", "max_idle_conns", defaultMysqlMaxIdleConns),
		MysqlMaxOpenConns: cfg.MustInt("mysql", "max_open_conns", defaultMysqlMaxOpenConns),

		RedisAddress:  cfg.MustValue("redis", "address", defaultRedisAddress),
		RedisPassword: cfg.MustValue("redis", "password", defaultRedisPassword),
		RedisDb:       cfg.MustInt("redis", "db", defaultRedisDb),

		JaegerAddress: cfg.MustValue("tracer", "jaeger_address", defaultJaegerAddress),
	}
}
