package util

import "fuxi/core"

type DispatchToFunc func(dispatch *core.Dispatch)
type SendToProvideeFunc func(pvid int32, msg core.Msg)

var (
	DispatchToClient   DispatchToFunc
	DispatchToProvidee DispatchToFunc
	SendToProvidee SendToProvideeFunc
)
