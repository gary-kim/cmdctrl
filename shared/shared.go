package shared

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

// PendingAction represents a action requested of a client
type PendingAction struct {
	// This action requires special attention from cmdctrl
	cmdctrlspec bool
	// This is the command to be run
	cmd string
	// This is the arguments with which to run the cmd command with
	args []string
	// This is the priority of the PendingAction
	priority int
}

// Message represents a message between the server and client
//
// If something non-critical is added here, the minor version must change. If something critical is added or something is removed, the major version must be changed.
type Message struct {
	Version       string
	Success       bool
	Action        string
	PendingAction PendingAction
}

// ToJSON provides a JSON output of the contents of PendingAction
func (p PendingAction) ToJSON() (string, error) {
	output, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(output), err
}

// FromJSON takes the input JSON and defines p from it
func (p PendingAction) FromJSON(input []byte) error {
	err := json.Unmarshal(input, &p)
	return err
}

// Run runs the PendingAction
func (p PendingAction) Run() error {
	if !p.cmdctrlspec {
		cmd := exec.Command(p.cmd, p.args...)
		cmd.Run()
		return nil
	}
	return nil
}

// Compare compares the priorities of the PendingActions
func (p PendingAction) Compare(other PendingAction) int {
	return p.priority - other.priority
}

// Compatible returns whether the server and client versions are Compatible
func Compatible(clientVersion string, serverVersion string) bool {
	clientV := strings.Split(clientVersion, ".")
	serverV := strings.Split(serverVersion, ".")
	if clientV[0] != serverV[0] {
		return false
	}
	clientMinor, err := strconv.ParseInt(clientV[1], 10, 64)
	if err != nil {
		return false
	}
	serverMinor, err := strconv.ParseInt(serverV[1], 10, 64)
	if err != nil {
		return false
	}
	if clientMinor < serverMinor {
		return false
	}
	return true
}
