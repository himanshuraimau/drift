# Drift - Distributed File System

A peer-to-peer distributed file system built in Go with content-addressable storage, AES-256 encryption, and fault tolerance.

## Features

- **Distributed Storage**: Files replicated across multiple nodes
- **Content-Addressable**: SHA1-based file addressing and deduplication  
- **Encryption**: AES-256 encryption for all stored files
- **P2P Network**: Direct peer-to-peer communication
- **Fault Tolerance**: Continues operating if nodes fail

## Architecture

View the [system architecture diagram](diagrams/01_overall_architecture.md) for a complete overview of the distributed file system components and their interactions.

```text
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Node 1   │────│    Node 2   │────│    Node 3   │
│   :3000     │    │   :7000     │    │   :5000     │
└─────────────┘    └─────────────┘    └─────────────┘
       │                  │                  │
       └──────────────────┼──────────────────┘
                          │
                   P2P Mesh Network
```

## Quick Start

### Prerequisites

- Go 1.24 or later

### Build and Run

```bash
# Clone the repository
git clone https://github.com/himanshuraimau/drift.git
cd drift

# Install dependencies
go mod download

# Build the project
make build

# Run the demo (starts 3 nodes automatically)
make run
```

## Usage

The demo starts three interconnected nodes:

- Node 1: `:3000` (Bootstrap)
- Node 2: `:7000` (Bootstrap)  
- Node 3: `:5000` (Connects to both)

### Manual Usage

```go
// Create a new node
server := makeServer(":8000", ":3000", ":7000")
go server.Start()

// Store a file (replicated across network)
err := server.Store("myfile.txt", bytes.NewReader([]byte("content")))

// Retrieve a file (from local or network)
reader, err := server.Get("myfile.txt")

// Delete a file (removed from all nodes)
err := server.Delete("myfile.txt")
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race
```

## Build Commands

```bash
make build          # Build binary
make test           # Run tests
make clean          # Clean artifacts
make fmt            # Format code
make vet            # Vet code
```

## How It Works

1. **File Storage**: Files are encrypted with AES-256 and stored using SHA1-based paths
2. **Network Replication**: Each file operation is replicated across all connected peers
3. **Content Addressing**: Files are deduplicated using content-based addressing
4. **Peer Discovery**: New nodes can join by connecting to existing bootstrap nodes
5. **Fault Tolerance**: If a file isn't found locally, it's retrieved from the network

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new features
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
