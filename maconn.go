package websocket_path

import (
	ma "github.com/multiformats/go-multiaddr"
	"net"
)

type maconn struct {
	net.Conn
	localAddr  ma.Multiaddr
	remoteAddr ma.Multiaddr
}

func (m *maconn) LocalMultiaddr() ma.Multiaddr {
	return m.localAddr
}

func (m *maconn) RemoteMultiaddr() ma.Multiaddr {
	return m.remoteAddr
}
