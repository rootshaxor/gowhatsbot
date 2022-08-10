package whats

import (
	"context"

	waProto "go.mau.fi/whatsmeow/binary/proto"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func SendTextMessage(event *events.Message, text string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (whatsmeow.SendResponse, error) {
	if text_message, err := NewExtendedMessage(text, ctx); err != nil {
		return whatsmeow.SendResponse{}, err
	} else {
		return client.SendMessage(context.Background(), event.Info.Chat, whatsmeow.GenerateMessageID(), text_message)
	}
}

func SendReactMessage(event *events.Message, react Reactions, client *whatsmeow.Client) (whatsmeow.SendResponse, error) {
	this_message := &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJid:   proto.String(event.Info.Chat.ToNonAD().String()),
				Participant: proto.String(event.Info.Sender.ToNonAD().String()),
				FromMe:      proto.Bool(event.Info.IsFromMe),
				Id:          &event.Info.ID,
			},
			Text: proto.String(string(react)),
		},
	}

	return client.SendMessage(context.Background(), event.Info.Chat, whatsmeow.GenerateMessageID(), this_message)
}

func SendChatPresence(event *events.Message, media waTypes.ChatPresenceMedia, state waTypes.ChatPresence, client *whatsmeow.Client) error {

	return client.SendChatPresence(event.Info.Chat, state, media)
}

func SendTyping(event *events.Message, client *whatsmeow.Client) error {
	return SendChatPresence(event, waTypes.ChatPresenceMediaText, waTypes.ChatPresenceComposing, client)
}

func SendStopTyping(event *events.Message, client *whatsmeow.Client) error {
	return SendChatPresence(event, waTypes.ChatPresenceMediaText, waTypes.ChatPresencePaused, client)
}

func SendRecording(event *events.Message, client *whatsmeow.Client) error {
	return SendChatPresence(event, waTypes.ChatPresenceMediaAudio, waTypes.ChatPresenceComposing, client)
}

func SendStopRecording(event *events.Message, client *whatsmeow.Client) error {
	return SendChatPresence(event, waTypes.ChatPresenceMediaAudio, waTypes.ChatPresencePaused, client)
}

func SendMessage(chat_jid waTypes.JID, message *waProto.Message, client *whatsmeow.Client) (whatsmeow.SendResponse, error) {
	return client.SendMessage(context.Background(), chat_jid, whatsmeow.GenerateMessageID(), message)
}
