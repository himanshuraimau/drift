# Drift - File Storage Architecture

## Storage System Overview
This diagram shows the detailed architecture of Drift's content-addressable storage system, including path transformation, encryption, and file organization.

## Content-Addressable Storage (CAS) Architecture

```mermaid
graph TB
    subgraph "File Input"
        A[File Key: "picture.jpg"]
        B[File Data: Binary Content]
    end
    
    subgraph "Hash Generation"
        C[SHA1 Hash Generator]
        D[Generated Hash:<br/>3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b]
    end
    
    subgraph "Path Transformation"
        E[CASPathTransformFunc]
        F[Split hash into blocks of 5]
        G[Block 1: 3a4b5]
        H[Block 2: c6d7e]
        I[Block 3: 8f9a0]
        J[Block 4: b1c2d]
        K[Block 5: 3e4f5]
        L[Block 6: a6b7c]
        M[Block 7: 8d9e0]
        N[Block 8: f1a2b]
    end
    
    subgraph "Directory Structure"
        O[Root Directory]
        P[Node ID Directory]
        Q[3a4b5/]
        R[c6d7e/]
        S[8f9a0/]
        T[b1c2d/]
        U[3e4f5/]
        V[a6b7c/]
        W[8d9e0/]
        X[File: 3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b]
    end
    
    A --> C
    B --> C
    C --> D
    D --> E
    E --> F
    F --> G
    F --> H
    F --> I
    F --> J
    F --> K
    F --> L
    F --> M
    F --> N
    
    G --> O
    H --> O
    I --> O
    J --> O
    K --> O
    L --> O
    M --> O
    N --> O
    
    O --> P
    P --> Q
    Q --> R
    R --> S
    S --> T
    T --> U
    U --> V
    V --> W
    W --> X
    
    classDef input fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef hash fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef transform fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef storage fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B input
    class C,D hash
    class E,F,G,H,I,J,K,L,M,N transform
    class O,P,Q,R,S,T,U,V,W,X storage
```

## File Storage Hierarchy

```mermaid
graph TB
    subgraph "Storage Root Structure"
        Root[drift_root/]
        
        subgraph "Node-specific Storage"
            Node1[node1_id/]
            Node2[node2_id/]
            Node3[node3_id/]
        end
        
        subgraph "Example: Node1 File Storage"
            Dir1[3a4b5/]
            Dir2[c6d7e/]
            Dir3[8f9a0/]
            Dir4[b1c2d/]
            File1[3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b]
        end
        
        subgraph "File Metadata"
            Size[File Size: 1024 bytes]
            Encrypted[Encrypted: Yes]
            Created[Created: timestamp]
            Hash[Hash: SHA1]
        end
        
        Root --> Node1
        Root --> Node2
        Root --> Node3
        
        Node1 --> Dir1
        Dir1 --> Dir2
        Dir2 --> Dir3
        Dir3 --> Dir4
        Dir4 --> File1
        
        File1 --> Size
        File1 --> Encrypted
        File1 --> Created
        File1 --> Hash
    end
    
    classDef rootDir fill:#ffebee,stroke:#c62828,stroke-width:2px
    classDef nodeDir fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef hashDir fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef fileNode fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef metadata fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    
    class Root rootDir
    class Node1,Node2,Node3 nodeDir
    class Dir1,Dir2,Dir3,Dir4 hashDir
    class File1 fileNode
    class Size,Encrypted,Created,Hash metadata
```

## Storage Operations Flow

```mermaid
sequenceDiagram
    participant Client
    participant Store
    participant PathTransform
    participant FileSystem
    participant Crypto
    
    Note over Client,Crypto: File Store Operation
    Client->>Store: Write(nodeID, "file.jpg", data)
    Store->>PathTransform: CASPathTransformFunc("file.jpg")
    PathTransform->>PathTransform: Generate SHA1 hash
    PathTransform->>PathTransform: Create directory path
    PathTransform-->>Store: Return PathKey
    
    Store->>FileSystem: Create directory structure
    FileSystem->>FileSystem: mkdir -p /root/nodeID/3a4b5/c6d7e/...
    FileSystem-->>Store: Directories created
    
    Store->>Crypto: Encrypt file data
    Crypto->>Crypto: Generate IV + AES encrypt
    Crypto-->>Store: Encrypted data
    
    Store->>FileSystem: Write encrypted data to file
    FileSystem-->>Store: File written successfully
    Store-->>Client: Return file size
    
    Note over Client,Crypto: File Read Operation
    Client->>Store: Read(nodeID, "file.jpg")
    Store->>PathTransform: CASPathTransformFunc("file.jpg")
    PathTransform-->>Store: Return PathKey
    
    Store->>FileSystem: Read file from path
    FileSystem-->>Store: Return encrypted data
    
    Store->>Crypto: Decrypt file data
    Crypto->>Crypto: Extract IV + AES decrypt
    Crypto-->>Store: Decrypted data
    
    Store-->>Client: Return file size and reader
```

