# Mitaine Project Guidelines

## Build Commands
- Build project: `go build ./cmd/mitaine`
- Run project: `go run ./cmd/mitaine/main.go`
- Run all tests: `go test ./...`
- Run specific test: `go test -v github.com/a-dot/mitaine/pkg/cards -run TestColor`
- Run tests for package: `go test ./pkg/cards`

## Code Style Guidelines
- **Imports**: Standard library first, then third-party, then project packages
- **Types**: Use type declarations with iota for enums and constants
- **Naming**: PascalCase for exported items, camelCase for unexported items
- **Receivers**: Pointer receivers when modifying state, value receivers for query methods
- **Error handling**: Check errors immediately, use log.Fatal or return as appropriate
- **Testing**: Use testify/require for assertions, table-driven tests with clear test cases
- **Organization**: Keep related types and functions in domain-specific packages
- **Variables**: Use var blocks for related variable declarations
- **Initialization**: Prefer factory functions for complex initializations