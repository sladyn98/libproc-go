package libproc

import (
	"syscall"
	"unsafe"
)

// proc_pidinfo flavor constants
const (
	PROC_PIDLISTFDS       = 1
	PROC_PIDTBSDINFO      = 3
	PROC_PIDTASKINFO      = 4
	PROC_PIDTHREADINFO    = 5
	PROC_PIDVNODEPATHINFO = 9
)

// File descriptor types
const (
	PROX_FDTYPE_ATALK     = 0
	PROX_FDTYPE_VNODE     = 1
	PROX_FDTYPE_SOCKET    = 2
	PROX_FDTYPE_PSHM      = 3
	PROX_FDTYPE_PSEM      = 4
	PROX_FDTYPE_KQUEUE    = 5
	PROX_FDTYPE_PIPE      = 6
	PROX_FDTYPE_FSEVENTS  = 7
	PROX_FDTYPE_NETPOLICY = 9
)

// ProcBSDInfo represents basic process information
type ProcBSDInfo struct {
	Flags     uint32
	Status    uint32
	Xstatus   uint32
	Pid       uint32
	Ppid      uint32
	Uid       uint32
	Gid       uint32
	Ruid      uint32
	Rgid      uint32
	Svuid     uint32
	Svgid     uint32
	_         uint32 // padding
	Comm      [16]byte
	Name      [32]byte
	Nfiles    uint32
	Pgid      uint32
	Pjobc     uint32
	_         uint32 // e_tdev
	_         uint32 // e_tpgid
	Nice      int32
	StartTime [8]byte // struct timeval
	_         [8]byte // padding to match C struct size (136 bytes total)
}

// ProcFdInfo represents file descriptor information
type ProcFdInfo struct {
	Fd     int32
	FdType uint32
}

//go:cgo_import_dynamic libc_proc_listallpids proc_listallpids "/usr/lib/libproc.dylib"
//go:cgo_import_dynamic libc_proc_name proc_name "/usr/lib/libproc.dylib"
//go:cgo_import_dynamic libc_proc_pidinfo proc_pidinfo "/usr/lib/libproc.dylib"

var libc_proc_listallpids_trampoline_addr uintptr
var libc_proc_name_trampoline_addr uintptr
var libc_proc_pidinfo_trampoline_addr uintptr

//go:linkname syscall_rawSyscall syscall.rawSyscall
func syscall_rawSyscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:linkname syscall_syscall6 syscall.syscall6
func syscall_syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

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

//go:nosplit
func proc_pidinfo(pid int, flavor int, arg uint64, buffer unsafe.Pointer, buffersize int) (ret int, err error) {
	r0, _, e1 := syscall_syscall6(
		libc_proc_pidinfo_trampoline_addr,
		uintptr(pid),
		uintptr(flavor),
		uintptr(arg),
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

func PidName(pid int) (string, error) {
	const bufSize = 2 * 1024 // MAXCOMLEN is typically 16, but buffer for safety
	buf := make([]byte, bufSize)
	n, err := proc_name(pid, unsafe.Pointer(&buf[0]), bufSize)
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

// ListAllPids returns the count of all running processes on the system
func ListAllPids() (int, error) {
	numBytes, err := proc_listallpids(nil, 0)
	if err != nil {
		return 0, err
	}
	numPids := numBytes / 4
	return numPids, nil
}

func ProcPidInfo(pid int, flavor int, arg uint64, buffer unsafe.Pointer, buffersize int) (int, error) {
	return proc_pidinfo(pid, flavor, arg, buffer, buffersize)
}

// GetProcBSDInfo retrieves basic BSD process information for the given PID
func GetProcBSDInfo(pid int) (*ProcBSDInfo, error) {
	var info ProcBSDInfo
	size := int(unsafe.Sizeof(info))

	n, err := proc_pidinfo(pid, PROC_PIDTBSDINFO, 0, unsafe.Pointer(&info), size)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, syscall.EINVAL
	}

	return &info, nil
}

// ListPidFileDescriptors returns a list of file descriptors for the given PID
func ListPidFileDescriptors(pid int) ([]ProcFdInfo, error) {
	// First call to get the required buffer size
	n, err := proc_pidinfo(pid, PROC_PIDLISTFDS, 0, nil, 0)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return []ProcFdInfo{}, nil
	}

	// Allocate buffer and get the actual data
	count := n / int(unsafe.Sizeof(ProcFdInfo{}))
	fds := make([]ProcFdInfo, count)

	n, err = proc_pidinfo(pid, PROC_PIDLISTFDS, 0, unsafe.Pointer(&fds[0]), n)
	if err != nil {
		return nil, err
	}

	actualCount := n / int(unsafe.Sizeof(ProcFdInfo{}))
	return fds[:actualCount], nil
}
