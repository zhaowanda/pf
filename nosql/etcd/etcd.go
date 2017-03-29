package etcd

import (
	"context"
	client "github.com/coreos/etcd/clientv3"
	"container/list"
	"strings"
	"time"
	"github.com/zhaowanda/pf/config"
	"github.com/zhaowanda/pf/config/yaml"
)

type EtcdClient struct {
	Client *client.Client
	Context context.Context
	Cancel context.CancelFunc
}


func InitEtcdClient() (EtcdClient, error) {
	// TODO 读取配置文件获取相应的配置
	endpoints, err := config.GetConfig("etcd.endpoints", yaml.YamlConfigure)
	if err != nil {
		return EtcdClient{}, err
	}
	timeout, err := config.GetConfig("etcd.timeout", yaml.YamlConfigure)
	if err != nil {
		return EtcdClient{}, err
	}
	username, err := config.GetConfig("etcd.username", yaml.YamlConfigure)
	if err != nil {
		return EtcdClient{}, err
	}
	password, err := config.GetConfig("etcd.password", yaml.YamlConfigure)
	if err != nil {
		return EtcdClient{}, err
	}
	contextTimeout, err := config.GetConfig("etcd.context.timeout", yaml.YamlConfigure)
	if err != nil {
		return EtcdClient{}, err
	}
	cfg := client.Config {
		Endpoints:   strings.Split(endpoints.(string), ","),
		DialTimeout: time.Duration(timeout.(int64)) * time.Second,
		Username: username.(string),
		Password: password.(string),
	}

	cli, err := client.New(cfg)
	if err != nil {
		return EtcdClient{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(contextTimeout.(int64)) * time.Second)
	client := EtcdClient{}
	client.Context = ctx
	client.Cancel = cancel
	client.Client = cli
	return client, nil
}



func (ec EtcdClient) Put(key, value string) (bool, error) {
	defer ec.Client.Close()
	_, err := ec.Client.Put(ec.Context, key, value)
	ec.Cancel()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (ec EtcdClient) Get(key string) (string, error) {
	defer ec.Client.Close()
	resp, err := ec.Client.Get(ec.Context, key)
	ec.Cancel()
	if err != nil {
		return "", err
	} else {
		result := resp.Kvs
		var resultValue string
		for _, value := range result {
			resultValue = string(value.Value)
		}
		return resultValue, nil
	}
}

func (ec EtcdClient) Del(key string) (bool, error) {
	defer ec.Client.Close()
	_, err := ec.Client.Delete(ec.Context, key)
	ec.Cancel()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ec EtcdClient) BatchPut(mapKey map[string]string) (bool, error) {
	flag := true
	var errMsg error
	for key, _ := range mapKey {
		result, err := ec.Put(key, mapKey[key])
		if err != nil {
			flag = result
			errMsg = err
			break
		}
	}
	return flag, errMsg
}

func (ec EtcdClient) BatchDel(keys list.List) (bool, error) {
	flag := true
	var errMsg error
	for it := keys.Front(); it != nil; it.Next() {
		key := it.Value.(string)
		result, err := ec.Del(key)
		if err != nil {
			flag = result
			break
		}
	}
	return flag, errMsg
}
