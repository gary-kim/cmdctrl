package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gary-kim/cmdctrl/cmd"
	"github.com/gary-kim/cmdctrl/shared"
	"github.com/golang-collections/go-datastructures/queue"
)

type clients struct {
	Clients    []*client
	SharedPass string
}

type client struct {
	clientID string
	queue    *queue.PriorityQueue
}

// CmdCtrlServer is the cmdctrl server itself
var CmdCtrlServer clients

// Options represents options to be set for running the cmdctrl server
type Options struct {
	SharedPass string
}

func (c client) Info() string {
	return c.clientID
}

// RunServer begins the cmdctrl server
func RunServer(addr string, opt Options) {
	CmdCtrlServer := clients{
		SharedPass: opt.SharedPass,
	}
	CmdCtrlServer.Run(addr, opt)
}

// Run starts the cmdctrl server
func (c *clients) Run(addr string, opt Options) {
	// start the server
	http.HandleFunc("/post", c.handle)
	http.HandleFunc("/webui", c.webui)
	http.HandleFunc("/clients", c.clientList)
	go func() {
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
	fmt.Println("Navigate to http://localhost/webui to control. There is a command line but the webui is far easier to use.")
	stdin := bufio.NewReader(os.Stdin)
	for {
		currentLineBytes, err := stdin.ReadBytes('\n')
		if err != nil {
			os.Exit(1)
		}
		currentLine := strings.ReplaceAll(string(currentLineBytes), "\n", "")
		switch strings.Split(currentLine, " ")[0] {
		case "exit":
			fmt.Println("cmdctrl is now exiting")
			os.Exit(0)
			break
		case "exec":
			if len(strings.Split(currentLine, " ")) < 3 {
				fmt.Println("client id must be the second argument and a priority must be the 2nd argument")
				continue
			}
			currentClient, err := c.getClient(strings.Split(currentLine, " ")[1])
			if err != nil {
				fmt.Printf("Could not find client %s registed with the server\n", strings.Split(currentLine, " ")[1])
			}
			priority, err := strconv.Atoi(strings.Split(currentLine, " ")[2])
			if err != nil {
				fmt.Printf("Could not parse time %s\n", strings.Split(currentLine, " ")[2])
			}
			currentClient.queue.Put(&shared.PendingAction{
				Cmdctrlspec: false,
				Cmd:         strings.Split(currentLine, " ")[3],
				Args:        shared.BadSplitter(strings.Join(strings.Split(currentLine, " ")[4:], " ")),
				Priority:    priority,
			})
			break
		case "clients":
			for _, curr := range c.Clients {
				fmt.Println(curr.Info())
			}
		}
	}
}

func (c *clients) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Incorrect request", http.StatusForbidden)
		fmt.Println("Incorrect Request")
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
		fmt.Println("Malformed Request")
	}
	forClient := shared.Message{
		Version:    cmd.Version,
		SharedPass: c.SharedPass,
	}
	clientID := r.FormValue("client")

	switch r.FormValue("q") {
	case "RegisterClient":
		c.Clients = append(c.Clients, &client{clientID: clientID, queue: queue.NewPriorityQueue(1)})
		forClient.Action = "ClientRegistered"
		forClient.Success = true
		sendToClient(forClient, w)
		break
	case "RequestedAction":
		currentClient, err := c.getClient(clientID)
		if err != nil {
			http.Error(w, "Cannot find client", http.StatusInternalServerError)
			return
		}
		if currentClient.queue.Empty() {
			forClient.Action = "NoAction"
			forClient.Success = true
			sendToClient(forClient, w)
			return
		}
		pendingAction, err := currentClient.queue.Get(1)
		if err != nil {
			http.Error(w, "Failed to find job for client", http.StatusInternalServerError)
		}
		forClient.Action = "PendingAction"
		forClient.PendingAction = *(pendingAction[0].(*shared.PendingAction))
		forClient.Success = true
		sendToClient(forClient, w)
	case "Info":
		fmt.Printf("Info from a client: %s", r.FormValue("info"))
	}
}

func sendToClient(forClient shared.Message, w http.ResponseWriter) {
	tr, err := json.Marshal(forClient)
	if err != nil {
		http.Error(w, "Failed to format message for client", http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(tr)
}

func (c *clients) getClient(clientID string) (*client, error) {
	for _, curr := range c.Clients {
		if strings.Index(curr.clientID, clientID) == 0 {
			return curr, nil
		}
	}
	return nil, errors.New("Cannot find client")
}

func (c *clients) webui(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(ControlUITemplate))
		return
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			return
		}
		currentClient, err := c.getClient(r.FormValue("client"))
		if err != nil {
			return
		}
		priority, _ := strconv.Atoi(r.FormValue("priority"))
		cmd := "math"
		args := []string{r.FormValue("input")}
		if r.FormValue("action") == "cmd" {
			cmd = strings.Fields(r.FormValue("input"))[0]
			args = shared.BadSplitter(strings.Join(strings.Fields(r.FormValue("input"))[1:], " "))
		}
		currentClient.queue.Put(&shared.PendingAction{
			Priority:    priority,
			Cmdctrlspec: (r.FormValue("action") != "cmd"),
			Cmd:         cmd,
			Args:        args,
		})
		fmt.Println("added")
		w.Write([]byte(ControlUITemplate))
	}
}

func (c *clients) clientList(w http.ResponseWriter, r *http.Request) {
	clientList := "Client List:\n"
	for _, curr := range c.Clients {
		clientList = clientList + curr.Info() + "\n"
	}
	w.Write([]byte(clientList))
}
