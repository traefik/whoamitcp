package tcp

import (
	"crypto/tls"
	"net"
)

type TlsHandler struct {
	Next   Handler
	Config *tls.Config
}

func (t *TlsHandler) ServeTCP(conn net.Conn) {
	t.Next.ServeTCP(tls.Server(conn, t.Config))
}
