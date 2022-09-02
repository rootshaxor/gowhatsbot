package types

type pluginList struct {
	list []*Plugin
}

type iPluginList interface {
	Add(name string, v Validator) *Plugin
	All() []*Plugin
	GetOne(name string) *Plugin
	SetDisable(name string) bool
	SetEnable(name string) bool
	SwitchDisable(name string, state bool) bool
	SwitchDisableCommand(name string, state bool) bool
}

func NewPluginList() iPluginList {
	var list pluginList = pluginList{
		list: []*Plugin{},
	}
	return &list
}

func (el *pluginList) Add(name string, v Validator) *Plugin {
	var newExt = Plugin{
		Name:     name,
		Validate: v,
	}

	el.list = append(el.list, &newExt)

	return &newExt
}

func (el *pluginList) All() []*Plugin {

	return el.list
}

func (el *pluginList) GetOne(name string) *Plugin {
	for _, plug := range el.All() {
		if plug.Name == name {
			return plug
		}
	}

	return nil
}

func (el *pluginList) SetDisable(name string) bool {
	return el.SwitchDisable(name, true)
}

func (el *pluginList) SetEnable(name string) bool {
	return el.SwitchDisable(name, false)
}

func (el *pluginList) SwitchDisable(name string, state bool) bool {
	if plg := el.GetOne(name); plg != nil {
		plg.Disabled = state

		return true
	} else {
		return false
	}
}

func (el *pluginList) SwitchDisableCommand(name string, state bool) bool {
	var success bool
	for _, plgs := range el.All() {
		for _, cmd := range plgs.Commands {
			for _, c := range cmd.Cmd {
				if c == name {
					cmd.Disabled = state

					success = success && (cmd.Disabled == state)
				}
			}
		}
	}

	return success
}
