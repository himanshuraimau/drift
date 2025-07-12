package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// Message represents a message sent between peers
type Message struct {
	Header MessageHeader
	Data   []byte
}

// MessageHeader contains metadata about the message
type MessageHeader struct {
	Type     byte
	Length   uint32
	Checksum uint32
}

// MessageType represents different types of messages
type MessageType byte

const (
	MessageTypeHandshake MessageType = iota
	MessageTypeStore
	MessageTypeGet
	MessageTypeDelete
	MessageTypeHeartbeat
)