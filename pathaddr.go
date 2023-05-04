package websocket_path

import (
	"bytes"
	"fmt"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/multiformats/go-varint"
	"strings"
)

const (
	P_WS_PATH = 0x3001dd
)

var (
	protoPath = multiaddr.Protocol{
		Name:       "ws+path",
		Code:       P_WS_PATH,
		VCode:      multiaddr.CodeToVarint(P_WS_PATH),
		Size:       multiaddr.LengthPrefixedVarSize,
		Path:       true,
		Transcoder: multiaddr.TranscoderUnix,
	}
)

func init() {
	multiaddr.AddProtocol(protoPath)
	manet.RegisterFromNetAddr(ParseWebsocketNetAddr, "websocket+path")
}

func stringToBytes(s string) ([]byte, error) {
	// consume trailing slashes
	s = strings.TrimRight(s, "/")

	var b bytes.Buffer
	sp := strings.Split(s, "/")

	if sp[0] != "" {
		return nil, fmt.Errorf("failed to parse multiaddr %q: must begin with /", s)
	}

	// consume first empty elem
	sp = sp[1:]

	if len(sp) == 0 {
		return nil, fmt.Errorf("failed to parse multiaddr %q: empty multiaddr", s)
	}

	for len(sp) > 0 {
		name := sp[0]

		var p multiaddr.Protocol
		if name == protoPath.Name {
			p = protoPath
		} else {
			p = multiaddr.ProtocolWithName(name)
		}
		if p.Code == 0 {
			return nil, fmt.Errorf("failed to parse multiaddr %q: unknown protocol %s", s, sp[0])
		}
		_, _ = b.Write(p.VCode)
		sp = sp[1:]

		if p.Size == 0 { // no length.
			continue
		}

		if len(sp) < 1 {
			return nil, fmt.Errorf("failed to parse multiaddr %q: unexpected end of multiaddr", s)
		}

		if p.Code == protoPath.Code {
			remaining := strings.Join(sp, "/")[1:]
			endPos := strings.Index(remaining, ")")
			if endPos < 0 {
				return nil, fmt.Errorf("failed to parse multiaddr %q: unexpected end of multiaddr", s)
			}
			subpath := remaining[:endPos]

			_, _ = b.Write(varint.ToUvarint(uint64(len(subpath))))
			b.Write([]byte(subpath))

			remaining = remaining[endPos+2:]
			sp = strings.Split(remaining, "/")
		} else {
			if p.Path {
				// it's a path protocol (terminal).
				// consume the rest of the address as the next component.
				sp = []string{"/" + strings.Join(sp, "/")}
			}

			a, err := p.Transcoder.StringToBytes(sp[0])
			if err != nil {
				return nil, fmt.Errorf("failed to parse multiaddr %q: invalid value %q for protocol %s: %s", s, sp[0], p.Name, err)
			}
			if p.Size < 0 { // varint size.
				_, _ = b.Write(varint.ToUvarint(uint64(len(a))))
			}
			b.Write(a)
			sp = sp[1:]
		}
	}

	return b.Bytes(), nil
}

func NewMultiaddr(address string) (multiaddr.Multiaddr, error) {
	bytes, err := stringToBytes(address)
	if err != nil {
		return nil, err
	}
	return multiaddr.NewMultiaddrBytes(bytes)
}
