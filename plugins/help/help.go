package help

import (
	"fmt"
	gerrors "main/core/errors"
	"main/core/helper"
	"main/core/texts"
	"main/core/types"
	"main/core/whats"
	"main/plugins"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var PlugHelp = plugins.Add("Helper", helpValidator)

func init() {

	PlugHelp.CommandAddMany([]types.Command{
		{
			Cmd:         []string{".h", ".help", ".menu", ".cmd", ".command"},
			Description: "Show all command menu",
			Usage:       "{cmd}",
			Execute:     helpCommand,
		},
	})

}

func helpValidator(i interface{}, c *whatsmeow.Client) (bool, error) {
	switch e := i.(type) {
	case *events.Message:
		if ct, err := c.Store.Contacts.GetContact(e.Info.Sender); err != nil {
			return false, gerrors.NewPlugin("Helper", err.Error())
		} else {
			if ct.Found || e.Info.IsFromMe {
				return false, nil
			} else {
				return false, gerrors.NewPlugin("Helper", "jid not in contact")
			}
		}
	}

	return false, gerrors.NewPlugin("Helper", "Not a Message")
}

func helpCommand(pattern string, args []string, cmd types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {

	var res []string

	for _, ex := range plugins.List.All() {
		if len(ex.Commands) > 0 && !ex.Disabled {

			res = append(res, texts.QuoteBy("# "+ex.Name+" :", types.QuoteBold))
			for _, cmd := range ex.Commands {
				if !cmd.Disabled && len(cmd.Cmd) > 0 {

					res = append(res, texts.AddTab(1, texts.QuoteBy(cmd.Description, types.QuoteItalic)))
					// res = append(res, AddTab(1, strings.Join(cmd.Cmd, ", ")))

					if rep, err := helper.MapMe(cmd); err != nil {
						return gerrors.NewPlugin(PlugHelp.Name, err.Error())
					} else {

						for _, c := range cmd.Cmd {
							var usage = cmd.Usage
							rep["cmd"] = texts.QuoteBy(c, types.QuoteBold, types.QuoteItalic)
							for key, val := range rep {
								switch rval := val.(type) {
								case bool:
									usage = strings.ReplaceAll(usage, "{"+key+"}", fmt.Sprintf("%b", rval))
								case string:
									usage = strings.ReplaceAll(usage, "{"+key+"}", rval)

								}

							}
							res = append(res, texts.AddTab(1, usage))
						}
						res = append(res, "")
					}
				}
			}
		}
	}

	sentctx, _ := whats.SanitizeContext(event, event.Info.IsGroup, client)
	whats.SendTextMessage(event, strings.Join(res, "\n"), sentctx, client)

	return nil
}
