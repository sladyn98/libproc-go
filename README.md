# libproc

[![Go Reference](https://pkg.go.dev/badge/github.com/sladyn98/libproc-go.svg)](https://pkg.go.dev/github.com/sladyn98/libproc-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sladyn98/libproc-go)](https://goreportcard.com/report/github.com/sladyn98/libproc-go)
[![CI](https://github.com/sladyn98/libproc-go/workflows/CI/badge.svg)](https://github.com/sladyn98/libproc-go/actions)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

A pure Go library providing access to macOS process information through the libproc system library.

## Why This Library?

### Stability Over Direct Syscalls

This library uses **libproc.dylib** instead of direct syscalls because:

- **Apple maintains libproc.dylib as a stable API** - The dylib interface remains consistent across macOS versions, even when underlying syscall numbers change
- **Syscall numbers are unstable** - Apple frequently changes syscall numbers between macOS releases, breaking code that uses direct syscalls
- **Binary compatibility** - Programs built with this library continue to work across macOS updates without recompilation
- **Official interface** - libproc.dylib is the documented way Apple provides for accessing process information

Direct syscall approach:
```go
// ❌ FRAGILE: Syscall number 336 may change between macOS versions
syscall.Syscall(336, ...) // May break on next macOS update
```

This library's approach:
```go
// ✅ STABLE: Uses libproc.dylib which Apple maintains across versions
libproc.ListAllPids() // Works reliably across macOS updates
```

### How It Works

Uses CGO dynamic linking with assembly trampolines to call libproc.dylib functions:
- `//go:cgo_import_dynamic` - Links to macOS system library at runtime
- Assembly trampolines - Minimal overhead jump stubs
- `syscall.rawSyscall` - Direct function pointer invocation
- Zero C code - Pure Go with platform-specific assembly

## Supported Platforms

- macOS (darwin) - AMD64 (Intel)
- macOS (darwin) - ARM64 (Apple Silicon)

## Installation

```bash
go get github.com/sladyn98/libproc-go
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sladyn98/libproc-go"
)

func main() {
    // Get count of running processes
    count, err := libproc.ListAllPids()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Running processes: %d\n", count)

    // Get process name by PID
    name, err := libproc.PidName(os.Getpid())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Current process: %s\n", name)
}
```

## API

### ListAllPids

```go
func ListAllPids() (int, error)
```

Returns the count of all running processes on the system.

### PidName

```go
func PidName(pid int) (string, error)
```

Returns the process name for a given process ID. System processes (like kernel_task or launchd) may return empty strings due to macOS security restrictions.

## Development

### Building

```bash
go build
```

### Testing

```bash
# Run all tests
go test ./tests/

# With coverage
go test -coverprofile=coverage.out ./tests/
go tool cover -html=coverage.out

# With race detection
go test -race ./tests/
```

### Project Structure

```
libproc/
├── README.md              # This file
├── LICENSE                # BSD-3-Clause license
├── CONTRIBUTING.md        # Contribution guidelines
├── doc.go                 # Package documentation
├── syscall.go             # Go implementation
├── syscall_arm64.s        # ARM64 assembly trampolines
└── tests/                 # Tests
    └── libproc_test.go    # Unit tests
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

BSD 3-Clause License - see [LICENSE](LICENSE) for details.

## See Also

- [Apple libproc Documentation](https://developer.apple.com/documentation/kernel)
- [Go Assembly Guide](https://go.dev/doc/asm)
