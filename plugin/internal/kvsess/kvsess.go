package kvsess

import (
	"github.com/go-redis/redis"
)

type RedisSession struct {
	Client  *redis.Client
	SvcName string
}

func (s *RedisSession) Service() string  { return s.SvcName }
func (s *RedisSession) Close() error     { return s.Client.Close() }
func (s *RedisSession) Raw() interface{} { return s.Client }

func (s *RedisSession) Get(key string) ([]byte, error) {
	val, err := s.Client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (s *RedisSession) Keys(pattern string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	for {
		keys, next, err := s.Client.Scan(cursor, pattern, 100).Result()
		if err != nil {
			return allKeys, err
		}
		allKeys = append(allKeys, keys...)
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return allKeys, nil
}
