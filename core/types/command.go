package types

import (
	"encoding/json"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

type CommandExecute func(string, []string, Command, *events.Message, *waProto.ContextInfo, *whatsmeow.Client) error
type Validator func(interface{}, *whatsmeow.Client) (bool, error)

type Command struct {
	Cmd         []string       `json:"cmd"`
	Description string         `json:"description"`
	Usage       string         `json:"usage"`
	Passed      bool           `json:"passed"`
	Disabled    bool           `json:"disabled"`
	Execute     CommandExecute `json:"execute"`
}

func (c *Command) SetDisabled() {
	c.Disabled = true
}

func (c *Command) GetUsage() string {
	c.Disabled = true

	return ""
}

func (c *Command) SetEnable() {
	c.Disabled = false
}

func (c *Command) Exists(cm string) bool {
	for _, ca := range c.Cmd {
		if ca == cm {
			return true
		}
	}
	return false
}

func (c Command) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Cmd         []string `json:"cmd"`
		Description string   `json:"description"`
		Usage       string   `json:"usage"`
		Passed      bool     `json:"passed"`
		Disabled    bool     `json:"disabled"`
	}{
		Cmd:         c.Cmd,
		Description: c.Description,
		Usage:       c.Usage,
		Passed:      c.Passed,
		Disabled:    c.Disabled,
	})
}
