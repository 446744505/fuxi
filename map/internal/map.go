package internal

import (
	"fuxi/core"
	"fuxi/providee"
	"github.com/davyxu/golog"
	"strconv"
)

var Map *mmap
var Log = golog.New("map")

type mmap struct {
	core.NetControlImpl

	Pvid int32
}

func NewMap() *mmap {
	Map = &mmap{}
	pvid, _ := strconv.Atoi(core.Args.Get("pvid"))
	Map.Pvid = int32(pvid)
	p := providee.NewProvidee(int32(pvid), core.ServerMap)
	p.SetEventHandler(&mapEventHandler{})
	Map.AddService(p)
	return Map
}

func (self *mmap) SendToClient(providerName string, clientSid int64, msg core.Msg) bool {
	if provider := providee.Providee.GetProvider(providerName); provider == nil {
		return false
	} else {
		msg.SetToID(clientSid)
		msg.SetToType(core.MsgToClient)
		provider.Send(msg)
	}
	return true
}

func (self *mmap) SendToProvidee(pvid int32, msg core.Msg) bool {
	return providee.Providee.SendToProvidee(pvid, msg)
}