package websocket_path

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

// Addr is an implementation of net.Addr for WebSocket.
type Addr struct {
	*url.URL
}

var _ net.Addr = Addr{}

// Network returns the network type for a WebSocket, "websocket".
func (a Addr) Network() string {
	return "websocket+path"
}

// NewLocalAddrWithScheme creates a new Addr using the given host string. isSecure
// should be true for WSS connections and false for WS.
func NewLocalAddrWithScheme(host string, isSecure bool) Addr {
	scheme := "ws"
	if isSecure {
		scheme = "wss"
	}
	return Addr{
		URL: &url.URL{
			Scheme: scheme,
			Host:   host,
		},
	}
}

func ConvertWebsocketMultiaddrToNetAddr(maddr ma.Multiaddr) (net.Addr, error) {
	url, err := parseMultiaddr(maddr)
	if err != nil {
		return nil, err
	}
	return Addr{URL: url}, nil
}

func ParseWebsocketNetAddr(a net.Addr) (ma.Multiaddr, error) {
	wsa, ok := a.(Addr)
	if !ok {
		return nil, fmt.Errorf("not a websocket address")
	}

	var (
		tcpma ma.Multiaddr
		err   error
		port  int
		host  = wsa.Hostname()
	)

	// Get the port
	if portStr := wsa.Port(); portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse port '%q': %s", portStr, err)
		}
	} else {
		return nil, fmt.Errorf("invalid port in url: '%q'", wsa.URL)
	}

	// NOTE: Ignoring IPv6 zones...
	// Detect if host is IP address or DNS
	if ip := net.ParseIP(host); ip != nil {
		// Assume IP address
		tcpma, err = manet.FromNetAddr(&net.TCPAddr{
			IP:   ip,
			Port: port,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Assume DNS name
		tcpma, err = ma.NewMultiaddr(fmt.Sprintf("/dns/%s/tcp/%d", host, port))
		if err != nil {
			return nil, err
		}
	}

	wsma, err := ma.NewMultiaddr("/" + wsa.Scheme)
	if err != nil {
		return nil, err
	}

	return tcpma.Encapsulate(wsma), nil
}

func parseMultiaddr(maddr ma.Multiaddr) (*url.URL, error) {
	parsed, err := parseWebsocketMultiaddr(maddr)
	if err != nil {
		return nil, err
	}

	scheme := "ws"
	if parsed.isWSS {
		scheme = "wss"
	}

	path, err := maddr.ValueForProtocol(P_WS_PATH)
	if err != nil {
		return nil, err
	}

	network, host, err := manet.DialArgs(parsed.restMultiaddr)
	if err != nil {
		return nil, err
	}
	switch network {
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, fmt.Errorf("unsupported websocket network %s", network)
	}
	return &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}, nil
}

type parsedWebsocketMultiaddr struct {
	isWSS bool
	// sni is the SNI value for the TLS handshake, and for setting HTTP Host header
	sni *ma.Component
	// the rest of the multiaddr before the /tls/sni/example.com/ws or /ws or /wss
	restMultiaddr ma.Multiaddr
	path          *ma.Component
}

func parseWebsocketMultiaddr(a ma.Multiaddr) (parsedWebsocketMultiaddr, error) {
	out := parsedWebsocketMultiaddr{}

	rest := a
	for {
		var head *ma.Component
		rest, head = ma.SplitLast(rest)
		if head == nil || rest == nil {
			break
		}

		if head.Protocol().Code == P_WS_PATH {
			out.path = head
			out.restMultiaddr = rest
		} else if head.Protocol().Code == ma.P_SNI {
			out.sni = head
			out.restMultiaddr = rest
		} else if head.Protocol().Code == ma.P_TLS {
			out.isWSS = true
			out.restMultiaddr = rest
			break
		}
	}

	if out.path == nil {
		return out, fmt.Errorf("not a websocket multiaddr")
	}

	return out, nil
}

func (pwma *parsedWebsocketMultiaddr) toMultiaddr() ma.Multiaddr {
	if !pwma.isWSS {
		return pwma.restMultiaddr.Encapsulate(pwma.path)
	}

	if pwma.sni == nil {
		return pwma.restMultiaddr.Encapsulate(tlsComponent).Encapsulate(pwma.path)
	}

	return pwma.restMultiaddr.Encapsulate(tlsComponent).Encapsulate(pwma.sni).Encapsulate(pwma.path)
}
