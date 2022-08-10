package types

type pluginList struct {
	list []*Plugin
}

type iPluginList interface {
	Add(name string, v Validator) *Plugin
	All() []*Plugin
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
