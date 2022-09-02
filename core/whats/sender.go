package whats

import (
	"context"

	waProto "go.mau.fi/whatsmeow/binary/proto"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type HandlerSentMessage func(waTypes.JID, string, *waProto.Message)

var handlersSentMessage = []HandlerSentMessage{}

func AddSentHandler(handler HandlerSentMessage) {
	handlersSentMessage = append(handlersSentMessage, handler)
}

func SendTextMessage(jid waTypes.JID, text string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (whatsmeow.SendResponse, error) {
	if text_message, err := NewExtendedMessage(text, ctx); err != nil {
		return whatsmeow.SendResponse{}, err
	} else {
		if resp, err := client.SendMessage(context.Background(), jid, whatsmeow.GenerateMessageID(), text_message); err != nil {
			return resp, err
		} else {
			for _, handler := range handlersSentMessage {
				handler(jid, resp.ID, text_message)
			}

			return resp, err
		}
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

	if resp, err := client.SendMessage(context.Background(), event.Info.Chat, whatsmeow.GenerateMessageID(), this_message); err != nil {
		return resp, err
	} else {
		for _, handler := range handlersSentMessage {
			handler(event.Info.Chat, resp.ID, this_message)
		}

		return resp, err
	}
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

func SendMessage(jid waTypes.JID, message *waProto.Message, client *whatsmeow.Client) (whatsmeow.SendResponse, error) {
	if resp, err := client.SendMessage(context.Background(), jid, whatsmeow.GenerateMessageID(), message); err != nil {
		return resp, err
	} else {
		for _, handler := range handlersSentMessage {
			handler(jid, resp.ID, message)
		}

		return resp, err
	}
}
