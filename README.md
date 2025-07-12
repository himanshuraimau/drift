# Drift - Distributed File System

A peer-to-peer distributed file system built in Go, featuring content-addressable storage, encryption, and fault tolerance.

## 🚀 Features

- **Distributed Storage**: Files are replicated across multiple nodes for redundancy
- **Content-Addressable Storage**: Files are stored using SHA1-based addressing
- **Encryption**: All files are encrypted with AES-256 before storage
- **P2P Network**: Direct peer-to-peer communication without central server
- **Fault Tolerance**: Network can continue operating even if some nodes fail
- **Bootstrap Network**: New nodes can easily join existing networks

## 📋 Architecture

### Core Components

- **P2P Network Layer**: TCP-based transport with peer discovery and management
- **File Server**: Main orchestrator for file operations (store, get, delete)
- **Store**: Content-addressable storage with path transformation
- **Crypto**: AES encryption/decryption for secure file storage
- **Message System**: GOB-based serialization for peer communication

### Network Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Node 1   │────│    Node 2   │────│    Node 3   │
│   :3000     │    │   :7000     │    │   :5000     │
└─────────────┘    └─────────────┘    └─────────────┘
       │                  │                  │
       └──────────────────┼──────────────────┘
                          │
                   P2P Network
```

## 🛠️ Installation

### Prerequisites

- Go 1.24 or later
- Git

### Build from Source

```bash
# Clone the repository
git clone https://github.com/himanshuraimau/drift.git
cd drift

# Install dependencies
make deps

# Build the project
make build

# Run tests
make test
```

## 🚦 Usage

### Running a Single Node

```bash
# Build and run
make run
```

### Running Multiple Nodes

The system automatically starts three nodes in the demo:

- Node 1: `:3000` (Bootstrap node)
- Node 2: `:7000` (Bootstrap node)  
- Node 3: `:5000` (Connects to both bootstrap nodes)

### Manual Node Configuration

```go
// Create a new node
s := makeServer(":8000", ":3000", ":7000") // Listen on :8000, connect to :3000 and :7000

// Start the server
go s.Start()

// Store a file
err := s.Store("myfile.txt", bytes.NewReader([]byte("file content")))

// Retrieve a file
reader, err := s.Get("myfile.txt")

// Delete a file
err := s.Delete("myfile.txt")
```

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

### Run Tests with Race Detection

```bash
make test-race
```

### Run Benchmarks

```bash
make bench
```

## 📦 Available Commands

```bash
# Build commands
make build          # Build the binary
make build-prod     # Build optimized binary for production
make build-all      # Build for all platforms

# Test commands
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make test-race      # Run tests with race detection

# Run commands
make run            # Build and run the application
make run-test       # Run tests then the application

# Utility commands
make clean          # Clean build artifacts
make clean-test     # Clean test data
make deps           # Install dependencies
make update-deps    # Update dependencies
make fmt            # Format code
make vet            # Vet code
make check          # Run fmt, vet, lint, and test
```

## 🔧 Configuration

### File Server Options

```go
type FileServerOpts struct {
    ID                string              // Unique node identifier
    EncKey            []byte              // Encryption key (32 bytes)
    StorageRoot       string              // Root directory for file storage
    PathTransformFunc PathTransformFunc   // Path transformation function
    Transport         p2p.Transport       // Network transport layer
    BootstrapNodes    []string            // List of bootstrap node addresses
}
```

### Transport Options

```go
type TCPTransportOpts struct {
    ListenAddr    string              // Address to listen on
    HandshakeFunc HandshakeFunc       // Peer handshake function
    Decoder       Decoder             // Message decoder
    OnPeer        func(Peer) error    // Callback for new peer connections
}
```

## 🔐 Security

- **Encryption**: All files are encrypted with AES-256 in CTR mode
- **Key Management**: Each node generates its own encryption key
- **Peer Authentication**: Basic handshake mechanism for peer verification
- **Data Integrity**: SHA1 hashes ensure file integrity

## 📊 Performance

- **Concurrent Operations**: Supports multiple simultaneous file operations
- **Efficient Storage**: Content-addressable storage eliminates duplicates
- **Network Optimization**: Efficient binary protocols for data transfer
- **Scalability**: Designed to handle large numbers of nodes

## 🔍 Monitoring

The system provides detailed logging for:

- Node connections and disconnections
- File operations (store, get, delete)
- Network events and errors
- Performance metrics

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by distributed systems concepts from BitTorrent, IPFS, and Kademlia
- Built with Go's powerful concurrency primitives
- Thanks to the Go community for excellent libraries and tools

## 📞 Support

For support, questions, or contributions:

- Open an issue on GitHub
- Check the documentation
- Review the test files for usage examples

## 🗺️ Roadmap

- [ ] Web UI for file management
- [ ] REST API for external integrations
- [ ] Advanced peer discovery mechanisms
- [ ] Improved replication strategies
- [ ] Docker containerization
- [ ] Kubernetes deployment configs
- [ ] Metrics and monitoring dashboard
- [ ] Load balancing and sharding

---

**Built with ❤️ using Go**