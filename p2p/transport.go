package p2p

import "net"

// RPC represents a remote procedure call between peers
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}

// Peer represents a remote node in the network
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport handles communication between nodes in the network
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}