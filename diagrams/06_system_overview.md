# Drift - Complete System Overview

## System Architecture Overview
This diagram provides a comprehensive view of the entire Drift distributed file system, showing how all components interact and work together to provide a robust, secure, and scalable file storage solution.

## Complete System Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        Client[Client Application]
        API[File Operations API]
        Client --> API
    end
    
    subgraph "Drift Network"
        subgraph "Node 1 (:3000) - Bootstrap"
            subgraph "Core Components 1"
                FS1[File Server]
                Store1[Store]
                Crypto1[Crypto]
                TCP1[TCP Transport]
            end
            
            subgraph "Storage 1"
                CAS1[Content-Addressable Storage]
                Disk1[Local Disk]
                CAS1 --> Disk1
            end
            
            FS1 --> Store1
            FS1 --> Crypto1
            FS1 --> TCP1
            Store1 --> CAS1
        end
        
        subgraph "Node 2 (:7000) - Bootstrap"
            subgraph "Core Components 2"
                FS2[File Server]
                Store2[Store]
                Crypto2[Crypto]
                TCP2[TCP Transport]
            end
            
            subgraph "Storage 2"
                CAS2[Content-Addressable Storage]
                Disk2[Local Disk]
                CAS2 --> Disk2
            end
            
            FS2 --> Store2
            FS2 --> Crypto2
            FS2 --> TCP2
            Store2 --> CAS2
        end
        
        subgraph "Node 3 (:5000) - Regular"
            subgraph "Core Components 3"
                FS3[File Server]
                Store3[Store]
                Crypto3[Crypto]
                TCP3[TCP Transport]
            end
            
            subgraph "Storage 3"
                CAS3[Content-Addressable Storage]
                Disk3[Local Disk]
                CAS3 --> Disk3
            end
            
            FS3 --> Store3
            FS3 --> Crypto3
            FS3 --> TCP3
            Store3 --> CAS3
        end
        
        subgraph "P2P Network Mesh"
            TCP1 <--> TCP2
            TCP1 <--> TCP3
            TCP2 <--> TCP3
        end
    end
    
    subgraph "Operations Flow"
        StoreFlow[Store Operation]
        GetFlow[Get Operation]
        DeleteFlow[Delete Operation]
        
        StoreFlow --> Encrypt[Encrypt File]
        StoreFlow --> Hash[Generate Hash]
        StoreFlow --> Replicate[Replicate to Peers]
        
        GetFlow --> CheckLocal[Check Local Storage]
        GetFlow --> FetchRemote[Fetch from Network]
        GetFlow --> Decrypt[Decrypt File]
        
        DeleteFlow --> RemoveLocal[Remove Local Copy]
        DeleteFlow --> NotifyPeers[Notify Peers]
        DeleteFlow --> Cleanup[Cleanup Storage]
    end
    
    API --> StoreFlow
    API --> GetFlow
    API --> DeleteFlow
    
    StoreFlow --> FS1
    StoreFlow --> FS2
    StoreFlow --> FS3
    
    GetFlow --> FS1
    GetFlow --> FS2
    GetFlow --> FS3
    
    DeleteFlow --> FS1
    DeleteFlow --> FS2
    DeleteFlow --> FS3
    
    classDef clientLayer fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef bootstrapNode fill:#ffebee,stroke:#c62828,stroke-width:2px
    classDef regularNode fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef networkLayer fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef operationsLayer fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef storageLayer fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    
    class Client,API clientLayer
    class FS1,Store1,Crypto1,TCP1,CAS1,Disk1,FS2,Store2,Crypto2,TCP2,CAS2,Disk2 bootstrapNode
    class FS3,Store3,Crypto3,TCP3,CAS3,Disk3 regularNode
    class StoreFlow,GetFlow,DeleteFlow,Encrypt,Hash,Replicate,CheckLocal,FetchRemote,Decrypt,RemoveLocal,NotifyPeers,Cleanup operationsLayer
