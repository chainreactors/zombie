package plugin

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/go-redis/redis"
)

type RedisService struct {
	*pkg.Task
	conn       *redis.Client
	Additional string
	Input      string
}

func (s *RedisService) Query() bool {
	return false
}

func (s *RedisService) GetInfo() bool {
	return false
}

func (s *RedisService) Connect() error {
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

func (s *RedisService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return pkg.NilConnError{s.Service}
}

func (s *RedisService) SetQuery(query string) {
	s.Input = query
}

func (s *RedisService) Output(res interface{}) {

}
