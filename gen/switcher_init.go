package msg

import "fuxi/net"

func init() {
	net.RegisterMsg(&Dispatch{})
	net.RegisterMsg(&PDispatch{})
}