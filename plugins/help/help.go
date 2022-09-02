package help

import (
	"fmt"
	gerrors "main/core/errors"
	"main/core/helper"
	"main/core/texts"
	"main/core/types"
	"main/core/whats"
	"main/plugins"
	"math/rand"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var PlugHelp = plugins.Add("Helper", plugins.UserValidator)

func init() {

	// PlugHelp.SetDisabled()

	PlugHelp.CommandAddMany([]*types.Command{
		{
			Cmd:         []string{".h", ".help", ".menu"},
			Description: "Show all command menu",
			Usage:       "{cmd}",
			Execute:     helpCommand,
		},
	})

}

func helpCommand(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {

	var res []string

	for _, ex := range plugins.All() {
		var _, err = ex.Validate(event, client)

		if len(ex.Commands) > 0 && !ex.Disabled && err == nil {

			res = append(res, texts.QuoteBy("🚀 "+ex.Name+" :", texts.QuoteBold))
			for _, cmd := range ex.Commands {
				if !cmd.Disabled && len(cmd.Cmd) > 0 {

					res = append(res, texts.AddTab(1, texts.QuoteBy(cmd.Description, texts.QuoteItalic)))
					// res = append(res, AddTab(1, strings.Join(cmd.Cmd, ", ")))

					if rep, err := helper.MapMe(cmd); err != nil {
						return gerrors.NewPlugin(PlugHelp.Name, err.Error())
					} else {

						var the_cmd = append([]string{}, cmd.Cmd...)
						// var other_cmd = []string{}

						if len(cmd.Cmd) > 1 {
							rand.Seed(time.Now().Unix())
							var rand_index = rand.Intn(len(cmd.Cmd))
							the_cmd = []string{cmd.Cmd[rand_index]}
							// for i, cm := range cmd.Cmd {
							// 	if i != rand_index {
							// 		other_cmd = append(other_cmd, texts.QuoteBy(cm, texts.QuoteBold))
							// 	}
							// }
						}

						for _, c := range the_cmd {
							var usage = cmd.Usage
							rep["cmd"] = texts.QuoteBy(c, texts.QuoteBold, texts.QuoteItalic)
							for key, val := range rep {
								switch rval := val.(type) {
								case bool:
									usage = strings.ReplaceAll(usage, "{"+key+"}", fmt.Sprintf("%v", rval))
								case string:
									usage = strings.ReplaceAll(usage, "{"+key+"}", rval)

								}
							}

							// if len(other_cmd) > 0 {
							// 	res = append(res, texts.AddTab(1, texts.QuoteBy("Alt : "+strings.Join(other_cmd, " , ")+"", texts.QuoteItalic)))
							// }

							res = append(res, texts.AddTab(1, usage))
						}
						res = append(res, "")
					}
				}
			}
		}
	}

	sentctx, _ := whats.SanitizeContext(event, event.Info.IsGroup, client)
	whats.SendTextMessage(event.Info.Chat, strings.Join(res, "\n"), sentctx, client)

	return nil
}
