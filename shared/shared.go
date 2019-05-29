package shared

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gary-kim/cmdctrl/shared/ccmath"
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/pkg/errors"
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
func (p PendingAction) Run(addr string) error {
	if !p.Cmdctrlspec {
		cmd := exec.Command(p.Cmd, p.Args...)
		cmd.Run()
		return nil
	}
	switch p.Cmd {
	case "math":
		value, err := ccmath.Solve(strings.Join(p.Args, " "))
		if err != nil {
			return errors.Wrap(err, "Math solver")
		}
		res, err := http.PostForm(addr, url.Values{"q": {"Info"}, "info": {strings.Join(p.Args, " ") + " is " + strconv.FormatFloat(value, 'f', 5, 64)}})
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}
	return nil
}

// Compare compares the priorities of the PendingActions
func (p PendingAction) Compare(other queue.Item) int {
	return p.Priority - other.(PendingAction).Priority
}

// BadSplitter is a bad but working cli parser (sort of)
func BadSplitter(input string) []string {
	tr := []string{""}
	open := false
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case ' ':
			if !open {
				tr = append(tr, "")
			} else {
				tr[len(tr)-1] += " "
			}
			break
		case '"':
			open = !open
			break
		default:
			tr[len(tr)-1] += string(input[i])
		}
	}
	return tr
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
