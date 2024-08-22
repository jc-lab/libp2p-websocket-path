package websocket_path

import (
	"testing"
)

func Test_parseWebsocketMultiaddr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(p parsedWebsocketMultiaddr) bool
		wantErr error
	}{
		{
			"with-tls",
			"/dns4/example.com/tcp/443/tls/ws+path/(/ipfs/ws)/p2p/12D3KooWBe6Lj7WALN1e3mLzNntnizqeDq3CZ3JvEUc1dXTXqegE",
			func(p parsedWebsocketMultiaddr) bool {
				if p.path.String() != "/ws+path//ipfs/ws" {
					return false
				}
				if p.isWSS != true || p.sni != nil {
					return false
				}
				if p.restMultiaddr.String() != "/dns4/example.com/tcp/443" {
					return false
				}
				return true
			},
			nil,
		},
		{
			"without-tls",
			"/dns4/example.com/tcp/80/ws+path/(/ipfs/ws)/p2p/12D3KooWBe6Lj7WALN1e3mLzNntnizqeDq3CZ3JvEUc1dXTXqegE",
			func(p parsedWebsocketMultiaddr) bool {
				if p.path.String() != "/ws+path//ipfs/ws" {
					return false
				}
				if p.isWSS != false || p.sni != nil {
					return false
				}
				if p.restMultiaddr.String() != "/dns4/example.com/tcp/80" {
					return false
				}
				return true
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := NewMultiaddr(tt.input)
			if err != nil {
				t.Error(err)
				return
			}
			got, err := parseWebsocketMultiaddr(addr)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("parseWebsocketMultiaddr() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !tt.check(got) {
				t.Errorf("parseWebsocketMultiaddr() got = %v", got)
			}
		})
	}
}
