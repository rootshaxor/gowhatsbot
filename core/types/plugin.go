package types

import (
	"main/core/whats"

	"go.mau.fi/whatsmeow"
)

type Validator func(interface{}, *whatsmeow.Client) (bool, error)

type Plugin struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Commands    []*Command `json:"commands"`
	Tags        []string   `json:"tags"`
	Permissions []string   `json:"permissions"`
	Disabled    bool       `json:"disabled"`
	Validate    Validator  `json:"validate"`
}

func (p *Plugin) CommandAdd(cmd *Command) int {
	p.Commands = append(p.Commands, cmd)

	return len(p.Commands)
}

func (p *Plugin) CommandAddMany(cmds []*Command) int {
	for _, c := range cmds {
		p.CommandAdd(c)
	}

	return len(p.Commands)
}

func (p *Plugin) CommandCopy(n string) *Command {
	for _, c := range p.Commands {
		for _, cc := range c.Cmd {
			if n == cc {
				return c
			}
		}
	}

	return &Command{}
}

func (p *Plugin) Call(i interface{}, a *whatsmeow.Client) []error {
	var resp []error

	if !p.Disabled {
		for _, x := range p.Commands {
			if !x.Disabled {

				var ev, ctx, cm, args, _ = whats.PrepareExec(i)
				if x.Exists(cm) || x.Passed {
					if err := x.Execute(cm, args, x, ev, ctx, a); err != nil {
						resp = append(resp, err)
					}
				}
			}
		}
	}

	return resp
}

func (p *Plugin) SetEnable() {
	p.Disabled = false
}

func (p *Plugin) SetDisabled() {
	p.Disabled = true
}
