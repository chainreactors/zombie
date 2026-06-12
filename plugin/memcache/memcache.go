package memcache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/chainreactors/zombie/pkg"
)

// memcacheSession implements pkg.Session over a memcache client.
type memcacheSession struct {
	service string
	client  *memcache.Client
}

func (s *memcacheSession) Service() string  { return s.service }
func (s *memcacheSession) Raw() interface{} { return s.client }

func (s *memcacheSession) Close() error {
	// Memcache client doesn't have a close method
	return nil
}

// MemcachePlugin is stateless; all connection state lives in memcacheSession.
type MemcachePlugin struct{}

func (p *MemcachePlugin) Name() string { return "memcached" }

func (p *MemcachePlugin) Open(task *pkg.Task) (pkg.Session, error) {
	client := memcache.New(fmt.Sprintf("%s:%s", task.IP, task.Port))
	// Memcache doesn't support authentication by default
	return &memcacheSession{service: task.Service, client: client}, nil
}

func (p *MemcachePlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	// Memcache has no auth, so unauth always returns a session
	client := memcache.New(fmt.Sprintf("%s:%s", task.IP, task.Port))
	return &memcacheSession{service: task.Service, client: client}, nil
}
