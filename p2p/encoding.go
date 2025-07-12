package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
)

// Decoder interface for decoding messages
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

// GOBDecoder implements the Decoder interface using GOB encoding
type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}

// DefaultDecoder is the default decoder implementation
type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return err
	}

	// In case of a stream we are not decoding what is being sent over the network
	// We are just setting stream true so we can handle that in our logic.
	stream := peekBuf[0] == IncomingStream
	if stream {
		rpc.Stream = true
		return nil
	}

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	rpc.Payload = buf[:n]

	return nil
}

// Encoder interface for encoding messages
type Encoder interface {
	Encode(io.Writer, *RPC) error
}

// GOBEncoder implements the Encoder interface using GOB encoding
type GOBEncoder struct{}

func (enc GOBEncoder) Encode(w io.Writer, rpc *RPC) error {
	return gob.NewEncoder(w).Encode(rpc)
}

// DefaultEncoder is the default encoder implementation
type DefaultEncoder struct{}

func (enc DefaultEncoder) Encode(w io.Writer, rpc *RPC) error {
	if rpc.Stream {
		// For streams, we just write the stream indicator
		_, err := w.Write([]byte{IncomingStream})
		return err
	}

	// For regular messages, write the message indicator followed by the payload
	if _, err := w.Write([]byte{IncomingMessage}); err != nil {
		return err
	}

	if len(rpc.Payload) == 0 {
		return fmt.Errorf("empty payload")
	}

	_, err := w.Write(rpc.Payload)
	return err
}