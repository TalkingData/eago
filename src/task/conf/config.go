package conf

import (
	"github.com/Unknwon/goconfig"
	"os"
	"strings"
)

const (
	_DEFAULT_LOG_LEVEL = "debug"
	_DEFAULT_LOG_PATH  = "./logs"

	// Etcd默认配置
	_DEFAULT_ETCD_ADDRESSES = "127.0.0.1:2379,127.0.0.1:2379,127.0.0.1:2379"
	_DEFAULT_ETCD_USERNAME  = "root"
	_DEFAULT_ETCD_PASSWORD  = "root"

	// Mysql默认配置
	_DEFAULT_MYSQL_ADDRESS              = "127.0.0.1:3306"
	_DEFAULT_MYSQL_DB_NAME              = "eago_task"
	_DEFAULT_MYSQL_USER                 = "root"
	_DEFAULT_MYSQL_PASSWORD             = "root"
	_DEFAULT_MYSQL_MAX_OPEN_CONNECTIONS = 20
	_DEFAULT_MYSQL_MAX_IDLE_CONNECTIONS = 5
)

var Config *Conf

// Conf 配置
type Conf struct {
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
}

// 验证配置文件
func (c *Conf) validateConfig() error {
	return nil
}

// init 初始化配置文件
func init() {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATHNAME)
	if err != nil {
		panic(err)
	}

	Config = &Conf{
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
	}

	_ = os.Mkdir(Config.LogPath, 0755)
}
