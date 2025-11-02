package libproc_test

import (
	"os"
	"testing"
	"unsafe"

	libproc "github.com/sladyn98/libproc-go"
)

func TestListAllPids(t *testing.T) {
	count, err := libproc.ListAllPids()
	if err != nil {
		t.Fatalf("ListAllPids() returned error: %v", err)
	}

	// Sanity checks
	if count < 10 {
		t.Errorf("ListAllPids() = %d, expected at least 10 processes", count)
	}

	if count > 10000 {
		t.Errorf("ListAllPids() = %d, unexpectedly high process count", count)
	}

	t.Logf("ListAllPids() returned %d processes", count)
}

func TestPidName_CurrentProcess(t *testing.T) {
	myPid := os.Getpid()
	name, err := libproc.PidName(myPid)
	if err != nil {
		t.Fatalf("PidName(%d) returned error: %v", myPid, err)
	}

	if name == "" {
		t.Errorf("PidName(%d) returned empty string for current process", myPid)
	}

	t.Logf("PidName(%d) = %q", myPid, name)
}

func TestPidName_ParentProcess(t *testing.T) {
	ppid := os.Getppid()
	name, err := libproc.PidName(ppid)
	if err != nil {
		t.Fatalf("PidName(%d) returned error: %v", ppid, err)
	}
	if name == "" {
		t.Errorf("PidName(%d) returned empty string for parent process", ppid)
	}
	t.Logf("PidName(%d) = %q (parent process)", ppid, name)
}

func TestPidName_SystemProcesses(t *testing.T) {
	tests := []struct {
		name string
		pid  int
	}{
		{"kernel_task", 0},
		{"launchd", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := libproc.PidName(tt.pid)

			// System processes may return empty string or error due to security restrictions
			// We just log what we got
			if err != nil {
				t.Logf("PidName(%d) returned error: %v (may be expected for system process)", tt.pid, err)
			} else {
				t.Logf("PidName(%d) = %q", tt.pid, name)
			}
		})
	}
}

func TestPidName_InvalidPid(t *testing.T) {
	tests := []struct {
		name string
		pid  int
	}{
		{"very_high", 99999},
		{"negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := libproc.PidName(tt.pid)

			// Invalid PIDs should either return error or empty string
			// Both are acceptable behaviors
			if err == nil && name != "" {
				t.Errorf("PidName(%d) = %q, nil; expected error or empty string for invalid PID", tt.pid, name)
			}

			t.Logf("PidName(%d) = %q, %v", tt.pid, name, err)
		})
	}
}

// TestPidName_TableDriven demonstrates table-driven test pattern
func TestPidName_TableDriven(t *testing.T) {
	// Get some valid PIDs
	myPid := os.Getpid()
	ppid := os.Getppid()

	tests := []struct {
		name      string
		pid       int
		wantErr   bool
		wantEmpty bool // true if empty string is acceptable
	}{
		{
			name:      "current_process",
			pid:       myPid,
			wantErr:   false,
			wantEmpty: false, // Current process should always have a name
		},
		{
			name:      "parent_process",
			pid:       ppid,
			wantErr:   false,
			wantEmpty: false, // Parent process should have a name
		},
		{
			name:      "system_process_launchd",
			pid:       1,
			wantErr:   false,
			wantEmpty: true, // System processes may return empty
		},
		{
			name:      "invalid_high_pid",
			pid:       99999,
			wantErr:   false, // May or may not error
			wantEmpty: true,  // Likely to be empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := libproc.PidName(tt.pid)

			if (err != nil) != tt.wantErr {
				t.Errorf("PidName(%d) error = %v, wantErr %v", tt.pid, err, tt.wantErr)
				return
			}

			if !tt.wantEmpty && got == "" {
				t.Errorf("PidName(%d) = %q, want non-empty string", tt.pid, got)
			}

			t.Logf("PidName(%d) = %q, error = %v", tt.pid, got, err)
		})
	}
}