```

## End-to-End Operation Flow

```mermaid
sequenceDiagram
    participant Client
    participant Node1
    participant Node2
    participant Node3
    participant Storage1
    participant Storage2
    participant Storage3
    
    Note over Client,Storage3: Complete File Store Operation
    Client->>Node1: Store("document.pdf", fileData)
    
    Note over Node1: Local Processing
    Node1->>Node1: Generate SHA1 hash
    Node1->>Node1: Create directory structure
    Node1->>Node1: Encrypt file data
    Node1->>Storage1: Write encrypted file
    Storage1-->>Node1: Confirm local storage
    
    Note over Node1,Node3: Network Replication
    Node1->>Node2: MessageStoreFile + encrypted data
    Node1->>Node3: MessageStoreFile + encrypted data
    
    Node2->>Storage2: Write encrypted file
    Node3->>Storage3: Write encrypted file
    
    Storage2-->>Node2: Confirm storage
    Storage3-->>Node3: Confirm storage
    
    Node2-->>Node1: Replication complete
    Node3-->>Node1: Replication complete
    
    Node1-->>Client: Store operation successful
    
    Note over Client,Storage3: Complete File Retrieve Operation
    Client->>Node2: Get("document.pdf")
    
    Node2->>Storage2: Check local storage
    
    alt File exists locally
        Storage2-->>Node2: Return encrypted file
        Node2->>Node2: Decrypt file data
        Node2-->>Client: Return file data
    else File not found locally
        Node2->>Node1: Request file
        Node1->>Storage1: Read encrypted file
        Storage1-->>Node1: Return encrypted file
        Node1-->>Node2: Send encrypted file
        Node2->>Node2: Decrypt file data
        Node2->>Storage2: Cache decrypted file
        Node2-->>Client: Return file data
    end
```

## Security and Encryption Flow

```mermaid
flowchart TD
    subgraph "Security Architecture"
        subgraph "File Input Security"
            A[Original File]
            B[File Integrity Check]
            C[Content Validation]
        end
        
        subgraph "Encryption Layer"
            D[Generate Unique Key per Node]
            E[AES-256 CTR Mode]
            F[Random IV Generation]
            G[Encrypted File Data]
        end
        
        subgraph "Network Security"
            H[Peer Authentication]
            I[Secure Transport]
            J[Message Integrity]
            K[Handshake Verification]
        end
        
        subgraph "Storage Security"
            L[Hash-based Addressing]
            M[Encrypted Storage]
            N[Access Control]
            O[Secure Deletion]
        end
        
        A --> B
        B --> C
        C --> D
        D --> E
        E --> F
        F --> G
        
        G --> H
        H --> I
        I --> J
        J --> K
        
        K --> L
        L --> M
        M --> N
        N --> O
    end
    
    classDef input fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef encryption fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef network fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef storage fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A,B,C input
    class D,E,F,G encryption
    class H,I,J,K network
    class L,M,N,O storage
```

## Fault Tolerance and Recovery

```mermaid
graph TB
    subgraph "Fault Tolerance Architecture"
        subgraph "Detection Layer"
            A[Health Monitoring]
            B[Connection Monitoring]
            C[Performance Monitoring]
            D[Network Monitoring]
        end
        
        subgraph "Recovery Layer"
            E[Automatic Reconnection]
            F[Peer Replacement]
            G[Data Replication]
            H[Load Redistribution]
        end
        
        subgraph "Consistency Layer"
            I[Conflict Resolution]
            J[State Synchronization]
            K[Version Control]
            L[Consensus Mechanism]
        end
        
        subgraph "Resilience Layer"
            M[Graceful Degradation]
            N[Backup Strategies]
            O[Disaster Recovery]
            P[Network Partitioning]
        end
        
        A --> E
        B --> F
        C --> G
        D --> H
        
        E --> I
        F --> J
        G --> K
        H --> L
        
        I --> M
        J --> N
        K --> O
        L --> P
    end
    
    classDef detection fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef recovery fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef consistency fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef resilience fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B,C,D detection
    class E,F,G,H recovery
    class I,J,K,L consistency
    class M,N,O,P resilience
