package providee

import "fuxi/core"

type ProvideeServiceConf interface {
	core.ServiceConf
	PVID() int32
}

type ProvideeServiceConfProp struct {
	pvid int32
}

func (self *ProvideeServiceConfProp) PVID() int32 {
	return self.pvid
}