// BenchmarkListAllPids measures performance of ListAllPids
func BenchmarkListAllPids(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := libproc.ListAllPids()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPidName measures performance of PidName
func BenchmarkPidName(b *testing.B) {
	pid := os.Getpid()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := libproc.PidName(pid)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPidName_Parallel measures parallel performance
func BenchmarkPidName_Parallel(b *testing.B) {
	pid := os.Getpid()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := libproc.PidName(pid)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// TestProcPidInfo tests the general proc_pidinfo wrapper
func TestProcPidInfo(t *testing.T) {
	myPid := os.Getpid()

	t.Run("GetBSDInfo", func(t *testing.T) {
		var info libproc.ProcBSDInfo
		size := int(unsafe.Sizeof(info))

		n, err := libproc.ProcPidInfo(myPid, libproc.PROC_PIDTBSDINFO, 0, unsafe.Pointer(&info), size)
		if err != nil {
			t.Fatalf("ProcPidInfo() error = %v", err)
		}

		if n != size {
			t.Errorf("ProcPidInfo() returned %d bytes, expected %d", n, size)
		}

		if info.Pid != uint32(myPid) {
			t.Errorf("ProcBSDInfo.Pid = %d, expected %d", info.Pid, myPid)
		}

		t.Logf("Process info: pid=%d, ppid=%d, comm=%s", info.Pid, info.Ppid, string(info.Comm[:]))
	})

	t.Run("GetFileDescriptors", func(t *testing.T) {
		// First get the required size
		n, err := libproc.ProcPidInfo(myPid, libproc.PROC_PIDLISTFDS, 0, nil, 0)
		if err != nil {
			t.Fatalf("ProcPidInfo() error = %v", err)
		}

		if n > 0 {
			// Allocate buffer and get FD info
			fds := make([]libproc.ProcFdInfo, n/int(unsafe.Sizeof(libproc.ProcFdInfo{})))
			n2, err := libproc.ProcPidInfo(myPid, libproc.PROC_PIDLISTFDS, 0, unsafe.Pointer(&fds[0]), n)
			if err != nil {
				t.Fatalf("ProcPidInfo() error = %v", err)
			}

			actualCount := n2 / int(unsafe.Sizeof(libproc.ProcFdInfo{}))
			fds = fds[:actualCount]

			t.Logf("Found %d file descriptors", len(fds))
			if len(fds) > 0 {
				t.Logf("First FD: fd=%d, type=%d", fds[0].Fd, fds[0].FdType)
			}
		}
	})

	t.Run("InvalidPid", func(t *testing.T) {
		var info libproc.ProcBSDInfo
		size := int(unsafe.Sizeof(info))

		_, err := libproc.ProcPidInfo(99999, libproc.PROC_PIDTBSDINFO, 0, unsafe.Pointer(&info), size)
		if err == nil {
			t.Logf("ProcPidInfo with invalid PID did not error (may be acceptable)")
		} else {
			t.Logf("ProcPidInfo with invalid PID returned error: %v", err)
		}
	})
}

// TestGetProcBSDInfo tests the higher-level BSD info function
func TestGetProcBSDInfo(t *testing.T) {
	myPid := os.Getpid()

	info, err := libproc.GetProcBSDInfo(myPid)
	if err != nil {
		t.Fatalf("GetProcBSDInfo() error = %v", err)
	}

	if info.Pid != uint32(myPid) {
		t.Errorf("ProcBSDInfo.Pid = %d, expected %d", info.Pid, myPid)
	}

	// Extract process name from Comm field
	var commStr string
	for i, c := range info.Comm {
		if c == 0 {
			commStr = string(info.Comm[:i])
			break
		}
	}

	t.Logf("Process info: pid=%d, ppid=%d, uid=%d, gid=%d, comm=%q",
		info.Pid, info.Ppid, info.Uid, info.Gid, commStr)
}

// TestListPidFileDescriptors tests the FD listing function
func TestListPidFileDescriptors(t *testing.T) {
	myPid := os.Getpid()

	fds, err := libproc.ListPidFileDescriptors(myPid)
	if err != nil {
		t.Fatalf("ListPidFileDescriptors() error = %v", err)
	}

	// A running process should have at least stdin, stdout, stderr
	if len(fds) < 3 {
		t.Logf("Warning: ListPidFileDescriptors() returned only %d FDs (expected at least 3)", len(fds))
	}

	t.Logf("Found %d file descriptors", len(fds))
	for i, fd := range fds {
		if i < 5 { // Log first 5 FDs
			t.Logf("FD %d: fd=%d, type=%d", i, fd.Fd, fd.FdType)
		}
	}
}

// BenchmarkProcPidInfo measures performance of ProcPidInfo
func BenchmarkProcPidInfo(b *testing.B) {
	pid := os.Getpid()
	var info libproc.ProcBSDInfo
	size := int(unsafe.Sizeof(info))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := libproc.ProcPidInfo(pid, libproc.PROC_PIDTBSDINFO, 0, unsafe.Pointer(&info), size)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetProcBSDInfo measures performance of GetProcBSDInfo
func BenchmarkGetProcBSDInfo(b *testing.B) {
	pid := os.Getpid()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := libproc.GetProcBSDInfo(pid)
		if err != nil {
			b.Fatal(err)
		}
	}
}
