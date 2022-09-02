package plugins

import (
	"log"
	"main/core/types"

	"go.mau.fi/whatsmeow"
)

var List = types.NewPluginList()

var Add = List.Add
var All = List.All
var GetOne = List.GetOne
var SetDisable = List.SetDisable
var SetEnable = List.SetEnable
var SwitchDisable = List.SwitchDisable
var SwitchDisableCommand = List.SwitchDisableCommand

func PluginExecutor(e interface{}, client *whatsmeow.Client) {

	for _, a := range All() {
		if a.Validate != nil {
			if show, err := a.Validate(e, client); err != nil {
				if show {
					log.Println(err)
				}
			} else {
				if errs := a.Call(e, client); len(errs) > 0 {
					log.Println(errs)
				}
			}
		}
	}
}

func init() {
	log.Println("Plugins package loaded")

}
