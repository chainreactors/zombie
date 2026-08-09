package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chainreactors/utils/httpx"
	"github.com/chainreactors/utils/parsers"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	InterruptError      = errors.New("interrupt")
	ErrorWrongUserOrPwd = errors.New("wrong username or password")
	NotImplUnauthorized = errors.New("not implemented unauthorized")
)

type TimeoutError struct {
	err     error
	timeout int
	service string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("%s spended out of %ds, %s", e.service, e.timeout, e.err.Error())
}

func (e TimeoutError) Unwrap() error { return e.err }

var UnknownService = &Service{Name: "unknown", DefaultPort: "", Source: "unknown"}

var Services = services{
	Plugins: map[string]*Service{},
	Aliases: map[string]*Service{},
}

type services struct {
	mu      sync.RWMutex
	Plugins map[string]*Service
	Aliases map[string]*Service
}

func (ss *services) Get(name string) (*Service, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	if s, ok := ss.Plugins[name]; ok {
		return s, true
	}
	if s, ok := ss.Aliases[name]; ok {
		return s, true
	}
	return UnknownService, false
}

func (ss *services) Register(s *Service) bool {
	if s == nil {
		return false
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if _, ok := ss.Plugins[s.Name]; !ok {
		ss.Plugins[s.Name] = s
	}
	for _, a := range s.Alias {
		if _, ok := ss.Aliases[a]; !ok {
			ss.Aliases[a] = s
		}
	}
	return true
}

// All returns a snapshot of the registered services.
func (ss *services) All() map[string]*Service {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	services := make(map[string]*Service, len(ss.Plugins))
	for name, service := range ss.Plugins {
		services[name] = service
	}
	return services
}

func (ss *services) DefaultPort(service string) string {
	if s, ok := ss.Get(service); ok {
		return s.DefaultPort
	}
	return ""
}

// SupportedServiceNames 返回所有已注册服务名(已排序),供未知服务的友好报错使用。
func SupportedServiceNames() string {
	services := Services.All()
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

const (
	PluginSource  = "plugin"
	NeutronSource = "neutron"
)

type Service struct {
	Name        string
	Alias       []string
	DefaultPort string
	Source      string
}

func (s Service) String() string {
	return s.Name
}

func GetDefault(port string) string {
	for _, s := range Services.All() {
		if s.DefaultPort == port {
			return s.Name
		}
	}
	return UnknownService.Name
}

// DialFunc 是与 proxyclient.Dial 兼容的拨号函数签名。使用普通函数类型而非
// 直接依赖 proxyclient，避免给 zombie 引入更高的 Go 版本要求（proxyclient
// 需要 go1.24，而 zombie 仍为 go1.16）。SDK 层可直接把 proxyclient.Dial /
// dialer.DialContext 赋值给该字段。
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)

// DialTimeoutFunc 与 NewSocketWithDialer / Task.DialTimeout 的签名一致，
// 供 socket 风格的插件（如 rsync）传递代理拨号器。
type DialTimeoutFunc func(network, address string, timeout time.Duration) (net.Conn, error)

type Task struct {
	*parsers.ZombieResult
	Timeout   int                `json:"-"`
	Context   context.Context    `json:"-"`
	Cancel    context.CancelFunc `json:"-"`
	Completed chan struct{}      `json:"-"`
	Raw       bool               `json:"-"`
	// ProxyDial 非 nil 时，插件应使用它建立连接而非直接 net.Dial。
	ProxyDial DialFunc `json:"-"`
}

func (t *Task) Duration() time.Duration {
	return time.Duration(t.Timeout) * time.Second
}

// DialTimeout 按 task 配置建立连接：设置了 ProxyDial 则走代理，否则直连。
// network 通常为 "tcp"。
func (t *Task) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	if t.ProxyDial != nil {
		ctx := t.Context
		if ctx == nil {
			ctx = context.Background()
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return t.ProxyDial(ctx, network, address)
	}
	return net.DialTimeout(network, address, timeout)
}

// HTTPClient 返回一个 per-task 的 *http.Client（统一经 utils/httpx 构造，零全局）。
// 设置了 ProxyDial 时连接走代理；否则直连。所有 http 系插件应使用它，
// 替代 http.DefaultClient，以保证代理生效且并发隔离。
func (t *Task) HTTPClient(followRedirects bool) *http.Client {
	cfg := httpx.ClientConfig{
		Timeout:            t.Duration(),
		FollowRedirects:    followRedirects,
		InsecureSkipVerify: true,
	}
	if t.ProxyDial != nil {
		cfg.DialContext = httpx.DialContextFunc(t.ProxyDial)
	}
	return httpx.NewHTTPClient(cfg)
}

func NewResult(task *Task, err error) *Result {
	result := &Result{Task: task, Err: err}
	if task == nil || task.ZombieResult == nil {
		return result
	}
	task.OK = err == nil
	if err != nil {
		task.ErrString = err.Error()
	} else {
		task.ErrString = ""
	}
	return result
}

type Result struct {
	*Task         `json:",inline"`
	Err           error           `json:"-"`
	ActionResults []*ActionResult `json:"-"`
}

func (r *Result) Merge(ar *ActionResult) {
	if ar == nil {
		return
	}
	r.Extracteds = append(r.Extracteds, ar.Extracteds...)
	for k, v := range ar.Vulns {
		if r.Vulns == nil {
			r.Vulns = make(parsers.Vulns)
		}
		r.Vulns[k] = v
	}
	for k, v := range ar.Loot {
		if r.Loot == nil {
			r.Loot = map[string][]byte{}
		}
		r.Loot[k] = v
	}
	r.ActionResults = append(r.ActionResults, ar)
}

func (r *Result) Format(form string) string {
	if r == nil || r.Task == nil || r.ZombieResult == nil {
		return ""
	}
	switch form {
	case parsers.ZombieFormatJSON, parsers.ZombieFormatJSONLine:
		bs, err := json.Marshal(r)
		if err != nil {
			return ""
		}
		return string(bs) + "\n"
	default:
		out := r.ZombieResult.Format(form)
		if len(r.Extracteds) == 0 {
			return out
		}
		return strings.TrimRight(out, "\n") + " " + r.Extracteds.String()
	}
}

func ParseMethod(input string, raw bool) (string, string) {
	if raw {
		return "", input
	}
	if strings.HasPrefix(input, "pk:") {
		return "pk", input[3:]
	} else if strings.HasPrefix(input, "hash:") {
		return "hash", input[5:]
	} else if strings.HasPrefix(input, "raw:") {
		return "raw", input[4:]
	} else {
		return "", input
	}
}
