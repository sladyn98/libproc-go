package libproc_test

import (
	"os"
	"testing"

	"github.com/sladyn98/libproc-go"
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

	// Parent process name should not be empty
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
