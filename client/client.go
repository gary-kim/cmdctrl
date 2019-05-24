package client

import (
	"crypto/tls"
	"net"

	"github.com/gorilla/websocket"
)

// RemoteServer represents a single cmdctrl server that this client is configured to connect to
type RemoteServer struct {
	addr     string
	RESTMode bool
	conn     *net.Conn
	wsconn   *websocket.Conn
}

// Options represents the options given by the user when cmdctrl was started in client mode
type Options struct {
	RESTMode bool
}

var primaryServer RemoteServer

// RunClient runs cmdctrl in client mode
func RunClient(addr string, opt Options) {
	primaryServer = RemoteServer{
		addr:     addr,
		RESTMode: opt.RESTMode,
	}
}

// dial makes a REST mode connection to a cmdctrl server
func (s RemoteServer) dial() (net.Conn, error) {
	TLSConfig := &tls.Config{
		ServerName: s.addr,
	}
	tconn, err := tls.Dial("tcp", s.addr, TLSConfig)
	if err != nil {
		return nil, err
	}
	return tconn, err
}

// getConn returns a connection to the command and control server
func (s RemoteServer) getConn() (net.Conn, error) {
	return s.dial()
}
