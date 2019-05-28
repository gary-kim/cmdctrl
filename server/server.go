package server

import (
	"encoding/json"
	"errors"
	"net/http"

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

// RunServer starts the cmdctrl server
func Run(addr string, opt Options) {
	// start the server

}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Incorrect request", http.StatusForbidden)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
	}
	forClient := shared.Message{
		Version: cmd.Version,
	}
	clientID := r.FormValue("client")
	switch r.FormValue("q") {
	case "RegisterClient":
		clients = append(clients, &client{clientID: clientID})
		forClient.Action = "ClientRegistered"
		tr, err := json.Marshal(forClient)
		if err != nil {
			http.Error(w, "Failed to format message for client", http.StatusInternalServerError)
		}
		w.Write(tr)
		break
	case "RequestedAction":
		currentClient, err := getClient(clientID)
		if err != nil {
			http.Error(w, "Cannot find client", http.StatusInternalServerError)
		}
		if currentClient.queue.Empty() {
			return
		}
		pendingAction, err := currentClient.queue.Get(1)
		if err != nil {
			http.Error(w, "Failed to find job for client", http.StatusInternalServerError)
		}
		forClient.PendingAction = pendingAction[0].(shared.PendingAction)

	}
}

func getClient(clientID string) (*client, error) {
	for _, curr := range clients {
		if curr.clientID == clientID {
			return curr, nil
		}
	}
	return nil, errors.New("Cannot find client")
}
