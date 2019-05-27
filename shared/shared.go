package shared

import (
	"encoding/json"
	"os/exec"
)

// PendingAction represents a action requested of a client
type PendingAction struct {
	// This action requires special attention from cmdctrl
	cmdctrlspec bool
	// This is the command to be run
	cmd string
	// This is the arguments with which to run the cmd command with
	args []string
	// This is the action that must be completed by cmdctrl if cmdctrlspec is true
	action string
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
