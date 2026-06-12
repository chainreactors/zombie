package pkg

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/chainreactors/parsers"
)

func TestClassifyError_Timeout(t *testing.T) {
	cases := []error{
		TimeoutError{err: errors.New("dial"), timeout: 5, service: "ssh"},
		fmt.Errorf("read tcp: i/o timeout"),
		fmt.Errorf("context deadline exceeded"),
	}
	for _, err := range cases {
		if got := ClassifyError(err); got != ErrCatTimeout {
			t.Errorf("ClassifyError(%q) = %d, want ErrCatTimeout", err, got)
		}
	}
}

func TestClassifyError_Refused(t *testing.T) {
	cases := []error{
		fmt.Errorf("dial tcp 127.0.0.1:22: connection refused"),
		fmt.Errorf("connection reset by peer"),
		fmt.Errorf("connect: no route to host"),
	}
	for _, err := range cases {
		if got := ClassifyError(err); got != ErrCatRefused {
			t.Errorf("ClassifyError(%q) = %d, want ErrCatRefused", err, got)
		}
	}
}

func TestClassifyError_Auth(t *testing.T) {
	cases := []error{
		ErrorWrongUserOrPwd,
		fmt.Errorf("ssh: unable to authenticate"),
		fmt.Errorf("Access denied for user 'root'"),
		fmt.Errorf("authentication failed"),
	}
	for _, err := range cases {
		if got := ClassifyError(err); got != ErrCatAuth {
			t.Errorf("ClassifyError(%q) = %d, want ErrCatAuth", err, got)
		}
	}
}

func TestClassifyError_WrappedTimeout(t *testing.T) {
	inner := TimeoutError{err: errors.New("dial"), timeout: 5, service: "ssh"}
	wrapped := fmt.Errorf("open failed: %w", inner)
	if got := ClassifyError(wrapped); got != ErrCatTimeout {
		t.Errorf("ClassifyError(wrapped TimeoutError) = %d, want ErrCatTimeout", got)
	}
}

func TestClassifyError_NetOpError(t *testing.T) {
	err := &net.OpError{Op: "dial", Net: "tcp", Err: fmt.Errorf("connection refused")}
	if got := ClassifyError(err); got != ErrCatRefused {
		t.Errorf("ClassifyError(net.OpError) = %d, want ErrCatRefused", got)
	}
}

func TestClassifyError_Other(t *testing.T) {
	if got := ClassifyError(errors.New("something unexpected")); got != ErrCatOther {
		t.Errorf("ClassifyError(unknown) = %d, want ErrCatOther", got)
	}
}

func TestClassifyError_Nil(t *testing.T) {
	if got := ClassifyError(nil); got != ErrCatOther {
		t.Errorf("ClassifyError(nil) = %d, want ErrCatOther", got)
	}
}

func TestStatistor_RecordError(t *testing.T) {
	stat := &Statistor{Tasks: make(map[string]int)}

	stat.RecordError(fmt.Errorf("i/o timeout"))
	stat.RecordError(fmt.Errorf("i/o timeout"))
	stat.RecordError(fmt.Errorf("connection refused"))
	stat.RecordError(ErrorWrongUserOrPwd)
	stat.RecordError(errors.New("random"))
	stat.RecordError(nil)

	if stat.ErrTimeout != 2 {
		t.Errorf("ErrTimeout = %d, want 2", stat.ErrTimeout)
	}
	if stat.ErrRefused != 1 {
		t.Errorf("ErrRefused = %d, want 1", stat.ErrRefused)
	}
	if stat.ErrAuth != 1 {
		t.Errorf("ErrAuth = %d, want 1", stat.ErrAuth)
	}
	if stat.ErrOther != 1 {
		t.Errorf("ErrOther = %d, want 1", stat.ErrOther)
	}
}

func TestStatistor_RecordResult(t *testing.T) {
	stat := &Statistor{Tasks: make(map[string]int)}
	task := &Task{ZombieResult: &parsers.ZombieResult{IP: "1.1.1.1", Port: "22", Service: "ssh"}}

	stat.RecordResult(&Result{Task: task, OK: true, Extracteds: make(parsers.Extracteds, 3), Loot: map[string][]byte{"a": {}, "b": {}}})
	stat.RecordResult(&Result{Task: task, OK: true})
	stat.RecordResult(&Result{Task: task, OK: false, Err: fmt.Errorf("connection refused")})

	if stat.Success != 2 {
		t.Errorf("Success = %d, want 2", stat.Success)
	}
	if stat.Extracteds != 3 {
		t.Errorf("Extracteds = %d, want 3", stat.Extracteds)
	}
	if stat.Loot != 2 {
		t.Errorf("Loot = %d, want 2", stat.Loot)
	}
	if stat.ErrRefused != 1 {
		t.Errorf("ErrRefused = %d, want 1", stat.ErrRefused)
	}
}

func TestStatistor_ErrorString(t *testing.T) {
	stat := &Statistor{Tasks: make(map[string]int)}
	if s := stat.ErrorString(); s != "" {
		t.Errorf("empty stat should return empty string, got %q", s)
	}

	stat.ErrTimeout = 10
	stat.ErrRefused = 5
	stat.ErrAuth = 20
	stat.ErrOther = 3
	s := stat.ErrorString()
	if s != "errors: timeout=10, refused=5, auth_fail=20, other=3" {
		t.Errorf("unexpected ErrorString: %q", s)
	}
}

func TestStatistor_SummaryString(t *testing.T) {
	stat := &Statistor{Tasks: make(map[string]int), Total: 100, Success: 3}
	s := stat.SummaryString()
	if s != "total: 100, success: 3" {
		t.Errorf("basic summary = %q", s)
	}

	stat.Extracteds = 5
	stat.Loot = 2
	stat.ErrTimeout = 10
	s = stat.SummaryString()
	expect := "total: 100, success: 3, extracteds: 5, loot: 2, errors: timeout=10, refused=0, auth_fail=0, other=0"
	if s != expect {
		t.Errorf("full summary = %q, want %q", s, expect)
	}
}
