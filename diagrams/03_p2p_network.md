# Drift - P2P Network Communication

## P2P Network Overview
This diagram shows how nodes communicate in the Drift distributed file system, including peer discovery, message routing, and data transfer.

## Network Topology

```mermaid
graph TB
    subgraph "P2P Network Mesh"
        Node1[Node 1<br/>:3000<br/>Bootstrap]
        Node2[Node 2<br/>:7000<br/>Bootstrap]
        Node3[Node 3<br/>:5000<br/>Regular]
        Node4[Node 4<br/>:8000<br/>Regular]
        Node5[Node 5<br/>:9000<br/>Regular]
        
        Node1 <--> Node2
        Node1 <--> Node3
        Node1 <--> Node4
        Node1 <--> Node5
        Node2 <--> Node3
        Node2 <--> Node4
        Node2 <--> Node5
        Node3 <--> Node4
        Node3 <--> Node5
        Node4 <--> Node5
    end
    
    subgraph "Connection Types"
        Bootstrap[Bootstrap Connection]
        Regular[Regular P2P Connection]
        Data[Data Transfer Stream]
    end
    
    classDef bootstrapNode fill:#ffebee,stroke:#c62828,stroke-width:3px
    classDef regularNode fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef connectionType fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class Node1,Node2 bootstrapNode
    class Node3,Node4,Node5 regularNode
    class Bootstrap,Regular,Data connectionType
```

## Node Bootstrap Process

```mermaid
sequenceDiagram
    participant NewNode as New Node (:8000)
    participant Bootstrap1 as Bootstrap Node 1 (:3000)
    participant Bootstrap2 as Bootstrap Node 2 (:7000)
    participant ExistingPeer as Existing Peer (:5000)
    
    Note over NewNode: Node starts with bootstrap addresses
    NewNode->>Bootstrap1: TCP Dial (:3000)
    NewNode->>Bootstrap2: TCP Dial (:7000)
    
    Bootstrap1->>Bootstrap1: Accept connection
    Bootstrap2->>Bootstrap2: Accept connection
    
    Bootstrap1->>NewNode: Handshake
    Bootstrap2->>NewNode: Handshake
    
    NewNode->>Bootstrap1: Handshake response
    NewNode->>Bootstrap2: Handshake response
    
    Bootstrap1->>Bootstrap1: Add peer to peer list
    Bootstrap2->>Bootstrap2: Add peer to peer list
    
    Note over Bootstrap1,Bootstrap2: Nodes can discover each other through gossip
    Bootstrap1->>ExistingPeer: Notify new peer
    Bootstrap2->>ExistingPeer: Notify new peer
    
    ExistingPeer->>NewNode: Establish connection
    NewNode->>ExistingPeer: Accept connection
    
    Note over NewNode,ExistingPeer: Full mesh network established
```

## Message Broadcasting

```mermaid
flowchart TD
    subgraph "Message Broadcast Flow"
        A[Node 1 initiates<br/>file store operation]
        B[Create MessageStoreFile]
        C[Encode message with GOB]
        D[Get all connected peers]
        E[Send to Peer 1]
        F[Send to Peer 2]
        G[Send to Peer 3]
        H[Send to Peer N]
    end
    
    subgraph "Message Processing"
        I[Peer receives message]
        J[Decode GOB message]
        K[Identify message type]
        L[Route to handler]
        M[Process operation]
        N[Send response/confirmation]
    end
    
    A --> B
    B --> C
    C --> D
    D --> E
    D --> F
    D --> G
    D --> H
    
    E --> I
    F --> I
    G --> I
    H --> I
    
    I --> J
    J --> K
    K --> L
    L --> M
    M --> N
    
    classDef initiator fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef broadcast fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef processing fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A,B,C initiator
    class D,E,F,G,H broadcast
    class I,J,K,L,M,N processing
```

## TCP Transport Layer

```mermaid
graph TB
    subgraph "TCP Transport Architecture"
        subgraph "Node A"
            ListenerA[TCP Listener<br/>:3000]
            DecoderA[Message Decoder]
            HandlerA[Message Handler]
            PeerManagerA[Peer Manager]
            
            ListenerA --> DecoderA
            DecoderA --> HandlerA
            HandlerA --> PeerManagerA
        end
        
        subgraph "Node B"
            ListenerB[TCP Listener<br/>:7000]
            DecoderB[Message Decoder]
            HandlerB[Message Handler]
            PeerManagerB[Peer Manager]
            
            ListenerB --> DecoderB
            DecoderB --> HandlerB
            HandlerB --> PeerManagerB
        end
        
        subgraph "Network Layer"
            TCP[TCP Connections]
            RPC[RPC Messages]
            Stream[Data Streams]
        end
        
        PeerManagerA <--> TCP
        PeerManagerB <--> TCP
        
        TCP --> RPC
        TCP --> Stream
    end
    
    classDef nodeComponent fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef networkComponent fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    
    class ListenerA,DecoderA,HandlerA,PeerManagerA,ListenerB,DecoderB,HandlerB,PeerManagerB nodeComponent
    class TCP,RPC,Stream networkComponent
```

