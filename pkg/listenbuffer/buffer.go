/*
listenbuffer uses the kernel's listening backlog functionality to queue
connections, allowing applications to start listening immediately and handle
connections later. This is signaled by closing the activation channel passed to
the constructor.

The maximum amount of queued connections depends on the configuration of your
kernel (typically called SOMAXXCON) and cannot be configured in Go with the
net package. See `src/net/sock_platform.go` in the Go tree or consult your
kernel's manual.

	activator := make(chan struct{})
	buffer, err := NewListenBuffer("tcp", "localhost:4000", activator)
	if err != nil {
		panic(err)
	}

	// will block until activator has been closed or is sent an event
	client, err := buffer.Accept()

Somewhere else in your application once it's been booted:

	close(activator)

`buffer.Accept()` will return the first client in the kernel listening queue, or
continue to block until a client connects or an error occurs.
*/
package listenbuffer

import (
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"strconv"
	"syscall"
)

// NewListenBuffer returns a net.Listener listening on addr with the protocol
// passed. The channel passed is used to activate the listenbuffer when the
// caller is ready to accept connections.
func NewListenBuffer(proto, addr string, activate chan struct{}) (net.Listener, error) {
	wrapped, err := net.Listen(proto, addr)
	if err != nil {
		return nil, err
	}

	return &defaultListener{
		wrapped:  wrapped,
		activate: activate,
	}, nil
}

// defaultListener is the buffered wrapper around the net.Listener
type defaultListener struct {
	wrapped  net.Listener  // The net.Listener wrapped by listenbuffer
	ready    bool          // Whether the listenbuffer has been activated
	activate chan struct{} // Channel to control activation of the listenbuffer
}

// Close closes the wrapped socket.
func (l *defaultListener) Close() error {
	return l.wrapped.Close()
}

// Addr returns the listening address of the wrapped socket.
func (l *defaultListener) Addr() net.Addr {
	return l.wrapped.Addr()
}

type CredConn struct {
	net.Conn
	Cred     *syscall.Ucred
	Loginuid int
}

func newCredConn(conn net.Conn) (net.Conn, error) {
	if _, ok := conn.(*net.UnixConn); ok {
		fd := int(reflect.ValueOf(&conn).Elem().Elem().Elem().FieldByName("conn").FieldByName("fd").Elem().FieldByName("sysfd").Int())
		if ucred, err := syscall.GetsockoptUcred(fd, syscall.SOL_SOCKET, syscall.SO_PEERCRED); err == nil {
			if data, err := ioutil.ReadAll(fmt.Sprintf("/proc/%d/loginuid", ucred.Pid)); err == nil {
				if luid, err := strconv.Atoi(string(data)); err == nil {
					return CredConn{conn, ucred, luid}, nil
				} else {
					return nil, fmt.Errorf("Invalid loginuid %q", string(data))
				}
			} else {
				return nil, fmt.Errorf("Unable to open `/proc/%d/loginuid`", ucred.Pid)
			}
		} else {
			return nil, err
		}
	}

	return conn, nil
}

// Accept returns a client connection on the wrapped socket if the listen buffer
// has been activated. To active the listenbuffer the activation channel passed
// to NewListenBuffer must have been closed or sent an event.
func (l *defaultListener) Accept() (net.Conn, error) {
	// if the listen has been told it is ready then we can go ahead and
	// start returning connections
	if l.ready {
		conn, err := l.wrapped.Accept()
		if err == nil {
			return newCredConn(conn)
		}
		return conn, err
	}
	<-l.activate
	l.ready = true
	return l.Accept()
}
