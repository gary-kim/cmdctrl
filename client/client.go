package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
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
	opt      Options
	conn     *net.Conn
}

// RemoteWSServer represents a single cmdctrl server that this client is configured to connect to in Websocket mode
type RemoteWSServer struct {
	addr     string
	clientID string
	opt      Options
	conn     *websocket.Conn
}

// RemoteServer interface represents a single cmdctrl server that this client is configured to connect to.
type RemoteServer interface {
	run()
}

// Options represents the options given by the user when cmdctrl was started in client mode
type Options struct {
	RESTMode           bool
	RESTUpdateInterval int
	SharedPass         string
}

var primaryServer RemoteRESTServer

// RunClient runs cmdctrl in client mode
func RunClient(addr string, opt Options) {
	if opt.RESTMode {
		primaryServer = RemoteRESTServer{
			addr: addr,
			opt:  opt,
		}
	}
	primaryServer.run()
}

// query will query the cmdctrl server with the given url.Values. It can also be given an int for an expected status code. If the expected status code is -1, it will ignore the status code.
func (s RemoteRESTServer) query(query url.Values, status int) (*shared.Message, error) {
	tr := shared.Message{}
	query.Add("client", s.clientID)
	res, err := http.PostForm(s.addr, query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if status != -1 && res.StatusCode != status {
		return nil, errors.Errorf("Unexpected http status code: %v", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tr)
	if err != nil {
		return nil, err
	}
	if !shared.Compatible(cmd.Version, tr.Version) {
		return nil, errors.New("Server and client versions are incompatible")
	}
	if s.opt.SharedPass != tr.SharedPass {
		return nil, errors.New("Cannot verify server identity")
	}
	return &tr, nil
}

// queryCommand querys the server for new PendingAction(s)
func (s RemoteRESTServer) queryCommand() (*shared.PendingAction, error) {
	pa := &shared.PendingAction{}
	queryReturn, err := s.query(url.Values{"q": {"RequestedAction"}}, 200)
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
func (s *RemoteRESTServer) registerClient() error {
	possibleLetters := []byte("123567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < 20; i++ {
		s.clientID += string(possibleLetters[rand.Intn(len(possibleLetters))])
	}

	// register with server
	res, err := s.query(url.Values{"q": {"RegisterClient"}}, 200)
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
		return
	}
	fmt.Printf("Client successfully registered with ID: %s\n", s.clientID)
	time.Sleep(1 * time.Second)
	for {
		pa, err := s.queryCommand()
		if err != nil {
			fmt.Printf("Could not query for command from server: %s", err)
		}
		pa.Run(s.addr)

		duration, err := time.ParseDuration(strconv.Itoa(s.opt.RESTUpdateInterval) + "s")
		if err != nil {
			fmt.Printf("Could not parse time duration %s", string(s.opt.RESTUpdateInterval)+"s")
		}
		time.Sleep(duration)
	}
}

func (s RemoteWSServer) run() {
	return
}

func (s RemoteRESTServer) verifySharedPass(pass string) bool {
	return s.opt.SharedPass == pass
}
