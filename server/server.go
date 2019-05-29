package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"log"

	"github.com/gary-kim/cmdctrl/cmd"
	"github.com/gary-kim/cmdctrl/shared"
	"github.com/golang-collections/go-datastructures/queue"
)

var clients []*client

type client struct {
	clientID string
	queue    *queue.PriorityQueue
}

// Options represents options to be set for running the cmdctrl server
type Options struct {
}

// AddToQueue adds the given PendingAction onto the Queue
func (c client) AddToQueue(pa shared.PendingAction) error {
	c.queue.Put(&pa)
	return nil
}

// Run starts the cmdctrl server
func Run(addr string, opt Options) {
	// start the server
	http.HandleFunc("/post", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request")
	if r.Method != "POST" {
		http.Error(w, "Incorrect request", http.StatusForbidden)
		fmt.Println("Incorrect Request")
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
		fmt.Println("Malformed Request")
	}
	forClient := shared.Message{
		Version: cmd.Version,
	}
	fmt.Println("Checks passed")
	clientID := r.FormValue("client")
	fmt.Println(clientID)

	switch r.FormValue("q") {
	case "RegisterClient":
		clients = append(clients, &client{clientID: clientID, queue: queue.NewPriorityQueue(1)})
		forClient.Action = "ClientRegistered"
		forClient.Success = true
		sendToClient(forClient, w)
		break
	case "RequestedAction":
		currentClient, err := getClient(clientID)
		if err != nil {
			http.Error(w, "Cannot find client", http.StatusInternalServerError)
			return
		}		
		currentClient.queue.Put(shared.PendingAction{
			Cmd: "notify-send",
			Priority: 1,
			Args: []string{"cmdctrl", "Complete"},
			Cmdctrlspec: false,
		})
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
		forClient.PendingAction = pendingAction[0].(shared.PendingAction)
		forClient.Success = true
		sendToClient(forClient, w)

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

func getClient(clientID string) (*client, error) {
	for _, curr := range clients {
		if curr.clientID == clientID {
			return curr, nil
		}
	}
	return nil, errors.New("Cannot find client")
}
