# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Test Commands
- Run all tests: `go test ./...`
- Run a specific test: `go test -run TestName ./path/to/package`
- Run a specific subtest: `go test -run TestName/SubTestName ./path/to/package`
- Run tests with verbose output: `go test -v ./...`
- Format code: `gofmt -w .`

## Code Style Guidelines
- **Naming**: PascalCase for exported types/functions/constants, camelCase for unexported
- **Formatting**: Follow standard Go formatting with gofmt
- **Types**: Use pointer arguments for mutability, return explicit errors
- **Error Handling**: Return errors rather than panicking in library code
- **Documentation**: Document all exported functions and types with comments
- **Testing**: Use subtests with descriptive names, use testify/assert for assertions
- **Pointers**: Use utility functions (P, VOr, OrP) for pointer conversion
- **Architecture**: Prefer composition over inheritance, use interfaces