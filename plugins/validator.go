package plugins

import (
	"encoding/json"
	"log"
	gerrors "main/core/errors"
	"main/core/utils"
	"os"
	"path"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type AuthType string

const (
	AuthMaster  AuthType = "master"
	AuthUser    AuthType = "user"
	AuthBlocked AuthType = "blocked"
)

var Auths map[string][]string = map[string][]string{}
var AuthPath = path.Join(utils.GetMyPath(), "auth.json")

func loadAuth() error {
	if _, err := os.Stat(AuthPath); err != nil {
		saveAuth()
		log.Println(err)
	}

	if jbyte, err := os.ReadFile(AuthPath); err != nil {
		log.Println(err)
		return err
	} else {

		if len(jbyte) <= 2 {
			jbyte = []byte("{}")
		}

		err := json.Unmarshal(jbyte, &Auths)
		if err != nil {
			log.Println(err)
		}

		return err
	}
}

func saveAuth() error {

	if jbyte, err := json.MarshalIndent(Auths, "", "  "); err != nil {
		return err
	} else {
		return os.WriteFile(AuthPath, jbyte, os.ModeAppend)
	}
}

func RegisterID(atype AuthType, id string) bool {
	loadAuth()

	if ExistsID(atype, id) {
		return true
	} else {
		Auths[string(atype)] = append(Auths[string(atype)], id)
		saveAuth()

		return true
	}
}

func RemoveID(atype AuthType, id string) bool {
	loadAuth()

	if !ExistsID(atype, id) {
		return true

	} else {
		var temp []string = []string{}
		var removed bool

		for _, cid := range Auths[string(atype)] {
			if id != cid {
				temp = append(temp, cid)
			} else {
				removed = removed || true
			}
		}

		Auths[string(atype)] = append([]string{}, temp...)

		saveAuth()
		return removed
	}
}

func ExistsID(atype AuthType, id string) bool {
	loadAuth()

	if _, ok := Auths[string(atype)]; !ok {
		Auths[string(atype)] = []string{}
		saveAuth()

		return ok
	} else {

		var exists bool

		for _, cid := range Auths[string(atype)] {
			exists = exists || (cid == id)
		}

		return exists

	}
}

func UserValidator(i interface{}, c *whatsmeow.Client) (bool, error) {
	switch e := i.(type) {
	case *events.Message:
		{
			var sender = e.Info.Sender.ToNonAD().String()
			var chat = e.Info.Chat.ToNonAD().String()
			if e.Info.Chat.User != "status" && ((ExistsID(AuthUser, sender) || ExistsID(AuthMaster, sender) || ExistsID(AuthUser, chat)) && !ExistsID(AuthBlocked, sender)) || e.Info.IsFromMe {
				return false, nil
			}
		}
	}

	return false, gerrors.NewPlugin("PluginValidator", "Not allowed")
}

func MasterValidator(i interface{}, c *whatsmeow.Client) (bool, error) {
	switch e := i.(type) {
	case *events.Message:
		{
			var sender = e.Info.Sender.ToNonAD().String()
			if e.Info.Chat.User != "status" && (ExistsID(AuthMaster, sender) && !ExistsID(AuthBlocked, sender)) || e.Info.IsFromMe {
				return false, nil
			}
		}
	}

	return false, gerrors.NewPlugin("PluginValidator", "Not allowed")
}
