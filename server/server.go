package server

import (
	"github.com/gary-kim/cmdctrl/shared"
)

type client struct {
	clientID string
	queue    []shared.PendingAction
}

// AddToQueue adds the given PendingAction onto the Queue
func (c client) AddToQueue(pa shared.PendingAction) error {
	queue = append(queue, pa)
	return nil
}
