/*
   Package to allow go applications to immediately start
   listening on a socket, unix, tcp, udp but hold connections
   until the application has booted and is ready to accept them
*/
package listenbuffer

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"reflect"
	"syscall"
)

// NewListenBuffer returns a listener listening on addr with the protocol.
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

type defaultListener struct {
	wrapped  net.Listener // the real listener to wrap
	ready    bool         // is the listner ready to start accpeting connections
	activate chan struct{}
}

func (l *defaultListener) Close() error {
	return l.wrapped.Close()
}

func (l *defaultListener) Addr() net.Addr {
	return l.wrapped.Addr()
}

type CredConn struct {
	net.Conn
	Cred *syscall.Ucred
}

func (l *defaultListener) Accept() (net.Conn, error) {
	// if the listen has been told it is ready then we can go ahead and
	// start returning connections
	if l.ready {
		conn, err := l.wrapped.Accept()
		if xconn, ok := &conn.(*net.UnixConn); ok {
			fd := int(reflect.ValueOf(xconn).Elem().Elem().Elem().FieldByName("conn").FieldByName("fd").Elem().FieldByName("sysfd").Int())
			if ucred, err := syscall.GetsockoptUcred(fd, syscall.SOL_SOCKET, syscall.SO_PEERCRED); err == nil {
				log.Printf("uid: %d, gid: %d, pid: %d", ucred.Uid, ucred.Gid, ucred.Pid)
				return CredConn{conn, ucred}, nil
			} else {
				log.Printf("Error: %s", err.Error())
			}
		}
		return conn, err
	}
	<-l.activate
	l.ready = true
	return l.Accept()
}
