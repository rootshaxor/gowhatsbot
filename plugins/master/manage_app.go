package master

import (
	"encoding/json"
	"log"
	"main/core/types"
	"main/core/utils"
	"main/core/whats"
	"main/plugins"
	"os"
	"path"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var configFile = path.Join(utils.GetMyPath(), "app.json")

func loadConfig() error {
	if _, err := os.Stat(configFile); err != nil {
		saveConfig()
		log.Println(err)
	}

	if jbyte, err := os.ReadFile(configFile); err != nil {
		log.Println(err)
		return err
	} else {

		if len(jbyte) <= 2 {
			jbyte = []byte("{}")
		}

		var appConfig = map[string]bool{}
		if err := json.Unmarshal(jbyte, &appConfig); err != nil {
			log.Println(err)
		} else {
			for appname, state := range appConfig {
				plugins.SwitchDisable(appname, state)
			}
		}

		return err
	}
}

func saveConfig() error {
	var appConfig = map[string]bool{}
	for _, cmd := range plugins.All() {
		appConfig[cmd.Name] = cmd.Disabled
	}

	if jbyte, err := json.MarshalIndent(appConfig, "", "  "); err != nil {
		return err
	} else {
		return os.WriteFile(configFile, jbyte, os.ModeAppend)
	}
}

var PlugMasterInit = plugins.Add("InitApp", initListener)

func initListener(i interface{}, client *whatsmeow.Client) (bool, error) {

	switch i.(type) {
	case *events.Connected:
		loadConfig()
		plugins.SwitchDisable("InitApp", true)
	}

	return false, nil
}

func init() {

	PlugMaster.CommandAddMany([]*types.Command{
		{
			Cmd:         []string{".app"},
			Description: "Set enable / disable App.",
			Usage:       "{cmd} <on|off> <appname> ... ",
			Execute:     switchApp,
		}, {
			Cmd:         []string{".cmd"},
			Description: "Set enable / disable Command.",
			Usage:       "{cmd} <on|off> <cmdname> ... ",
			Execute:     switchCmd,
		},
	})

}

func switchApp(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {
	whats.SendTyping(event, client)

	if len(args) > 1 {
		var state bool
		switch args[0] {
		case "on":
			state = false

		case "off":
			state = true
		}

		for _, name := range args[1:] {
			plugins.SwitchDisable(name, state)
			saveConfig()
		}

		whats.SendReactMessage(event, whats.ReactHandLike, client)
		whats.SendStopTyping(event, client)
	} else {
		whats.SendReactMessage(event, whats.ReactHandBad, client)
		whats.SendStopTyping(event, client)
	}

	return nil
}

func switchCmd(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {
	whats.SendTyping(event, client)

	if len(args) > 1 {
		var state bool
		switch args[0] {
		case "on":
			state = false

		case "off":
			state = true
		}

		for _, name := range args[1:] {
			plugins.SwitchDisableCommand(name, state)
			saveConfig()
		}

		whats.SendReactMessage(event, whats.ReactHandLike, client)
		whats.SendStopTyping(event, client)
	} else {
		whats.SendReactMessage(event, whats.ReactHandBad, client)
		whats.SendStopTyping(event, client)
	}

	return nil
}
