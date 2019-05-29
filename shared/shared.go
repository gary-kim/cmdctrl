package shared

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/golang-collections/go-datastructures/queue"
)

// PendingAction represents a action requested of a client
type PendingAction struct {
	// This action requires special attention from cmdctrl
	Cmdctrlspec bool
	// This is the command to be run
	Cmd string
	// This is the arguments with which to run the cmd command with
	Args []string
	// This is the priority of the PendingAction
	Priority int
}

// Message represents a message between the server and client
//
// If something non-critical is added here, the minor version must change. If something critical is added or something is removed, the major version must be changed.
type Message struct {
	Version       string
	Success       bool
	Action        string
	PendingAction PendingAction
	SharedPass    string
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
	if !p.Cmdctrlspec {
		cmd := exec.Command(p.Cmd, p.Args...)
		cmd.Run()
		return nil
	}
	return nil
}

// Compare compares the priorities of the PendingActions
func (p PendingAction) Compare(other queue.Item) int {
	return p.Priority - other.(PendingAction).Priority
}

// Compatible returns whether the server and client versions are Compatible
//
// cmdctrl follows semver. Patch versions are ignored.
// Client can be minor ahead of server but not behind
// Major version difference will automatically fail.
func Compatible(clientVersion string, serverVersion string) bool {
	clientV := strings.Split(clientVersion[1:], ".")
	serverV := strings.Split(serverVersion[1:], ".")
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
