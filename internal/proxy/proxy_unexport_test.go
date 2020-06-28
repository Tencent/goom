package proxy

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func TestPrintMock(t *testing.T) {
	var trampoline = func(a ...interface{}) (n int, err error) {
		return 0, nil
	}

	// 静态代理函数
	patch, err := StaticProxyByName("fmt.Print", func(a ...interface{}) (n int, err error) {
		// 调用原来的函数
		return fmt.Println("called fmt.Print, args:", a)
	}, &trampoline)
	if err != nil {
		t.Error("mock print err:", err)
	}

	fmt.Print("ok", "1")
	patch.Unpatch()
	fmt.Println("unpatched")
	fmt.Print("ok", "2")
}

func TestNetConnMock(t *testing.T) {

	// 原始函数
	var connWrite func(c *conn, b []byte) (int, error)

	// 使用gomonkey进行切面
	patch, err := StaticProxyByName("net.(*conn).Write", func(c *conn, b []byte) (int, error) {
		n, _ := connWrite(c, b)
		// 修改返回结果
		return n, errors.New("mocked")
	}, &connWrite)

	if err != nil {
		t.Error("mock print err:", err)
	}

	// 发起网络请求
	host := "100.65.4.24"
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

type conn struct {
	fd *netFD
}

func (c *conn) Read(b []byte) (n int, err error) { return 0, nil }

func (c *conn) Write(b []byte) (n int, err error) { return 0, nil }

func (c *conn) Close() error { return nil }

func (c *conn) LocalAddr() net.Addr { return nil }

func (c *conn) RemoteAddr() net.Addr { return nil }

func (c *conn) SetDeadline(t time.Time) error { return nil }

func (c *conn) SetReadDeadline(t time.Time) error { return nil }

func (c *conn) SetWriteDeadline(t time.Time) error { return nil }

func (c *conn) ok() bool { return c != nil && c.fd != nil }

func connIdentity(c conn) string {
	fd := c.fd
	ptr := (uintptr)(unsafe.Pointer(fd))
	return strconv.FormatInt(int64(ptr), 10)
}

// Network file descriptor.
type netFD struct {
	pfd FD

	// immutable until Close
	family      int
	sotype      int
	isConnected bool // handshake completed or use of association with peer
	net         string
	laddr       net.Addr
	raddr       net.Addr
}

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
	iovecs *[]syscall.Iovec

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

type fdMutex struct {
	state uint64
	rsema uint32
	wsema uint32
}

type pollDesc struct {
	runtimeCtx uintptr
}
