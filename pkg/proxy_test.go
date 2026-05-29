package pkg

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chainreactors/parsers"
)

// TestTaskHTTPClientProxy 验证 http 系插件统一使用的 task.HTTPClient 会经过 ProxyDial。
func TestTaskHTTPClientProxy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	var hits int32
	task := &Task{
		ZombieResult: &parsers.ZombieResult{},
		Context:      context.Background(),
		Timeout:      5,
		ProxyDial: func(ctx context.Context, network, address string) (net.Conn, error) {
			atomic.AddInt32(&hits, 1)
			return (&net.Dialer{}).DialContext(ctx, network, address)
		},
	}
	resp, err := task.HTTPClient(true).Get(srv.URL)
	if err != nil {
		t.Fatalf("http client get: %v", err)
	}
	resp.Body.Close()
	if atomic.LoadInt32(&hits) != 1 {
		t.Fatalf("task.HTTPClient did not route through ProxyDial, hits=%d", hits)
	}
}

// TestTaskHTTPClientDirect 验证未设置 ProxyDial 时 task.HTTPClient 直连。
func TestTaskHTTPClientDirect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	task := &Task{ZombieResult: &parsers.ZombieResult{}, Context: context.Background(), Timeout: 5}
	resp, err := task.HTTPClient(true).Get(srv.URL)
	if err != nil {
		t.Fatalf("direct http get: %v", err)
	}
	resp.Body.Close()
}

// TestTaskDialTimeoutDirect 验证未设置 ProxyDial 时走直连。
func TestTaskDialTimeoutDirect(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		c.Close()
	}()

	task := &Task{ZombieResult: &parsers.ZombieResult{}}
	conn, err := task.DialTimeout("tcp", ln.Addr().String(), 2*time.Second)
	if err != nil {
		t.Fatalf("direct dial: %v", err)
	}
	conn.Close()
}

// TestTaskDialTimeoutProxy 验证设置 ProxyDial 后连接经过代理拨号器。
func TestTaskDialTimeoutProxy(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()

	var hits int32
	task := &Task{
		ZombieResult: &parsers.ZombieResult{},
		Context:      context.Background(),
		ProxyDial: func(ctx context.Context, network, address string) (net.Conn, error) {
			atomic.AddInt32(&hits, 1)
			d := net.Dialer{}
			return d.DialContext(ctx, network, address)
		},
	}
	conn, err := task.DialTimeout("tcp", ln.Addr().String(), 2*time.Second)
	if err != nil {
		t.Fatalf("proxy dial: %v", err)
	}
	conn.Close()
	if atomic.LoadInt32(&hits) != 1 {
		t.Fatalf("expected proxy dialer to be used once, hits=%d", hits)
	}
}

// TestNewSocketWithDialer 验证 socket 风格的拨号器注入。
func TestNewSocketWithDialer(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		io.Copy(io.Discard, c)
		c.Close()
	}()

	var hits int32
	dial := func(network, address string, timeout time.Duration) (net.Conn, error) {
		atomic.AddInt32(&hits, 1)
		return net.DialTimeout(network, address, timeout)
	}
	s, err := NewSocketWithDialer("tcp", ln.Addr().String(), 2, dial)
	if err != nil {
		t.Fatalf("NewSocketWithDialer: %v", err)
	}
	s.Close()
	if atomic.LoadInt32(&hits) != 1 {
		t.Fatalf("expected dialer used once, hits=%d", hits)
	}
}
