package conf

const (
	defaultSrvListen             = "127.0.0.1:0"
	defaultApiListen             = "127.0.0.1:0"
	defaultGinModel              = "release"
	defaultMicroRegisterTtl      = 10
	defaultMicroRegisterInterval = 3
	defaultSchedulerRegisterTtl  = 10
	defaultSrvTokenTtlSecs       = 10

	defaultLogLevel = "debug"
	defaultLogPath  = "./logs"

	// Etcd默认配置
	defaultEtcdUsername = ""
	defaultEtcdPassword = ""

	// Mysql默认配置
	defaultMysqlAddress      = "127.0.0.1:3306"
	defaultMysqlDbName       = "eago_task"
	defaultMysqlUser         = "root"
	defaultMysqlPassword     = "root"
	defaultMysqlMaxOpenConns = 20
	defaultMysqlMaxIdleConns = 5

	// Redis默认配置
	defaultRedisAddress  = "127.0.0.1:6379"
	defaultRedisPassword = ""
	defaultRedisDb       = 1

	// Jaeger默认配置
	defaultJaegerAddress = "127.0.0.1:5775"
)
