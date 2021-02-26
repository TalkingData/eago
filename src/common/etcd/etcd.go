package etcd

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
)

var EtcdReg registry.Registry

func InitEtcd(address string, username string, password string) {
	EtcdReg = etcdv3.NewRegistry(
		registry.Addrs(address),
		etcdv3.Auth(username, password),
	)
}
