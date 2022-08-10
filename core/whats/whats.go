package whats

import (
	gerrors "main/core/errors"
	"reflect"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"

	// waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var expirationChat map[string]uint32 = map[string]uint32{}

func EventIs(e interface{}, i interface{}) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(i)
}

func GetTextContext(e *events.Message) (string, *waProto.ContextInfo) {
	var text string = ""
	var ctx *waProto.ContextInfo

	if msg := e.Message; msg != nil {

		if msg_type := msg.GetAudioMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()

		} else if msg_type := msg.GetButtonsMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetContentText()

		} else if msg_type := msg.GetButtonsResponseMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetSelectedButtonId()

		} else if msg_type := msg.GetContactMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetVcard()

		} else if msg_type := msg.GetConversation(); msg_type != "" {
			text = msg_type

		} else if msg_type := msg.GetDocumentMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetFileName()

		} else if msg_type := msg.GetExtendedTextMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetText()

		} else if msg_type := msg.GetImageMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetCaption()

		} else if msg_type := msg.GetListMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetDescription()

		} else if msg_type := msg.GetListResponseMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.SingleSelectReply.GetSelectedRowId()

		} else if msg_type := msg.GetProductMessage(); msg_type != nil {
			ctx = msg_type.ContextInfo
			text = msg_type.GetBody()

		} else if msg_type := msg.GetReactionMessage(); msg_type != nil {
			text = msg_type.GetText()

		} else if msg_type := msg.GetStickerMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()

		} else if msg_type := msg.GetTemplateButtonReplyMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetSelectedId()

		} else if msg_type := msg.GetTemplateMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()

			if msg_subtype := msg_type.GetFourRowTemplate(); msg_subtype != nil {
				text = msg_subtype.Content.GetNamespace()

			} else if msg_subtype := msg_type.GetHydratedTemplate(); msg_subtype != nil {
				text = msg_subtype.GetTemplateId()

			} else if msg_subtype := msg_type.GetHydratedFourRowTemplate(); msg_subtype != nil {
				text = msg_subtype.GetTemplateId()

			}

		} else if msg_type := msg.GetVideoMessage(); msg_type != nil {
			ctx = msg_type.GetContextInfo()
			text = msg_type.GetCaption()

		}

	}

	return text, ctx
}

func SanitizeContext(e *events.Message, quote bool, client *whatsmeow.Client) (*waProto.ContextInfo, error) {
	_, ctx := GetTextContext(e)

	var newctx waProto.ContextInfo
	if ctx != nil {
		newctx = waProto.ContextInfo{
			Expiration: ctx.Expiration,
		}
	}

	if exp, ok := expirationChat[e.Info.Chat.User]; ok {
		newctx.Expiration = proto.Uint32(exp)

	} else {
		if newctx.Expiration != nil {
			expirationChat[e.Info.Chat.User] = *newctx.Expiration
		}
	}

	if quote {
		newctx.StanzaId = proto.String(e.Info.ID)
		newctx.Participant = proto.String(e.Info.Sender.ToNonAD().String())
		newctx.QuotedMessage = e.Message
	}

	return &newctx, gerrors.NewCore("SanitizeContext", "Can't Sanitize ContextInfo")
}

func PrepareExec(e interface{}) (*events.Message, *waProto.ContextInfo, string, []string, error) {
	var event = e.(*events.Message)

	var text, ctx = GetTextContext(event)
	text = strings.TrimSpace(text)

	if len(text) > 0 {
		var args = strings.Split(text, " ")
		return event, ctx, args[0], args[1:], nil
	} else {
		return event, ctx, "", []string{}, nil
	}

}
