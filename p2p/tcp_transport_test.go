package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	opts := TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	
	assert.Equal(t, listenAddr, tr.Addr(), "Transport address should match")
	assert.NotNil(t, tr.rpcch, "RPC channel should be initialized")
}
