# Drift - Node Lifecycle and Bootstrap Process

## Node Lifecycle Overview
This diagram shows the complete lifecycle of a node in the Drift distributed file system, from initialization to joining the network and handling failures.

## Node Initialization Process

```mermaid
flowchart TD
    subgraph "Node Startup"
        A[Node Starts] --> B[Parse Configuration]
        B --> C[Generate Node ID]
        C --> D[Generate Encryption Key]
        D --> E[Initialize Store]
        E --> F[Setup TCP Transport]
        F --> G[Configure Message Handlers]
    end
    
    subgraph "Network Bootstrap"
        H[Start TCP Listener] --> I[Connect to Bootstrap Nodes]
        I --> J{Bootstrap Nodes Available?}
        J -->|Yes| K[Establish Connections]
        J -->|No| L[Start as Bootstrap Node]
        K --> M[Perform Handshake]
        L --> N[Wait for Incoming Connections]
    end
    
    subgraph "Network Integration"
        O[Register with Peers] --> P[Start Message Loop]
        P --> Q[Begin File Operations]
        Q --> R[Node Fully Operational]
    end
    
    G --> H
    M --> O
    N --> O
    
    classDef startup fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef bootstrap fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef integration fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A,B,C,D,E,F,G startup
    class H,I,J,K,L,M,N bootstrap
    class O,P,Q,R integration
```

## Bootstrap Network Formation

```mermaid
sequenceDiagram
    participant Node1 as Node 1<br/>(:3000)
    participant Node2 as Node 2<br/>(:7000)
    participant Node3 as Node 3<br/>(:5000)
    
    Note over Node1,Node3: Initial Bootstrap Nodes
    Node1->>Node1: Start as bootstrap (:3000)
    Node2->>Node2: Start as bootstrap (:7000)
    
    Note over Node1,Node2: Bootstrap nodes establish connection
    Node1->>Node2: Discover and connect
    Node2->>Node1: Accept connection
    Node1->>Node2: Exchange peer information
    Node2->>Node1: Acknowledge
    
    Note over Node1,Node3: New Node Joins Network
    Node3->>Node3: Start with bootstrap list [:3000, :7000]
    Node3->>Node1: Connect to :3000
    Node3->>Node2: Connect to :7000
    
    Node1->>Node3: Handshake
    Node2->>Node3: Handshake
    
    Node3->>Node1: Register as peer
    Node3->>Node2: Register as peer
    
    Note over Node1,Node3: Peer Discovery
    Node1->>Node3: Share peer list
    Node2->>Node3: Share peer list
    
    Node3->>Node1: Connect to additional peers
    Node3->>Node2: Connect to additional peers
    
    Note over Node1,Node3: Full Mesh Network Established
```

## Node State Management

```mermaid
stateDiagram-v2
    [*] --> Initializing
    Initializing --> Connecting: Configuration Complete
    Connecting --> Bootstrapping: TCP Listener Started
    Bootstrapping --> PeerDiscovery: Bootstrap Successful
    Bootstrapping --> StandaloneBoot: No Bootstrap Nodes
    PeerDiscovery --> Operational: Peers Discovered
    StandaloneBoot --> Operational: Ready for Connections
    
    Operational --> Operational: File Operations
    Operational --> Reconnecting: Connection Lost
    Operational --> Shutting: Graceful Shutdown
    
    Reconnecting --> Operational: Reconnection Successful
    Reconnecting --> Isolated: Max Retries Reached
    
    Isolated --> Reconnecting: Retry Timer
    Isolated --> Shutting: Manual Shutdown
    
    Shutting --> [*]: Cleanup Complete
    
    note right of Operational
        Node can perform:
        - File storage
        - File retrieval
        - Peer communication
        - Network replication
    end note
```

## Peer Discovery and Connection Management

```mermaid
flowchart TD
    subgraph "Peer Discovery Process"
        A[New Node Joins] --> B[Connect to Bootstrap Nodes]
        B --> C[Receive Peer List]
        C --> D[Attempt Connections]
        D --> E{Connection Successful?}
        E -->|Yes| F[Add to Peer Pool]
        E -->|No| G[Mark as Failed]
        F --> H[Exchange Metadata]
        G --> I[Try Next Peer]
        H --> J[Peer Fully Integrated]
        I --> D
    end
    
    subgraph "Connection Health Monitoring"
        K[Monitor Peer Connections] --> L{Peer Responsive?}
        L -->|Yes| M[Maintain Connection]
        L -->|No| N[Attempt Reconnection]
        M --> K
        N --> O{Reconnection Successful?}
        O -->|Yes| P[Restore to Peer Pool]
        O -->|No| Q[Remove from Pool]
        P --> K
        Q --> R[Update Peer List]
        R --> K
    end
    
    J --> K
    
    classDef discovery fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef monitoring fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef management fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A,B,C,D,E,F,G,H,I,J discovery
    class K,L,M,N,O,P,Q,R monitoring
```

## Node Configuration and Initialization

