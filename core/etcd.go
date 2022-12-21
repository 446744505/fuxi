package core

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var ETCD *etcd

type OnWatcher interface {
	OnAdd(key, val string)
	OnDelete(key, val string)
}

type node struct {
	OnWatcher
	kvs map[string]string
}

type etcd struct {
	conf clientv3.Config
	kvs map[string]string

	lock sync.RWMutex
	nodes map[string]*node

	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancelFunc     func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func InitEtcd(addr []string) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	ETCD = &etcd{
		conf: conf,
		kvs: make(map[string]string),
		nodes: make(map[string]*node),
	}
	go ETCD.init()
}

func (self *etcd) Put(key, val string) {
	self.kvs[key] = val
	if self.client != nil {
		self.put(key, val)
	}
}

func (self *etcd) Delete(key string) {
	delete(self.kvs, key)
	if self.client != nil {
		self.delete(key)
	}
}

func (self *etcd) Watch(prefix string, onWatcher OnWatcher) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.nodes[prefix] = &node{
		OnWatcher: onWatcher,
		kvs: make(map[string]string),
	}
	if self.client != nil {
		go self.watcher(prefix)
	}
}

func (self *etcd) init() error {
	if client, err := clientv3.New(self.conf); err != nil {
		return err
	} else {
		self.client = client
	}

	if err := self.setLease(5); err != nil {
		return err
	}
	for k, v := range self.kvs {
		self.put(k, v)
	}
	for prefix, _ := range self.nodes {
		go self.watcher(prefix)
	}
	go self.listenLeaseRespChan()
	return nil
}

//设置租约
func (self *etcd) setLease(timeNum int64) error {
	lease := clientv3.NewLease(self.client)

	//设置租约时间
	leaseResp, err := lease.Grant(context.TODO(), timeNum)
	if err != nil {
		return err
	}

	//设置续租
	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	self.lease = lease
	self.leaseResp = leaseResp
	self.cancelFunc = cancelFunc
	self.keepAliveChan = leaseRespChan
	return nil
}

//监听 续租情况
func (self *etcd) listenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-self.keepAliveChan:
			if leaseKeepResp == nil {
				Log.Infoln("etcd server closed")
				self.client = nil
				self.init()
			}
		}
	}
}

//通过租约 注册服务
func (self *etcd) put(key, val string) error {
	kv := clientv3.NewKV(self.client)
	if _, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(self.leaseResp.ID)); err != nil {
		return err
	}
	Log.Infof("etcd put kv %s = %s", key, val)
	return nil
}

func (self *etcd) delete(key string) error {
	kv := clientv3.NewKV(self.client)
	if _, err := kv.Delete(context.TODO(), key, clientv3.WithLease(self.leaseResp.ID)); err != nil {
		return err
	}
	Log.Infof("etcd delete kv %s", key)
	return nil
}

func (self *etcd) watcher(prefix string) error {
	rsp, err := self.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if rsp != nil && rsp.Kvs != nil {
		for _, kv := range rsp.Kvs {
			self.putNode(prefix, string(kv.Key), string(kv.Value))
		}
	}

	rch := self.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for rsp := range rch {
		for _, ev := range rsp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				self.putNode(prefix, string(ev.Kv.Key),string(ev.Kv.Value))
			case mvccpb.DELETE:
				self.deleteNode(prefix, string(ev.Kv.Key))
			}
		}
	}
	return nil
}

func (self *etcd) putNode(prefix, key, val string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.nodes[prefix].put(key, val)
}

func (self *etcd) deleteNode(prefix, key string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if node, ok := self.nodes[prefix]; ok {
		node.delete(key)
	}
}

//撤销租约
func (self *etcd) RevokeLease() error {
	self.cancelFunc()
	_, err := self.lease.Revoke(context.TODO(), self.leaseResp.ID)
	return err
}

func (self *node) put(key, val string) {
	self.kvs[key] = val
	self.OnAdd(key, val)
	Log.Infof("etcd add node %s = %s", key, val)
}

func (self *node) delete(key string) {
	if val, ok := self.kvs[key]; ok {
		delete(self.kvs, key)
		self.OnDelete(key, val)
		Log.Infof("etcd delete node key %s", key)
	}
}