## Message Types and Handling

```mermaid
flowchart TD
    subgraph "Incoming Message Processing"
        A[TCP Connection] --> B[Read Message]
        B --> C[Decode GOB]
        C --> D{Message Type?}
        
        D -->|Store| E[MessageStoreFile]
        D -->|Get| F[MessageGetFile]
        D -->|Delete| G[MessageDeleteFile]
        
        E --> H[handleMessageStoreFile]
        F --> I[handleMessageGetFile]
        G --> J[handleMessageDeleteFile]
        
        H --> K[Store file locally]
        I --> L[Read file from disk]
        J --> M[Delete file from disk]
        
        K --> N[Send confirmation]
        L --> O[Stream file data]
        M --> P[Send confirmation]
    end
    
    classDef transport fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef message fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef handler fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef operation fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef response fill:#ffebee,stroke:#c62828,stroke-width:2px
    
    class A,B,C transport
    class D,E,F,G message
    class H,I,J handler
    class K,L,M operation
    class N,O,P response
```

## Data Transfer Protocols

```mermaid
sequenceDiagram
    participant Sender as Sender Node
    participant Receiver as Receiver Node
    
    Note over Sender,Receiver: File Store Operation
    Sender->>Receiver: MessageStoreFile (metadata)
    Receiver->>Receiver: Prepare for incoming stream
    Sender->>Receiver: IncomingStream marker
    Sender->>Receiver: Encrypted file data (stream)
    Receiver->>Receiver: Write data to disk
    Receiver->>Sender: Confirmation
    
    Note over Sender,Receiver: File Get Operation
    Sender->>Receiver: MessageGetFile (request)
    Receiver->>Receiver: Read file from disk
    Receiver->>Sender: IncomingStream marker
    Receiver->>Sender: File size (8 bytes)
    Receiver->>Sender: Encrypted file data (stream)
    Sender->>Sender: Write data to disk
    
    Note over Sender,Receiver: File Delete Operation
    Sender->>Receiver: MessageDeleteFile (request)
    Receiver->>Receiver: Delete file from disk
    Receiver->>Sender: Confirmation
```

## Peer Management

```mermaid
graph TB
    subgraph "Peer Management System"
        A[Peer Discovery] --> B[Connection Establishment]
        B --> C[Handshake Process]
        C --> D[Peer Registration]
        D --> E[Active Peer Pool]
        
        E --> F[Peer Monitoring]
        F --> G{Peer Alive?}
        G -->|Yes| H[Maintain Connection]
        G -->|No| I[Remove from Pool]
        
        H --> F
        I --> J[Attempt Reconnection]
        J --> K{Reconnect Success?}
        K -->|Yes| D
        K -->|No| L[Mark as Failed]
        
        E --> M[Message Broadcasting]
        M --> N[Send to All Peers]
        
        E --> O[Load Balancing]
        O --> P[Select Best Peer]
    end
    
    classDef discovery fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef management fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef monitoring fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef operations fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class A,B,C,D discovery
    class E,H,I,J,K,L management
    class F,G monitoring
    class M,N,O,P operations
```

## Network Fault Tolerance

```mermaid
flowchart TD
    subgraph "Fault Tolerance Mechanisms"
        A[Network Partition] --> B[Detect Lost Peers]
        B --> C[Update Peer List]
        C --> D[Continue Operations]
        
        E[Node Failure] --> F[Connection Timeout]
        F --> G[Remove from Active Peers]
        G --> H[Redistribute Requests]
        
        I[Bootstrap Node Down] --> J[Use Alternative Bootstrap]
        J --> K[Maintain Network Connectivity]
        
        L[Message Failure] --> M[Retry Mechanism]
        M --> N{Retry Successful?}
        N -->|Yes| O[Continue Operation]
        N -->|No| P[Mark Peer as Failed]
        
        D --> Q[Network Healing]
        H --> Q
        K --> Q
        O --> Q
        P --> Q
        
        Q --> R[Reconnect to Available Peers]
        R --> S[Restore Full Connectivity]
    end
    
    classDef fault fill:#ffebee,stroke:#c62828,stroke-width:2px
    classDef detection fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef recovery fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef healing fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    
    class A,E,I,L fault
    class B,F,J,M,N detection
    class C,D,G,H,K,O,P recovery
    class Q,R,S healing
```

## Key P2P Features

1. **Mesh Topology**: Every node can communicate with every other node
2. **Bootstrap Nodes**: Special nodes that help new nodes join the network
3. **Peer Discovery**: Automatic discovery of peers through bootstrap nodes
4. **Message Broadcasting**: Efficient distribution of messages to all peers
5. **Fault Tolerance**: Network continues operating despite node failures
6. **Load Distribution**: Files are distributed across multiple nodes
7. **Stream-based Transfer**: Efficient transfer of large files
8. **Connection Pooling**: Reuse of TCP connections for multiple operations