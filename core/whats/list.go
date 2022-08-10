package whats

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func NewListMessage(title, description, toggle, footer string, section []*waProto.ListMessage_Section, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {
	var message, err = NewList(title, description, toggle, footer, section, ctx, client)
	return &waProto.Message{ListMessage: message}, err

}

func NewList(title, description, toggle, footer string, section []*waProto.ListMessage_Section, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.ListMessage, error) {

	var message = &waProto.ListMessage{
		Title:       &title,
		Description: &description,
		ButtonText:  &toggle,
		FooterText:  &footer,
		ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
		Sections:    section,
		ContextInfo: ctx,
	}

	return message, nil
}
