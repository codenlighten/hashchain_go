# Contributing to HashChain

Thank you for your interest in contributing to HashChain! This document provides guidelines and instructions for contributing to the project.

## Development Environment Setup

1. Install Go 1.18 (required for compatibility)
2. Install Protocol Buffer tools:
   - protoc compiler
   - protoc-gen-go v1.28.1
   - protoc-gen-go-grpc v1.2.0

## Building and Testing

1. Generate Protocol Buffer code:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/blockchain.proto
```

2. Run tests:
```bash
go test ./...
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Update documentation if needed
5. Run tests to ensure everything works
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to your branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Code Style

- Follow standard Go code style and conventions
- Use `gofmt` to format your code
- Add comments for exported functions and types
- Write meaningful commit messages

## Project Structure

```
hashchain_go/
├── cmd/           # Command-line applications
├── pkg/           # Library code
└── proto/         # Protocol Buffer definitions
```

## Current Focus Areas

1. Transaction Validation
   - Input/output validation
   - UTXO tracking
   - Basic scripting support

2. P2P Networking
   - Peer discovery
   - Block/transaction propagation
   - Node sync protocol

3. Security Features
   - Digital signatures
   - Public/private key pairs
   - Basic wallet functionality

## Questions or Need Help?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the codebase
- Improvement suggestions

## License

By contributing to HashChain, you agree that your contributions will be licensed under the MIT License.
