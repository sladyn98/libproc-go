// Package libproc provides access to macOS process information through the
// native libproc system library.
//
// This library uses a zero-C-code approach with CGO dynamic linking and
// assembly trampolines to efficiently call macOS libproc functions.
//
// # Supported Platforms
//
// This library currently supports:
//   - macOS (darwin) on AMD64 (Intel)
//   - macOS (darwin) on ARM64 (Apple Silicon)
//
// # Architecture
//
// The library implements a sophisticated CGO pattern that avoids writing C code:
//
//  1. Dynamic Library Import: Uses //go:cgo_import_dynamic to link against
//     /usr/lib/libproc.dylib at runtime
//
//  2. Assembly Trampolines: Small assembly stubs that jump to C functions
//     (syscall_amd64.s for Intel, syscall_arm64.s for Apple Silicon)
//
//  3. Raw Syscalls: Uses syscall.rawSyscall to invoke function pointers
//     with proper argument passing
//
//  4. Type-Safe Wrappers: Public Go functions handle type conversions
//     and error handling
//
// # Usage
//
// Basic usage example:
//
//	package main
//
//	import (
//		"fmt"
//		"log"
//		"os"
//
//		"github.com/yourusername/libproc"
//	)
//
//	func main() {
//		// Get count of running processes
//		count, err := libproc.ListAllPids()
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("Running processes: %d\n", count)
//
//		// Get process name
//		name, err := libproc.PidName(os.Getpid())
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("Current process: %s\n", name)
//	}
//
// # Security Considerations
//
// Some system processes (like kernel_task or launchd) may return empty strings
// or errors when querying their names due to macOS security restrictions.
// This is expected behavior and not a bug in the library.
//
// # Performance
//
// This library provides direct syscall access to libproc functions with minimal
// overhead. The assembly trampolines add negligible latency compared to direct
// C calls, while maintaining the benefits of pure Go code (no CGO compilation).
package libproc
