package memcache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/chainreactors/zombie/pkg"
)

type MemcachePlugin struct {
	*pkg.Task
	client *memcache.Client
}

func (s *MemcachePlugin) Name() string {
	return s.Service
}

func (s *MemcachePlugin) Unauth() (bool, error) {
	client := memcache.New(fmt.Sprintf("%s:%s", s.IP, s.Port))
	s.client = client
	return true, nil
}

func (s *MemcachePlugin) Login() error {
	client := memcache.New(fmt.Sprintf("%s:%s", s.IP, s.Port))
	s.client = client
	// Memcache doesn't support authentication by default
	return nil
}

func (s *MemcachePlugin) GetResult() *pkg.Result {
	// todo list items
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *MemcachePlugin) Close() error {
	// Memcache client doesn't have a close method
	return nil
}