```mermaid
graph TB
    subgraph "FileServerOpts Configuration"
        A[Node ID<br/>Generated/Provided]
        B[Encryption Key<br/>32 bytes random]
        C[Storage Root<br/>Local directory]
        D[Path Transform<br/>CAS function]
        E[Transport Layer<br/>TCP transport]
        F[Bootstrap Nodes<br/>List of addresses]
    end
    
    subgraph "Transport Configuration"
        G[Listen Address<br/>:port]
        H[Handshake Function<br/>NOPHandshakeFunc]
        I[Message Decoder<br/>GOB decoder]
        J[Peer Handler<br/>OnPeer callback]
    end
    
    subgraph "Store Configuration"
        K[Storage Root<br/>File system path]
        L[Path Transform<br/>Hash-based paths]
        M[Permissions<br/>Directory creation]
    end
    
    A --> G
    B --> H
    C --> I
    D --> J
    E --> K
    F --> L
    
    G --> M
    H --> M
    I --> M
    J --> M
    K --> M
    L --> M
    
    classDef config fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef transport fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef store fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A,B,C,D,E,F config
    class G,H,I,J transport
    class K,L,M store
```

## Network Partition Handling

```mermaid
sequenceDiagram
    participant Node1
    participant Node2
    participant Node3
    participant Network
    
    Note over Node1,Node3: Normal Operation
    Node1->>Node2: File operations
    Node2->>Node3: File operations
    Node3->>Node1: File operations
    
    Note over Network: Network Partition Occurs
    Network->>Network: Partition between Node1 and Node2,3
    
    Note over Node1: Partition Detection
    Node1->>Node2: Heartbeat (timeout)
    Node1->>Node3: Heartbeat (timeout)
    Node1->>Node1: Detect network partition
    
    Note over Node2,Node3: Remaining Nodes Continue
    Node2->>Node3: Continue operations
    Node3->>Node2: Continue operations
    
    Note over Node1: Isolated Node Behavior
    Node1->>Node1: Continue local operations
    Node1->>Node1: Queue replication operations
    
    Note over Network: Network Heals
    Network->>Network: Partition resolved
    
    Note over Node1,Node3: Reconnection Process
    Node1->>Node2: Reconnect attempt
    Node1->>Node3: Reconnect attempt
    
    Node2->>Node1: Accept reconnection
    Node3->>Node1: Accept reconnection
    
    Note over Node1,Node3: Synchronization
    Node1->>Node2: Sync queued operations
    Node1->>Node3: Sync queued operations
    
    Node2->>Node1: Sync missed operations
    Node3->>Node1: Sync missed operations
```

## Node Failure and Recovery

```mermaid
flowchart TD
    subgraph "Failure Detection"
        A[Connection Timeout] --> B[Heartbeat Failure]
        B --> C[Mark Node as Suspected]
        C --> D[Retry Connection]
        D --> E{Recovery Successful?}
        E -->|Yes| F[Restore Node Status]
        E -->|No| G[Mark as Failed]
    end
    
    subgraph "Network Adaptation"
        H[Remove from Peer Pool] --> I[Update Routing Tables]
        I --> J[Redistribute File Requests]
        J --> K[Continue Operations]
    end
    
    subgraph "Node Recovery"
        L[Node Comes Online] --> M[Attempt Bootstrap]
        M --> N[Reconnect to Network]
        N --> O[Sync Missed Operations]
        O --> P[Resume Normal Operations]
    end
    
    G --> H
    F --> K
    K --> L
    P --> K
    
    classDef detection fill:#ffebee,stroke:#c62828,stroke-width:2px
    classDef adaptation fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef recovery fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A,B,C,D,E,F,G detection
    class H,I,J,K adaptation
    class L,M,N,O,P recovery
```

## Node Metrics and Monitoring

```mermaid
graph TB
    subgraph "Node Health Metrics"
        A[Connection Status]
        B[Peer Count]
        C[Storage Usage]
        D[Network Bandwidth]
        E[Operation Latency]
        F[Error Rate]
    end
    
    subgraph "Performance Monitoring"
        G[File Operation Throughput]
        H[Network Message Rate]
        I[Storage I/O Rate]
        J[CPU Usage]
        K[Memory Usage]
        L[Disk Usage]
    end
    
    subgraph "Network Monitoring"
        M[Peer Connectivity]
        N[Message Success Rate]
        O[Replication Status]
        P[Network Partition Detection]
        Q[Bootstrap Success Rate]
        R[Handshake Failure Rate]
    end
    
    A --> G
    B --> H
    C --> I
    D --> J
    E --> K
    F --> L
    
    G --> M
    H --> N
    I --> O
    J --> P
    K --> Q
    L --> R
    
    classDef health fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef performance fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef network fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    
    class A,B,C,D,E,F health
    class G,H,I,J,K,L performance
    class M,N,O,P,Q,R network
```

## Key Lifecycle Features

1. **Automatic Bootstrap**: Nodes automatically discover and join the network
2. **Peer Discovery**: Dynamic discovery of peers through bootstrap nodes
3. **Fault Tolerance**: Graceful handling of node failures and network partitions
4. **Self-Healing**: Automatic reconnection and synchronization
5. **Health Monitoring**: Continuous monitoring of node and network health
6. **Configuration Management**: Flexible configuration of node parameters
7. **Graceful Shutdown**: Clean shutdown with proper resource cleanup
8. **Recovery Mechanisms**: Automatic recovery from various failure scenarios