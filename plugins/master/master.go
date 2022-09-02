package master

import (
	"fmt"
	"log"
	"main/core/types"
	"main/core/whats"
	"main/plugins"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var PlugMaster = plugins.Add("Master", plugins.MasterValidator)

func init() {

	PlugMaster.CommandAddMany([]*types.Command{
		{
			Cmd:         []string{"+user", "-user"},
			Description: "Add or Remove user",
			Usage:       "{cmd} @ mention user",
			Execute:     userManage,
		}, {
			Cmd:         []string{"+usern", "-usern"},
			Description: "Add or Remove user",
			Usage:       "{cmd} <number>",
			Execute:     userManageN,
		}, {
			Cmd:         []string{"+group", "-group"},
			Description: "Add or Remove this group",
			Usage:       "{cmd}",
			Execute:     userManage,
		}, {
			Cmd:         []string{"+master", "-master"},
			Description: "Add or Remove master",
			Usage:       "{cmd} @ mention user",
			Execute:     masterManage,
		}, {
			Cmd:         []string{"+block", "-block"},
			Description: "Block or Unblock user",
			Usage:       "{cmd} @ mention user",
			Execute:     blockManage,
		},
	})

}

func userManageN(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {

	var text = strings.Join(args, "")
	for _, old := range []string{"-", " ", "+"} {
		text = strings.ReplaceAll(text, old, "")
	}

	if len(text) > 10 && !strings.Contains(text, event.Info.Sender.Server) {
		var jid = fmt.Sprintf("%s@%s", text, event.Info.Sender.Server)
		log.Println(jid)

		if strings.HasPrefix(pattern, "+") {

			if plugins.RegisterID(plugins.AuthUser, jid) {
				whats.SendReactMessage(event, whats.ReactHandLike, client)
				log.Println(pattern, jid, true)

				return nil
			}
		} else if strings.HasPrefix(pattern, "-") {

			if plugins.RemoveID(plugins.AuthUser, jid) {
				whats.SendReactMessage(event, whats.ReactHandLike, client)
				log.Println(pattern, jid, true)

				return nil
			}
		}
	}

	whats.SendReactMessage(event, whats.ReactHandBad, client)
	return nil
}

func userManage(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {
	var jids = []string{}

	if ctx != nil {

		if ctx.MentionedJid != nil {
			jids = append(jids, ctx.MentionedJid...)
		} else {
			jids = append(jids, event.Info.Chat.ToNonAD().String())
		}
	} else {
		jids = append(jids, event.Info.Chat.ToNonAD().String())
	}

	if len(jids) > 0 {
		for _, jid := range jids {
			if strings.HasPrefix(pattern, "+") {

				if plugins.RegisterID(plugins.AuthUser, jid) {
					whats.SendReactMessage(event, whats.ReactHandLike, client)
					log.Println(pattern, jid, true)

					return nil
				}
			} else if strings.HasPrefix(pattern, "-") {

				if plugins.RemoveID(plugins.AuthUser, jid) {
					whats.SendReactMessage(event, whats.ReactHandLike, client)
					log.Println(pattern, jid, true)

					return nil
				}
			}
		}
	}

	whats.SendReactMessage(event, whats.ReactHandBad, client)
	return nil
}

func masterManage(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {
	var jids = []string{}

	if ctx != nil {

		if ctx.MentionedJid != nil {
			jids = append(jids, ctx.MentionedJid...)
		} else {
			jids = append(jids, event.Info.Chat.ToNonAD().String())
		}
	} else {
		jids = append(jids, event.Info.Chat.ToNonAD().String())
	}

	if len(jids) > 0 {
		for _, jid := range jids {
			if strings.HasPrefix(pattern, "+") {

				if plugins.RegisterID(plugins.AuthMaster, jid) {
					whats.SendReactMessage(event, whats.ReactHandLike, client)
					log.Println(pattern, jid, true)

					return nil
				}
			} else if strings.HasPrefix(pattern, "-") {

				if plugins.RemoveID(plugins.AuthMaster, jid) {
					whats.SendReactMessage(event, whats.ReactHandLike, client)
					log.Println(pattern, jid, true)

					return nil
				}
			}
		}
	}

	whats.SendReactMessage(event, whats.ReactHandBad, client)
	return nil
}

func blockManage(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {
	var jids = []string{}

	if ctx != nil {

		if ctx.MentionedJid != nil {
			jids = append(jids, ctx.MentionedJid...)
		} else {
			jids = append(jids, event.Info.Chat.ToNonAD().String())
		}
	} else {
		jids = append(jids, event.Info.Chat.ToNonAD().String())
	}

	if len(jids) > 0 {
		for _, jid := range jids {
			if strings.HasPrefix(pattern, "+") {

				if plugins.RegisterID(plugins.AuthBlocked, jid) {
					whats.SendReactMessage(event, whats.ReactHandLike, client)
					log.Println(pattern, jid, true)

					return nil
				}
			} else if strings.HasPrefix(pattern, "-") {

				if plugins.RemoveID(plugins.AuthBlocked, jid) {
					whats.SendReactMessage(event, whats.ReactHandLike, client)
					log.Println(pattern, jid, true)

					return nil
				}
			}
		}
	}

	whats.SendReactMessage(event, whats.ReactHandBad, client)
	return nil
}
