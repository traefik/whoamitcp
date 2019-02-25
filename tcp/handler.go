package tcp

import "net"

type Handler struct {
	ServeTCP func(conn net.Conn)
}
