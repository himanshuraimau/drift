# Drift - System Architecture

## Overview

Drift is a peer-to-peer distributed file system with encryption and content-addressable storage.

```mermaid
graph TB
    subgraph "Node 1 (:3000)"
        FS1[File Server]
        Store1[Store]
        Crypto1[Crypto]
        TCP1[TCP Transport]
        
        FS1 --> Store1
        FS1 --> Crypto1
        FS1 --> TCP1
    end
    
    subgraph "Node 2 (:7000)"
        FS2[File Server]
        Store2[Store]
        Crypto2[Crypto]
        TCP2[TCP Transport]
        
        FS2 --> Store2
        FS2 --> Crypto2
        FS2 --> TCP2
    end
    
    subgraph "Node 3 (:5000)"
        FS3[File Server]
        Store3[Store]
        Crypto3[Crypto]
        TCP3[TCP Transport]
        
        FS3 --> Store3
        FS3 --> Crypto3
        FS3 --> TCP3
    end
    
    TCP1 <--> TCP2
    TCP1 <--> TCP3
    TCP2 <--> TCP3
    
    Store1 --> Files1[Encrypted Files]
    Store2 --> Files2[Encrypted Files]
    Store3 --> Files3[Encrypted Files]
    
    classDef node fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef storage fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class FS1,FS2,FS3,Store1,Store2,Store3,Crypto1,Crypto2,Crypto3,TCP1,TCP2,TCP3 node
    class Files1,Files2,Files3 storage
```

## Core Components

- **File Server**: Orchestrates file operations and peer management
- **Store**: Content-addressable storage with SHA1-based paths
- **Crypto**: AES-256 encryption for secure file storage
- **TCP Transport**: P2P network communication layer

## Key Features

- **Distributed Storage**: Files replicated across multiple nodes
- **Content-Addressable**: SHA1-based file addressing and deduplication
- **Encryption**: AES-256 encryption for all stored files
- **P2P Network**: Direct peer-to-peer communication
- **Fault Tolerance**: Continues operating if nodes fail
