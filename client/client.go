package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gary-kim/cmdctrl/shared"
	"github.com/gorilla/websocket"
)

// RemoteRESTServer represents a single cmdctrl server that this client is configured to connect to in REST Mode
type RemoteRESTServer struct {
	addr     string
	clientID string
	conn     *net.Conn
}

// RemoteWSServer represents a single cmdctrl server that this client is configured to connect to in Websocket mode
type RemoteWSServer struct {
	addr     string
	clientID string
	conn     *websocket.Conn
}

// RemoteServer interface represents a single cmdctrl server that this client is configured to connect to.
type RemoteServer interface {
	run()
}

// Options represents the options given by the user when cmdctrl was started in client mode
type Options struct {
	RESTMode bool
}

var primaryServer RemoteServer

// RunClient runs cmdctrl in client mode
func RunClient(addr string, opt Options) {
	if opt.RESTMode {
		primaryServer = RemoteRESTServer{
			addr: addr,
		}
	}
	primaryServer.run()
}

func (s RemoteRESTServer) queryCommand() (shared.PendingAction, error) {
	errReturn := shared.PendingAction{}
	res, err := http.PostForm(s.addr, url.Values{"client": {s.clientID}, "q": {"RequestedAction"}})
	if err != nil {
		return errReturn, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errReturn, errors.New("Query for command did not return status code 200")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errReturn, err
	}
	pa := shared.PendingAction{}
	err = pa.FromJSON(body)
	if err != nil {
		return errReturn, err
	}
	return pa, nil
}

func (s RemoteRESTServer) run() {
	for {
		pa, err := s.queryCommand()
		if err != nil {
			fmt.Printf("Could not query for command from server: %s", err)
		}
		pa.Run()
		time.Sleep(60 * time.Second)
	}
}

func (s RemoteWSServer) run() {
	return
}