```

## Performance and Scalability

```mermaid
graph TB
    subgraph "Performance Optimization"
        subgraph "Network Performance"
            A[Concurrent Connections]
            B[Streaming I/O]
            C[Connection Pooling]
            D[Message Batching]
        end
        
        subgraph "Storage Performance"
            E[Content Deduplication]
            F[Efficient Path Structure]
            G[Parallel I/O Operations]
            H[Caching Strategy]
        end
        
        subgraph "Scalability Features"
            I[Horizontal Scaling]
            J[Load Distribution]
            K[Dynamic Peer Discovery]
            L[Elastic Resource Usage]
        end
        
        subgraph "Monitoring & Analytics"
            M[Performance Metrics]
            N[Resource Usage Tracking]
            O[Bottleneck Detection]
            P[Optimization Recommendations]
        end
        
        A --> E
        B --> F
        C --> G
        D --> H
        
        E --> I
        F --> J
        G --> K
        H --> L
        
        I --> M
        J --> N
        K --> O
        L --> P
    end
    
    classDef network fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef storage fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef scalability fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef monitoring fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B,C,D network
    class E,F,G,H storage
    class I,J,K,L scalability
    class M,N,O,P monitoring
```

## System Integration Points

```mermaid
graph TB
    subgraph "External Integration"
        subgraph "API Layer"
            A[REST API]
            B[CLI Interface]
            C[SDK/Library]
            D[Web Interface]
        end
        
        subgraph "Storage Backends"
            E[Local File System]
            F[Cloud Storage]
            G[Network Attached Storage]
            H[Distributed Storage]
        end
        
        subgraph "Monitoring & Logging"
            I[Metrics Collection]
            J[Log Aggregation]
            K[Alert System]
            L[Dashboard]
        end
        
        subgraph "Security Integration"
            M[Authentication]
            N[Authorization]
            O[Audit Logging]
            P[Compliance]
        end
        
        A --> E
        B --> F
        C --> G
        D --> H
        
        E --> I
        F --> J
        G --> K
        H --> L
        
        I --> M
        J --> N
        K --> O
        L --> P
    end
    
    classDef api fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef storage fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef monitoring fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef security fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B,C,D api
    class E,F,G,H storage
    class I,J,K,L monitoring
    class M,N,O,P security
```

## Key System Characteristics

### 🔒 Security Features
- **End-to-End Encryption**: AES-256 encryption for all stored files
- **Secure Transport**: Encrypted communication between nodes
- **Peer Authentication**: Verified peer connections
- **Data Integrity**: SHA1 hashing for file verification

### 🌐 Network Architecture
- **P2P Mesh Network**: Direct peer-to-peer communication
- **Bootstrap Nodes**: Simplified network joining
- **Fault Tolerance**: Continues operation despite node failures
- **Dynamic Discovery**: Automatic peer discovery and management

### 💾 Storage System
- **Content-Addressable**: Hash-based file addressing
- **Deduplication**: Automatic removal of duplicate files
- **Distributed Replication**: Files stored across multiple nodes
- **Efficient Organization**: Optimized directory structure

### ⚡ Performance Optimizations
- **Concurrent Operations**: Parallel file processing
- **Streaming I/O**: Efficient handling of large files
- **Connection Pooling**: Reuse of network connections
- **Caching**: Intelligent file caching strategies

### 📈 Scalability
- **Horizontal Scaling**: Easy addition of new nodes
- **Load Distribution**: Balanced file distribution
- **Resource Efficiency**: Optimal resource utilization
- **Dynamic Adaptation**: Automatic scaling based on demand