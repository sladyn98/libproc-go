package libproc

import (
	"syscall"
	"unsafe"
)

//go:cgo_import_dynamic libc_proc_listallpids proc_listallpids "/usr/lib/libproc.dylib"
//go:cgo_import_dynamic libc_proc_name proc_name "/usr/lib/libproc.dylib"

var libc_proc_listallpids_trampoline_addr uintptr
var libc_proc_name_trampoline_addr uintptr

//go:linkname syscall_rawSyscall syscall.rawSyscall
func syscall_rawSyscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:nosplit
func proc_listallpids(buffer unsafe.Pointer, buffersize int) (ret int, err error) {
	r0, _, e1 := syscall_rawSyscall(
		libc_proc_listallpids_trampoline_addr,
		uintptr(buffer),
		uintptr(buffersize),
		0,
	)
	ret = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

//go:nosplit
func proc_name(pid int, buffer unsafe.Pointer, buffersize uint32) (ret int, err error) {
	r0, _, e1 := syscall_rawSyscall(
		libc_proc_name_trampoline_addr,
		uintptr(pid),
		uintptr(buffer),
		uintptr(buffersize),
	)
	ret = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func PidName(pid int) (string, error) {
	buf := make([]byte, 2*1024) // MAXCOMLEN is typically 16, but buffer for safety
	n, err := proc_name(pid, unsafe.Pointer(&buf[0]), uint32(len(buf)))
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", nil
	}
	// Find null terminator
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return string(buf[:i]), nil
		}
	}
	return string(buf[:n]), nil
}

// Public API function with proper error handling
func ListAllPids() (int, error) {
	numBytes, err := proc_listallpids(nil, 0)
	if err != nil {
		return 0, err
	}
	numPids := numBytes / 4
	return numPids, nil
}
