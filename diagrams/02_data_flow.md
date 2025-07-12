# Drift - Data Flow Diagram

## Data Flow Overview
This diagram illustrates how data flows through the Drift distributed file system during different operations (Store, Get, Delete).

## Store Operation Data Flow

```mermaid
sequenceDiagram
    participant Client
    participant FileServer
    participant Store
    participant Crypto
    participant Network
    participant RemotePeer1
    participant RemotePeer2
    
    Client->>FileServer: Store("file.txt", data)
    FileServer->>Store: Write(nodeID, "file.txt", data)
    Store->>Store: Generate SHA1 hash key
    Store->>Store: Create directory structure
    Store->>Store: Write file to local disk
    Store-->>FileServer: Return file size
    
    FileServer->>Crypto: Encrypt(data, encKey)
    Crypto->>Crypto: Generate IV
    Crypto->>Crypto: AES-256 CTR encryption
    Crypto-->>FileServer: Return encrypted data
    
    FileServer->>Network: Broadcast(MessageStoreFile)
    Network->>RemotePeer1: Send store message
    Network->>RemotePeer2: Send store message
    
    FileServer->>RemotePeer1: Stream encrypted data
    FileServer->>RemotePeer2: Stream encrypted data
    
    RemotePeer1->>RemotePeer1: Store encrypted file
    RemotePeer2->>RemotePeer2: Store encrypted file
    
    RemotePeer1-->>FileServer: Confirm storage
    RemotePeer2-->>FileServer: Confirm storage
    
    FileServer-->>Client: Store operation complete
```

## Get Operation Data Flow

```mermaid
sequenceDiagram
    participant Client
    participant FileServer
    participant Store
    participant Crypto
    participant Network
    participant RemotePeer
    
    Client->>FileServer: Get("file.txt")
    FileServer->>Store: Has(nodeID, "file.txt")?
    
    alt File exists locally
        Store-->>FileServer: File found
        Store->>Store: Read file from disk
        Store-->>FileServer: Return file data
        FileServer-->>Client: Return file data
    else File not found locally
        Store-->>FileServer: File not found
        FileServer->>Network: Broadcast(MessageGetFile)
        Network->>RemotePeer: Send get message
        
        RemotePeer->>RemotePeer: Check if file exists
        RemotePeer->>RemotePeer: Read file from disk
        RemotePeer->>Network: Send file size
        RemotePeer->>Network: Stream encrypted data
        
        Network->>FileServer: Receive file size
        Network->>FileServer: Receive encrypted data
        
        FileServer->>Crypto: Decrypt(encryptedData, encKey)
        Crypto->>Crypto: Extract IV
        Crypto->>Crypto: AES-256 CTR decryption
        Crypto-->>FileServer: Return decrypted data
        
        FileServer->>Store: WriteDecrypt(nodeID, "file.txt", data)
        Store->>Store: Store file locally
        Store-->>FileServer: Confirm storage
        
        FileServer->>Store: Read(nodeID, "file.txt")
        Store-->>FileServer: Return file data
        FileServer-->>Client: Return file data
    end
```

## Delete Operation Data Flow

```mermaid
sequenceDiagram
    participant Client
    participant FileServer
    participant Store
    participant Network
    participant RemotePeer1
    participant RemotePeer2
    
    Client->>FileServer: Delete("file.txt")
    FileServer->>Store: Has(nodeID, "file.txt")?
    
    alt File exists
        Store-->>FileServer: File found
        
        FileServer->>Network: Broadcast(MessageDeleteFile)
        Network->>RemotePeer1: Send delete message
        Network->>RemotePeer2: Send delete message
        
        RemotePeer1->>RemotePeer1: Delete file from disk
        RemotePeer2->>RemotePeer2: Delete file from disk
        
        RemotePeer1-->>FileServer: Confirm deletion
        RemotePeer2-->>FileServer: Confirm deletion
        
        FileServer->>Store: Delete(nodeID, "file.txt")
        Store->>Store: Remove file and directories
        Store-->>FileServer: Confirm deletion
        
        FileServer-->>Client: Delete operation complete
    else File not found
        Store-->>FileServer: File not found
        FileServer-->>Client: Error: File not found
    end
```

## Content-Addressable Storage Flow

```mermaid
flowchart TD
    A[Input: File Key] --> B[Generate SHA1 Hash]
    B --> C[Hash: e.g., 3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b]
    C --> D[Split Hash into Blocks]
    D --> E[Block 1: 3a4b5]
    D --> F[Block 2: c6d7e]
    D --> G[Block 3: 8f9a0]
    D --> H[Block 4: b1c2d]
    D --> I[Continue...]
    
    E --> J[Create Directory Structure]
    F --> J
    G --> J
    H --> J
    I --> J
    
    J --> K[Path: /3a4b5/c6d7e/8f9a0/b1c2d/...]
    K --> L[Full Path: /root/nodeID/3a4b5/c6d7e/.../fullhash]
    L --> M[Store File at Full Path]
    
    classDef hashClass fill:#ffe0b2,stroke:#f57c00,stroke-width:2px
    classDef pathClass fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef storageClass fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    
    class B,C,D,E,F,G,H,I hashClass
    class J,K,L pathClass
    class M storageClass
```

## Encryption/Decryption Flow

```mermaid
flowchart TD
    subgraph "Encryption Flow"
        A1[Original File Data] --> B1[Generate Random IV]
        B1 --> C1[AES-256 CTR Cipher]
        A1 --> C1
        D1[32-byte Encryption Key] --> C1
        C1 --> E1[Encrypted Data]
        B1 --> F1[Prepend IV to Data]
        E1 --> F1
        F1 --> G1[IV + Encrypted Data]
    end
    
    subgraph "Decryption Flow"
        A2[IV + Encrypted Data] --> B2[Extract IV]
        A2 --> C2[Extract Encrypted Data]
        B2 --> D2[AES-256 CTR Cipher]
        C2 --> D2
        E2[32-byte Encryption Key] --> D2
        D2 --> F2[Decrypted Data]
        F2 --> G2[Original File Data]
    end
    
    classDef inputClass fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef cryptoClass fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef outputClass fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A1,A2,G2 inputClass
    class B1,C1,B2,D2 cryptoClass
    class G1,F2 outputClass
```

## Network Message Flow

```mermaid
flowchart TD
    subgraph "Message Types"
        A[MessageStoreFile]
        B[MessageGetFile] 
        C[MessageDeleteFile]
    end
    
    subgraph "Message Processing"
        D[Receive Message] --> E[Decode GOB]
        E --> F{Message Type?}
        F -->|Store| G[handleMessageStoreFile]
        F -->|Get| H[handleMessageGetFile]
        F -->|Delete| I[handleMessageDeleteFile]
    end
    
    subgraph "File Operations"
        G --> J[Store File Locally]
        H --> K[Read File from Disk]
        I --> L[Delete File from Disk]
    end
    
    subgraph "Network Response"
        J --> M[Send Confirmation]
        K --> N[Stream File Data]
        L --> O[Send Confirmation]
    end
    
    classDef messageClass fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef processClass fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef opClass fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef responseClass fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class A,B,C messageClass
    class D,E,F,G,H,I processClass
    class J,K,L opClass
    class M,N,O responseClass
```

## Key Data Flow Principles

1. **Replication**: Every file operation is replicated across all connected peers
2. **Content-Addressable**: Files are stored using SHA1-based paths for deduplication
3. **Encryption**: All data is encrypted before storage and decrypted after retrieval
4. **Fault Tolerance**: If a file isn't found locally, it's fetched from the network
5. **Consistency**: All peers maintain synchronized copies of files