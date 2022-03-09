package builtin

import (
	"eago/common/log"
	"eago/common/redis"
	"eago/task/conf"
	"net"
)

// RegisterSrvWrapHandler 注册Srv
func RegisterSrvWrapHandler(conn net.Conn) {
	log.Info("RegisterSrvWrapHandler called.")
	defer log.Info("RegisterSrvWrapHandler end.")

	if conn != nil {
		RegisterAllowedSrv(
			conn.LocalAddr().String(),
			conn.RemoteAddr().String(),
		)
	}
}

// RegisterAllowedSrv 注册Srv白名单
func RegisterAllowedSrv(srvAddr, workerAddr string) {
	if err := redis.Redis.Set(genAllowedSrvKey(srvAddr), workerAddr, conf.Conf.AllowedSrvTtlSecs); err != nil {
		log.ErrorWithFields(log.Fields{
			"srv_addr":    srvAddr,
			"worker_addr": workerAddr,
			"error":       err,
		}, "An error occurred while RegisterAllowedSrv.")
	}
}

// UnregisterAllowedSrv 注销Srv白名单
func UnregisterAllowedSrv(srvAddr, workerAddr string) {
	if !redis.Redis.HasKey(genAllowedSrvKey(srvAddr)) {
		return
	}
	if err := redis.Redis.Del(genAllowedSrvKey(srvAddr)); err != nil {
		log.ErrorWithFields(log.Fields{
			"srv_addr":    srvAddr,
			"worker_addr": workerAddr,
			"error":       err,
		}, "An error occurred while UnregisterAllowedSrv, Skipped it.")
	}
}

// IsSrvAllowed 判断Srv是否已经在白名单内
func IsSrvAllowed(srvAddr, workerAddr string) bool {
	// 如果Key不存在，肯定没有注册srv
	if !redis.Redis.HasKey(genAllowedSrvKey(srvAddr)) {
		log.WarnWithFields(log.Fields{
			"srv_addr":    srvAddr,
			"worker_addr": workerAddr,
		}, "Srv was not registered.")
		return false
	}

	// 取出key对应的内容
	content, err := redis.Redis.Get(genAllowedSrvKey(srvAddr))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"srv_addr":    srvAddr,
			"worker_addr": workerAddr,
			"error":       err,
		}, "An error occurred while IsSrvAllowed.")
		return false
	}

	// 如果内容和Agent的地址不同，说明不允许
	if workerAddr != content {
		log.WarnWithFields(log.Fields{
			"srv_addr":    srvAddr,
			"worker_addr": workerAddr,
			"content":     content,
		}, "Srv not allowed.")
		return false
	}

	return true
}

// genAllowedSrvKey 生成白名单Key
func genAllowedSrvKey(srvAddr string) string {
	return "allowed_srv/" + srvAddr
}
