package plugins

import (
	"log"
	"main/core/types"

	"go.mau.fi/whatsmeow"
)

var List = types.NewPluginList()
var WhatsClient *whatsmeow.Client

func Add(name string, vd types.Validator) *types.Plugin {
	return List.Add(name, vd)
}

func init() {
	log.Println("Plugins package loaded")

}