## Encryption Layer Integration

```mermaid
graph TB
    subgraph "Storage with Encryption"
        subgraph "Input Layer"
            A[Original File Data]
            B[File Key]
            C[Node ID]
        end
        
        subgraph "Path Layer"
            D[Path Transform]
            E[Hash Generation]
            F[Directory Structure]
        end
        
        subgraph "Encryption Layer"
            G[Generate IV]
            H[AES-256 CTR]
            I[Encryption Key]
            J[Encrypted Data]
        end
        
        subgraph "Storage Layer"
            K[File System]
            L[Disk Storage]
            M[Stored File:<br/>IV + Encrypted Data]
        end
        
        A --> D
        B --> D
        C --> D
        D --> E
        E --> F
        
        A --> G
        G --> H
        I --> H
        H --> J
        
        J --> K
        F --> K
        K --> L
        L --> M
    end
    
    classDef inputLayer fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef pathLayer fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef cryptoLayer fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef storageLayer fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B,C inputLayer
    class D,E,F pathLayer
    class G,H,I,J cryptoLayer
    class K,L,M storageLayer
```

## Storage Deduplication

```mermaid
flowchart TD
    subgraph "File Deduplication Process"
        A[File 1: "hello.txt"]
        B[File 2: "greeting.txt"]
        C[File 3: "hello.txt"]
        
        A --> D[SHA1: abc123...]
        B --> E[SHA1: def456...]  
        C --> F[SHA1: abc123...]
        
        D --> G{Hash Exists?}
        E --> H{Hash Exists?}
        F --> I{Hash Exists?}
        
        G -->|No| J[Store New File]
        H -->|No| K[Store New File]
        I -->|Yes| L[Reference Existing File]
        
        J --> M[Disk Space: +1 file]
        K --> N[Disk Space: +1 file]
        L --> O[Disk Space: +0 files]
        
        subgraph "Storage Efficiency"
            P[3 Files Uploaded]
            Q[2 Files Stored]
            R[33% Space Saved]
        end
        
        M --> P
        N --> P
        O --> P
        P --> Q
        Q --> R
    end
    
    classDef fileInput fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef hashProcess fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef decision fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef storage fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef efficiency fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B,C fileInput
    class D,E,F hashProcess
    class G,H,I decision
    class J,K,L,M,N,O storage
    class P,Q,R efficiency
```

## Storage Performance Optimization

```mermaid
graph TB
    subgraph "Performance Optimizations"
        subgraph "Directory Structure"
            A[Shallow Directory Trees]
            B[Balanced Distribution]
            C[Fast File Lookups]
        end
        
        subgraph "I/O Optimization"
            D[Streaming I/O]
            E[Buffered Writes]
            F[Concurrent Operations]
        end
        
        subgraph "Caching Strategy"
            G[Recently Accessed Files]
            H[Metadata Cache]
            I[Path Cache]
        end
        
        subgraph "Monitoring"
            J[Storage Usage]
            K[Access Patterns]
            L[Performance Metrics]
        end
        
        A --> D
        B --> E
        C --> F
        
        D --> G
        E --> H
        F --> I
        
        G --> J
        H --> K
        I --> L
    end
    
    classDef structure fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef optimization fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef caching fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef monitoring fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B,C structure
    class D,E,F optimization
    class G,H,I caching
    class J,K,L monitoring
```

## Storage Consistency Model

```mermaid
sequenceDiagram
    participant Node1
    participant Node2
    participant Node3
    participant Client
    
    Note over Node1,Node3: File Store Consistency
    Client->>Node1: Store "file.txt"
    Node1->>Node1: Store locally
    Node1->>Node2: Replicate file
    Node1->>Node3: Replicate file
    
    Node2->>Node2: Store locally
    Node3->>Node3: Store locally
    
    Node2-->>Node1: Store complete
    Node3-->>Node1: Store complete
    
    Node1-->>Client: Store operation complete
    
    Note over Node1,Node3: File Get Consistency
    Client->>Node2: Get "file.txt"
    Node2->>Node2: Check local storage
    
    alt File exists locally
        Node2-->>Client: Return file data
    else File not found locally
        Node2->>Node1: Request file
        Node1-->>Node2: Send file data
        Node2->>Node2: Store locally
        Node2-->>Client: Return file data
    end
```

## Key Storage Features

1. **Content-Addressable**: Files stored using SHA1 hash-based addressing
2. **Deduplication**: Identical files stored only once across the network
3. **Encrypted Storage**: All files encrypted with AES-256 before storage
4. **Distributed**: Files replicated across multiple nodes for reliability
5. **Efficient Organization**: Hash-based directory structure for fast access
6. **Metadata Tracking**: File size, encryption status, and timestamps
7. **Fault Tolerance**: Files remain accessible even if some nodes fail
8. **Space Efficiency**: Automatic deduplication reduces storage requirements