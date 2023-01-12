package util

import "fuxi/core"

type DispatchToFunc func(dispatch *core.Dispatch)
type SendToProvideeFunc func(pvid int32, msg core.Msg)
type ClientBrokenFunc func(clientSid int64)

var (
	DispatchToClient   DispatchToFunc
	ClientToProvidee DispatchToFunc
	ProvideeToProvidee DispatchToFunc
	SendToProvidee     SendToProvideeFunc
	ClientBroken ClientBrokenFunc
)
