package util

var CtxTypeSessionInfo = "session_info"

type LinkerSessionInfo struct {
	RoleId int64
}

type ProvideeSessionInfo struct {
	Pvid int32
	Name string
}