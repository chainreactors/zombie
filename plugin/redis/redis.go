package redis

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/go-redis/redis"
)

type RedisPlugin struct {
	*pkg.Task
	conn       *redis.Client
	Additional string
	Input      string
}

func (s *RedisPlugin) Login() error {
	opt := redis.Options{Addr: s.Address(),
		Password: s.Password, DB: 0, DialTimeout: s.Duration()}
	client := redis.NewClient(&opt)
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	s.conn = client
	return nil

}

func (s *RedisPlugin) Unauth() (bool, error) {
	opt := redis.Options{Addr: s.Address(),
		Password: "", DB: 0, DialTimeout: s.Duration()}
	client := redis.NewClient(&opt)
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
