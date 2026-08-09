package memcache

import (
	"fmt"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/chainreactors/zombie/pkg"
)

type memcacheSession struct {
	service string
	client  *memcache.Client
}

func (s *memcacheSession) Service() string { return s.service }

func (s *memcacheSession) Close() error { return nil }

func (s *memcacheSession) Get(key string) ([]byte, error) {
	item, err := s.client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (s *memcacheSession) Keys(pattern string) ([]string, error) {
	return nil, fmt.Errorf("memcached does not support key enumeration")
}

func (s *memcacheSession) Command(name string, args ...string) (interface{}, error) {
	switch strings.ToUpper(name) {
	case "SET":
		if len(args) < 2 {
			return nil, fmt.Errorf("SET requires key and value")
		}
		return "OK", s.client.Set(&memcache.Item{Key: args[0], Value: []byte(args[1])})
	case "DELETE":
		if len(args) < 1 {
			return nil, fmt.Errorf("DELETE requires key")
		}
		return "OK", s.client.Delete(args[0])
	case "FLUSH", "FLUSH_ALL":
		return "OK", s.client.FlushAll()
	default:
		return nil, fmt.Errorf("unsupported memcached command: %s", name)
	}
}

type MemcachePlugin struct{}

func (p *MemcachePlugin) Open(task *pkg.Task) (pkg.Session, error) {
	client := memcache.New(fmt.Sprintf("%s:%s", task.IP, task.Port))
	return &memcacheSession{service: task.Service, client: client}, nil
}

func (p *MemcachePlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	client := memcache.New(fmt.Sprintf("%s:%s", task.IP, task.Port))
	return &memcacheSession{service: task.Service, client: client}, nil
}
