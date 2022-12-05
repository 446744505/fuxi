package providee

import (
	msg "fuxi/gen"
	"fuxi/net"
)

type ProvideeEventHandler struct {
	net.SampleEventHandler
}

func (self ProvideeEventHandler) OnSessionAdd(session *net.Session) {
	session.Send(&msg.BindPvid{PVID: })
}

func (self ProvideeEventHandler) OnSessionRemoved(session *net.Session) {

}