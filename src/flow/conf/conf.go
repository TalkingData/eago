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

	MysqlAddress      string
	MysqlDbName       string
	MysqlUser         string
	MysqlPassword     string
	MysqlMaxIdleConns int
	MysqlMaxOpenConns int

	KafkaAddresses []string

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

		ApiListen: cfg.MustValue("main", "api_listen", defaultApiListen),
		GinMode:   cfg.MustValue("main", "gin_mode", defaultGinModel),
		MicroRegisterTtl: time.Duration(cfg.MustInt(
			"main", "register_ttl", defaultMicroRegisterTtl,
		)) * time.Second,
		MicroRegisterInterval: time.Duration(cfg.MustInt(
			"main", "register_interval", defaultMicroRegisterInterval,
		)) * time.Second,

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		NotifyTitle:   cfg.MustValue("notify", "title", defaultNotifyTitle),
		NotifyBaseUrl: cfg.MustValue("notify", "base_url", defaultNotifyBaseUrl),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", global.DefaultConfigSeparator),
		EtcdUsername:  cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword:  cfg.MustValue("etcd", "password", defaultEtcdPassword),

		MysqlAddress:      cfg.MustValue("mysql", "address", defaultMysqlAddress),
		MysqlDbName:       cfg.MustValue("mysql", "db_name", defaultMysqlDbName),
		MysqlUser:         cfg.MustValue("mysql", "user", defaultMysqlUser),
		MysqlPassword:     cfg.MustValue("mysql", "password", defaultMysqlPassword),
		MysqlMaxIdleConns: cfg.MustInt("mysql", "max_idle_conns", defaultMysqlMaxIdleConns),
		MysqlMaxOpenConns: cfg.MustInt("mysql", "max_open_conns", defaultMysqlMaxOpenConns),

		KafkaAddresses: cfg.MustValueArray("kafka", "addresses", global.DefaultConfigSeparator),

		JaegerAddress: cfg.MustValue("tracer", "jaeger_address", defaultJaegerAddress),
	}
}
