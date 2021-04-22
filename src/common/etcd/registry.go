package etcd

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
)

var EtcdRegistry registry.Registry

// InitEtcd 初始化etcd
func InitEtcd(address, username, password string) {
	EtcdRegistry = etcdv3.NewRegistry(
		registry.Addrs(address),
		etcdv3.Auth(username, password),
	)
}
