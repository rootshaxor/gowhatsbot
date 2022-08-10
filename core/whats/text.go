package whats

import (
	_ "image/jpeg"
	_ "image/png"

	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func NewConversation(text string) (*waProto.Message, error) {

	return &waProto.Message{Conversation: &text}, nil
}

func NewExtendedMessage(text string, ctx *waProto.ContextInfo) (*waProto.Message, error) {

	var extendedtext_message = &waProto.ExtendedTextMessage{
		Text:        &text,
		ContextInfo: ctx,
	}

	return &waProto.Message{ExtendedTextMessage: extendedtext_message}, nil
}
