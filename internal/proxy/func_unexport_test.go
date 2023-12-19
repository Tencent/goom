// Package proxy_test 对 proxy 包的测试
package proxy_test

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/tencent/goom/internal/proxy"
)

// TestNetConnMock 测试 mock 网络连接
func TestNetConnMock(t *testing.T) {
	// 原始函数
	var connWrite func(c *conn, b []byte) (int, error)

	// 使用 patch 进行切面
	patch, err := proxy.FuncName("net.(*conn).Write", func(c *conn, b []byte) (int, error) {
		n, _ := connWrite(c, b)
		// 修改返回结果
		return n, errors.New("mocked")
	}, &connWrite)

	if err != nil {
		t.Error("mock print err:", err)
	}

	// 发起网络请求
	host := "127.0.0.1"
	port := 80

	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	fmt.Println("Connecting to " + host + ":" + strconv.Itoa(port))

	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	content := []byte{1, 2, 3}
	_, err = conn.Write(content)

	// 预期返回: err: mocked
	t.Log("err:", err)
	patch.Unpatch()
}

// // nolint
type conn struct {
	// nolint
	fd *netFD
}

// nolint
func (c *conn) Read([]byte) (n int, err error) { return 0, nil }

// nolint
func (c *conn) Write([]byte) (n int, err error) { return 0, nil }

// nolint
func (c *conn) Close() error { return nil }

// nolint
func (c *conn) LocalAddr() net.Addr { return nil }

// nolint
func (c *conn) RemoteAddr() net.Addr { return nil }

// nolint
func (c *conn) SetDeadline(time.Time) error { return nil }

// nolint
func (c *conn) SetReadDeadline(time.Time) error { return nil }

// nolint
func (c *conn) SetWriteDeadline(time.Time) error { return nil }

// nolint
// Network file descriptor.
type netFD struct {
	// nolint
	pfd FD

	// immutable until Close
	family int
	// nolint
	sotype      int
	isConnected bool // handshake completed or use of association with peer
	net         string
	laddr       net.Addr
	raddr       net.Addr
}

// nolint
// FD is a file descriptor. The net and os packages use this type as a
// field of a larger type representing a network connection or OS file.
type FD struct {
	// Lock sysfd and serialize access to Read and Write methods.
	fdmu fdMutex

	// System file descriptor. Immutable until Close.
	Sysfd int

	// I/O poller.
	pd pollDesc

	// Writev cache.
	iovecs *[]int

	// Semaphore signaled when file is closed.
	csema uint32

	// Non-zero if this file has been set to blocking mode.
	isBlocking uint32

	// Whether this is a streaming descriptor, as opposed to a
	// packet-based descriptor like a UDP socket. Immutable.
	IsStream bool

	// Whether a zero byte read indicates EOF. This is false for a
	// message based socket connection.
	ZeroReadIsEOF bool

	// Whether this is a file rather than a network socket.
	isFile bool
}

// nolint
type fdMutex struct {
	state uint64
	rsema uint32
	wsema uint32
}

// nolint
type pollDesc struct {
	runtimeCtx uintptr
}
