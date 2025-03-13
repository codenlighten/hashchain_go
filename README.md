# HashChain

A lightweight blockchain implementation in Go, inspired by Bitcoin SV (BSV). This project implements core blockchain functionality using modern Go practices and Protocol Buffers for efficient message serialization.

## Features

- BSV-inspired block structure
- Proof of Work consensus mechanism
- Transaction mempool for unconfirmed transactions
- P2P networking using gRPC
- Block validation and chain management
- Double SHA256 for block hashing
- Merkle tree for transaction validation
- Configurable mining difficulty

## Project Structure

```
hashchain_go/
├── cmd/
│   ├── client/     # Client implementation for testing
│   └── node/       # Node server implementation
├── pkg/
│   ├── blockchain/ # Core blockchain implementation
│   └── node/       # Node service and P2P networking
└── proto/          # Protocol Buffer definitions
```

## Requirements

- Go 1.18 or later
- Protocol Buffers compiler (protoc)
- protoc-gen-go and protoc-gen-go-grpc plugins

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/hashchain_go.git
cd hashchain_go
```

2. Install dependencies:
```bash
go mod download
```

3. Install Protocol Buffers plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

4. Generate Protocol Buffers code:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/blockchain.proto
```

## Usage

1. Start the node server:
```bash
go run cmd/node/main.go --port 50051
```

2. Run the test client:
```bash
go run cmd/client/main.go --server localhost:50051
```

The client will:
1. Create and broadcast a sample transaction
2. Create a new block with the transaction
3. Mine the block using Proof of Work
4. Broadcast the mined block to the network

## Current Status

- Basic blockchain structure
- Block mining with Proof of Work
- Transaction handling and mempool
- P2P communication via gRPC
- Block validation and chain management
- Merkle tree implementation
- Genesis block handling

## Roadmap

### Short-term Goals
1. Improve transaction validation
   - Add input/output validation
   - Implement UTXO (Unspent Transaction Output) tracking
   - Add basic scripting support

2. Enhance P2P Networking
   - Add peer discovery
   - Implement block and transaction propagation
   - Add sync protocol for new nodes

3. Improve Security
   - Add digital signatures
   - Implement public/private key pairs
   - Add basic wallet functionality

### Medium-term Goals
1. Chain Management
   - Handle chain reorganization
   - Implement checkpoints
   - Add pruning capabilities

2. Performance Optimization
   - Optimize block validation
   - Improve mining performance
   - Add block storage optimization

3. Developer Tools
   - Add CLI tools
   - Create block explorer
   - Add testing framework

### Long-term Goals
1. Advanced Features
   - Smart contract support
   - Custom consensus rules
   - Advanced scripting capabilities

2. Network Improvements
   - Add SPV (Simplified Payment Verification)
   - Implement mempool management policies
   - Add network health monitoring

3. Production Readiness
   - Add extensive documentation
   - Implement logging and metrics
   - Add production deployment guides

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
