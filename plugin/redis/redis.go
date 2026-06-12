package redis

import (
	"net"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin/internal/kvsess"
	"github.com/go-redis/redis"
)

// RedisPlugin is a stateless factory that satisfies the Plugin interface.
type RedisPlugin struct{}

func (RedisPlugin) Name() string { return "redis" }

// Open authenticates with the password from task and returns a KVSession.
func (RedisPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	return dial(task, task.Password)
}

// Unauth attempts an unauthenticated connection (empty password).
func (RedisPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return dial(task, "")
}

// dial builds redis.Options (with optional proxy Dialer), connects, pings,
// and wraps the client in a kvsess.RedisSession.
func dial(task *pkg.Task, password string) (pkg.Session, error) {
	opt := &redis.Options{
		Addr:        task.Address(),
		Password:    password,
		DB:          0,
		DialTimeout: task.Duration(),
	}
	if task.ProxyDial != nil {
		opt.Dialer = func() (net.Conn, error) {
			return task.DialTimeout("tcp", task.Address(), task.Duration())
		}
	}

	client := redis.NewClient(opt)
	if _, err := client.Ping().Result(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &kvsess.RedisSession{
		Client:  client,
		SvcName: "redis",
	}, nil
}
