package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gary-kim/cmdctrl/cmd"
	"github.com/gary-kim/cmdctrl/shared"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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

// query will query the cmdctrl server with the given url.Values. It can also be given an int for an expected status code. If the expected status code is -1, it will ignore the status code.
func (s RemoteRESTServer) query(query url.Values, status int) (*shared.Message, error) {
	tr := shared.Message{}
	res, err := http.PostForm(s.addr, query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if status != -1 && res.StatusCode != status {
		return nil, errors.New("Unexpected http status code")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tr)
	if err != nil {
		return nil, err
	}
	if shared.Compatible(cmd.Version, tr.Version) {
		return nil, errors.New("Server and client versions are incompatible")
	}
	return &tr, nil
}

func (s RemoteRESTServer) queryCommand() (*shared.PendingAction, error) {
	pa := &shared.PendingAction{}
	queryReturn, err := s.query(url.Values{"client": {s.clientID}, "q": {"RequestedAction"}}, 200)
	if err != nil {
		return pa, err
	}
	if queryReturn.Success != true {
		return pa, errors.New("Server unsuccessful")
	}
	if queryReturn.Action == "NoAction" {
		return pa, nil
	}
	if queryReturn.Action == "PendingAction" {
		return &queryReturn.PendingAction, nil
	}
	return pa, nil
}

// registerclient will create a client id and register the client id with the server
func (s RemoteRESTServer) registerClient() error {
	possibleLetters := []byte("123567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randomB := rand.Perm(len(possibleLetters))
	for _, curr := range randomB {
		s.clientID += string(possibleLetters[curr])
	}

	// register with server
	res, err := s.query(url.Values{"client": {s.clientID}, "q": {"RegisterClient"}}, 200)
	if err != nil {
		return err
	}
	if !res.Success || res.Action != "ClientRegistered" {
		return errors.New("Could not register client")
	}
	return nil
}

func (s RemoteRESTServer) run() {
	err := s.registerClient()
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to register client"))
	}
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
