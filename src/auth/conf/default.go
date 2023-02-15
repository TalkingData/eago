package conf

const (
	defaultSrvListen             = "127.0.0.1:0"
	defaultApiListen             = "127.0.0.1:0"
	defaultGinModel              = "release"
	defaultMicroRegisterTtl      = 10
	defaultMicroRegisterInterval = 3
	defaultWorkerRegisterTtl     = 10
	defaultTokenTtl              = 900
	defaultSecretKey             = "eago_default_secret_key"

	defaultLogLevel = "debug"
	defaultLogPath  = "./logs"

	// Etcd默认配置
	defaultEtcdUsername = ""
	defaultEtcdPassword = ""

	// Crowd默认配置
	defaultCrowdAddress     = "https://127.0.0.1/crowd"
	defaultCrowdAppName     = "eago"
	defaultCrowdAppPassword = "eago"

	// Mysql默认配置
	defaultMysqlAddress      = "127.0.0.1:3306"
	defaultMysqlDbName       = "eago_auth"
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
