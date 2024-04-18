package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chainreactors/fingers/common"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/parsers"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	InterruptError = errors.New("interrupt")
)

type NilConnError struct {
	Service string
}

func (e NilConnError) Error() string {
	return e.Service + " has nil conn"
}

type TimeoutError struct {
	err     error
	timeout int
	service string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("%s spended out of %ds, %s", e.service, e.timeout, e.err.Error())
}

func (e TimeoutError) Unwrap() error { return e.err }

type TaskMod int

const (
	TaskModBrute TaskMod = 0 + iota
	TaskModUnauth
	TaskModCheck
	TaskModSniper
)

func (m TaskMod) String() string {
	switch m {
	case TaskModBrute:
		return "brute"
	case TaskModUnauth:
		return "unauth"
	case TaskModCheck:
		return "check"
	case TaskModSniper:
		return "sniper"
	default:
		return "unknown"
	}
}

type Task struct {
	IP       string             `json:"ip"`
	Port     string             `json:"port"`
	Service  string             `json:"service"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Scheme   string             `json:"scheme"`
	Param    map[string]string  `json:"-"`
	Mod      TaskMod            `json:"-"`
	Timeout  int                `json:"-"`
	Context  context.Context    `json:"-"`
	Canceler context.CancelFunc `json:"-"`
	Locker   *sync.Mutex        `json:"-"`
}

func (t *Task) String() string {
	return fmt.Sprintf("%s://%s:%s", t.Service, t.IP, t.Port)
}

func (t *Task) Address() string {
	return t.IP + ":" + t.Port
}

func (t *Task) URI() string {
	if t.Scheme != "" {
		return t.Scheme + "://" + t.Address()
	} else {
		return t.Service + "://" + t.Address()
	}
}

func (t *Task) URL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s", t.Scheme, t.Username, t.Password, t.IP, t.Port)
}

func (t *Task) UintPort() uint16 {
	p, _ := strconv.Atoi(t.Port)
	return uint16(p)
}

func (t *Task) Duration() time.Duration {
	return time.Duration(t.Timeout) * time.Second
}

func NewResult(task *Task, err error) *Result {
	if err != nil {
		return &Result{
			Task: task,
			OK:   false,
			Err:  err,
		}
	} else {
		return &Result{
			Task: task,
			OK:   true,
		}
	}
}

type Result struct {
	*Task
	Vulns      common.Vulns
	Extracteds parsers.Extracteds
	OK         bool
	Err        error
}

func (r *Result) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("[%s] ", r.Mod.String()))
	s.WriteString(r.URI())
	if r.Username != "" {
		s.WriteString(" " + r.Username)
	}
	if r.Password != "" {
		s.WriteString(" " + r.Password)
	}
	if len(r.Param) != 0 {
		s.WriteString(" " + fmt.Sprintf("%v", r.Param))
	}

	s.WriteString(", " + r.Service + " login successfully\n")
	return s.String()
}

func (r *Result) Json() string {
	bs, err := json.Marshal(r)
	if err != nil {
		logs.Log.Error(err.Error())
		return ""
	}
	return string(bs) + "\n"
}

func (r *Result) Format(form string) string {
	switch form {
	case "json":
		return r.Json()
	case "csv":
		return ""
	default:
		return r.String()
	}
}

type Basic struct {
	Input string
	Data  string
}
