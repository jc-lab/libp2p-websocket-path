# go-libp2p-websocket-path

libp2p websocket transport with path

**Transport:**

```go
libp2p.Transport(path_websocket.New)
```

**Connect:**

```go
addr, err := websocket_path.NewMultiaddr("/dns4/your.ipfs.server.io/tcp/443/tls/ws+path/(/ipfs/ws)/p2p/12D3KooXXXXXXXXXX")
info, err := peer.AddrInfoFromP2pAddr(addr)
lite.Bootstrap([]peer.AddrInfo{*info})
```

# License

[MIT License](./LICENSE)
