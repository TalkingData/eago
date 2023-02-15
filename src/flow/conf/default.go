package conf

const (
	defaultApiListen             = "127.0.0.1:0"
	defaultGinModel              = "release"
	defaultMicroRegisterTtl      = 10
	defaultMicroRegisterInterval = 3

	defaultLogLevel = "debug"
	defaultLogPath  = "./logs"

	defaultNotifyTitle   = "[CMDB-Eago]Flow notify"
	defaultNotifyBaseUrl = "https://eago.tendcloud.com"

	// Etcd默认配置
	defaultEtcdUsername = ""
	defaultEtcdPassword = ""

	// Mysql默认配置
	defaultMysqlAddress      = "127.0.0.1:3306"
	defaultMysqlDbName       = "eago_flow"
	defaultMysqlUser         = "root"
	defaultMysqlPassword     = "root"
	defaultMysqlMaxOpenConns = 20
	defaultMysqlMaxIdleConns = 5

	// Jaeger默认配置
	defaultJaegerAddress = "127.0.0.1:5775"
)
