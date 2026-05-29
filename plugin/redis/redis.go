package redis

import (
	"net"

	"github.com/chainreactors/zombie/pkg"
	"github.com/go-redis/redis"
)

type RedisPlugin struct {
	*pkg.Task
	conn       *redis.Client
	Additional string
	Input      string
}

// options 构建 redis 连接参数，并在配置了代理时注入自定义 Dialer。
func (s *RedisPlugin) options(password string) *redis.Options {
	opt := &redis.Options{
		Addr:        s.Address(),
		Password:    password,
		DB:          0,
		DialTimeout: s.Duration(),
	}
	if s.ProxyDial != nil {
		opt.Dialer = func() (net.Conn, error) {
			return s.DialTimeout("tcp", s.Address(), s.Duration())
		}
	}
	return opt
}

func (s *RedisPlugin) Login() error {
	client := redis.NewClient(s.options(s.Password))
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	s.conn = client
	return nil

}

func (s *RedisPlugin) Unauth() (bool, error) {
	client := redis.NewClient(s.options(""))
	_, err := client.Ping().Result()
	if err != nil {
		return false, err
	}

	s.conn = client
	return true, nil
}

func (s *RedisPlugin) Name() string {
	return s.Service
}

func (s *RedisPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *RedisPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
