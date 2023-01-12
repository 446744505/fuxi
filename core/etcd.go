package core

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ETCD *etcd

const (
	NodeNameLinker   = "linker"
	NodeNameProvider = "provider"
	NodeNameProvidee = "providee"
)

type OnWatcher interface {
	OnAdd(key, val string)
	OnDelete(key, val string)
}

type SwitcherMeta struct {
	NodeName    string
	LinkerUrl   string
	ProviderUrl string
}

type ProvideeMeta struct {
	NodeName    string
	ServerName  string
	ProviderUrl string
	Pvid        int32
}

type node struct {
	OnWatcher
	kvs map[string]string
}

type pair struct {
	isDelete bool
	key      string
	val      interface{}
}

type etcd struct {
	conf clientv3.Config

	closeSig chan struct{}

	kvChan chan *pair
	kvs    map[string]string

	nodesChan chan *pair
	nodes     sync.Map

	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancelFunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func InitEtcd(addr []string) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	ETCD = &etcd{
		conf:      conf,
		kvChan:    make(chan *pair, 1000),
		kvs:       make(map[string]string),
		nodesChan: make(chan *pair, 1000),
		closeSig:  make(chan struct{}),
	}
	go ETCD.startWork()
}

func StopEtcd() {
	close(ETCD.closeSig)
}

func (self *etcd) Put(key, val string) error {
	select {
	case self.kvChan <- &pair{key: key, val: val}:
		return nil
	default:
		return errors.New("etcd kv chan is full")
	}
}

func (self *etcd) Delete(key string) error {
	select {
	case self.kvChan <- &pair{isDelete: true, key: key}:
		return nil
	default:
		return errors.New("etcd kv chan is full")
	}
}

func (self *etcd) Watch(prefix string, onWatcher OnWatcher) {
	n := &node{
		OnWatcher: onWatcher,
		kvs:       make(map[string]string),
	}
	self.nodesChan <- &pair{
		key: prefix,
		val: n,
	}
}

func (self *etcd) startWork() error {
	if client, err := clientv3.New(self.conf); err != nil {
		return err
	} else {
		self.client = client
	}

	if err := self.setLease(5); err != nil {
		Log.Errorf("etcd start setLease err %v", err)
		return err
	}

	for {
		select {
		case <-self.closeSig:
			self.clean()
			if err := self.doRevokeLease(); err != nil {
				Log.Errorf("etcd revoke lease err %v", err)
			}
			if err := self.client.Close(); err != nil {
				Log.Errorf("etcd close client err %v", err)
			}
			return nil
		case p := <-self.kvChan:
			if p.isDelete {
				delete(self.kvs, p.key)
				if err := self.doDelete(p.key); err != nil {
					Log.Errorf("etcd delete key %v err %v", p.key, err)
				}
			} else {
				self.kvs[p.key] = fmt.Sprint(p.val)
				if err := self.doPut(p.key, fmt.Sprint(p.val)); err != nil {
					Log.Errorf("etcd put key %v val %v err %v", p.key, p.val, err)
				}
			}
		case p := <-self.nodesChan:
			self.nodes.Store(p.key, p.val.(*node))
			go self.doWatcher(p.key)
		case rsp := <-self.keepAliveChan:
			if rsp == nil {
				Log.Infoln("etcd server closed")
				self.leaseResp = nil
				if self.cancelFunc != nil {
					self.cancelFunc()
				}
				if err := self.setLease(5); err != nil {
					Log.Errorf("etcd reconnect setLease err %v", err)
					continue
				}
				for k, v := range self.kvs {
					self.doPut(k, v)
				}
			}
		}
	}
}

//设置租约
func (self *etcd) setLease(timeNum int64) error {
	lease := clientv3.NewLease(self.client)

	//设置租约时间，阻塞到服务可用
	leaseResp, err := lease.Grant(context.TODO(), timeNum)
	if err != nil {
		return err
	}

	//设置续租
	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		cancelFunc()
		return err
	}

	self.lease = lease
	self.leaseResp = leaseResp
	self.cancelFunc = cancelFunc
	self.keepAliveChan = leaseRespChan

	return nil
}

func (self *etcd) doPut(key, val string) error {
	kv := clientv3.NewKV(self.client)
	if _, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(self.leaseResp.ID)); err != nil {
		return err
	}
	Log.Infof("etcd put kv %s = %s", key, val)
	return nil
}

func (self *etcd) doDelete(key string) error {
	kv := clientv3.NewKV(self.client)
	if _, err := kv.Delete(context.TODO(), key); err != nil {
		return err
	}
	Log.Infof("etcd delete kv %s", key)
	return nil
}

func (self *etcd) doWatcher(prefix string) error {
	rsp, err := self.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		Log.Errorf("etcd watcher prefix %v err %v", prefix, err)
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
				self.putNode(prefix, string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				self.deleteNode(prefix, string(ev.Kv.Key))
			}
		}
	}
	return nil
}

func (self *etcd) putNode(prefix, key, val string) {
	if value, ok := self.nodes.Load(prefix); ok {
		node := value.(*node)
		node.put(key, val)
	}
}

func (self *etcd) deleteNode(prefix, key string) {
	if value, ok := self.nodes.Load(prefix); ok {
		node := value.(*node)
		node.delete(key)
	}
}

//撤销租约
func (self *etcd) doRevokeLease() error {
	self.cancelFunc()
	_, err := self.lease.Revoke(context.TODO(), self.leaseResp.ID)
	return err
}

func (self *etcd) clean() {
	close(self.kvChan)
	close(self.nodesChan)
	self.kvs = nil
	self.nodes = sync.Map{}
}

func (self *node) put(key, val string) {
	self.kvs[key] = val
	Log.Infof("etcd add node %s = %s", key, val)
	self.OnAdd(key, val)
}

func (self *node) delete(key string) {
	if val, ok := self.kvs[key]; ok {
		delete(self.kvs, key)
		Log.Infof("etcd delete node key %s", key)
		self.OnDelete(key, val)
	}
}

func (self *SwitcherMeta) Path() string {
	return fmt.Sprintf("%s/%s/%s", self.NodeName, self.LinkerUrl, self.ProviderUrl)
}

func (self *SwitcherMeta) ValueOf(str string) *SwitcherMeta {
	arr := strings.Split(str, "/")
	self.NodeName = arr[0]
	self.LinkerUrl = arr[1]
	self.ProviderUrl = arr[2]
	return self
}

func (self *ProvideeMeta) Path() string {
	return fmt.Sprintf("%s/%s/%v", self.NodeName, self.ProviderUrl, self.Pvid)
}

func (self *ProvideeMeta) ValueOf(str string, val string) *ProvideeMeta {
	arr := strings.Split(str, "/")
	self.NodeName = arr[0]
	self.ProviderUrl = arr[1]
	self.ServerName = val
	pvid, _ := strconv.Atoi(arr[2])
	self.Pvid = int32(pvid)
	return self
}
