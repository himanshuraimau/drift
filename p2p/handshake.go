package p2p

import (
	"errors"
	"fmt"
)

// ErrInvalidHandshake is returned when the handshake fails
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakeFunc defines the handshake function signature
type HandshakeFunc func(Peer) error

// NOPHandshakeFunc is a no-operation handshake function
func NOPHandshakeFunc(Peer) error {
	return nil
}

// BasicHandshakeFunc performs a basic handshake with version checking
func BasicHandshakeFunc(peer Peer) error {
	// Send handshake message
	handshakeMsg := []byte("DRIFT-FS-v1.0")
	if err := peer.Send(handshakeMsg); err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}

	// Read response
	buf := make([]byte, len(handshakeMsg))
	if _, err := peer.Read(buf); err != nil {
		return fmt.Errorf("failed to read handshake response: %w", err)
	}

	// Verify handshake
	if string(buf) != string(handshakeMsg) {
		return ErrInvalidHandshake
	}

	return nil
}