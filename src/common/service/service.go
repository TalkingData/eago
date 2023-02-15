package service

type EagoSrv interface {
	// Start 启动EagoSrv
	Start() error
	// Stop 关闭EagoSrv
	Stop()
}
