package etcd

import (
	"context"
	client "github.com/coreos/etcd/clientv3"
)

type EtcdClient struct {
	Client *client.Client
	Context context.Context
	Cancel context.CancelFunc
}
