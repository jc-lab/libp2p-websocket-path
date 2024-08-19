package websocket_path

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	addr, err := NewMultiaddr("/dns4/example.com/tcp/443/tls/ws+path/(/ipfs/ws)/p2p/12D3KooWBe6Lj7WALN1e3mLzNntnizqeDq3CZ3JvEUc1dXTXqegE")
	if err != nil {
		log.Fatalln(errors.Wrap(err, "a"))
	}
	a1, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "b"))
	}
	_ = a1
}
