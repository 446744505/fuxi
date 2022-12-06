package providee

import "fuxi/core"

type ProvideeServiceConf interface {
	core.ServiceConf
	PVID() int16
}

type ProvideeServiceConfProp struct {
	pvid int16
}

func (self *ProvideeServiceConfProp) PVID() int16 {
	return self.pvid
}
