package pkg

import (
	"net"
	"time"

	"github.com/chainreactors/utils/httpx"
)

// Socket 复用 utils/httpx 的统一实现，消除与 gogo 的重复定义。
type Socket = httpx.Socket

func NewSocket(network, target string, delay int) (*Socket, error) {
	return NewSocketWithDialer(network, target, delay, nil)
}

// NewSocketWithDialer 使用指定 dialTimeout 创建 Socket。dialTimeout 为 nil 时
// 直连（net.DialTimeout）——不再依赖任何包级全局，天然并发安全。
func NewSocketWithDialer(network, target string, delay int, dialTimeout func(string, string, time.Duration) (net.Conn, error)) (*Socket, error) {
	return httpx.NewSocket(network, target, httpx.SocketConfig{
		Timeout:     time.Duration(delay) * time.Second,
		DialTimeout: dialTimeout,
	})
}
