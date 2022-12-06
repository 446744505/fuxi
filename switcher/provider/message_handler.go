package provider

import (
	"fuxi/core"
	msg "fuxi/gen"
)

func (self *ProviderEventHandler) OnPDispatch(p core.Msg) {

}

func (self *ProviderEventHandler) OnBindPvid(p core.Msg) {
	bind := p.(*msg.BindPvid)
	Provider.BindProvidee(bind.PVID, bind.Name, p.Session())
}

func (self *ProviderEventHandler) OnUnBindPvid(p core.Msg) {
	unbind := p.(*msg.UnBindPvid)
	Provider.UnBindProvidee(unbind.PVID)
}