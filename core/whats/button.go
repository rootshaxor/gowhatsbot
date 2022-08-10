package whats

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func NewButtonsMessage(header interface{}, content, footer string, buttons []*waProto.ButtonsMessage_Button, ctx *waProto.ContextInfo) (*waProto.Message, error) {
	var message, err = NewButtons(header, content, footer, buttons, ctx)
	return &waProto.Message{ButtonsMessage: message}, err
}

func NewButtons(header interface{}, content, footer string, buttons []*waProto.ButtonsMessage_Button, ctx *waProto.ContextInfo) (*waProto.ButtonsMessage, error) {

	var message = &waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_EMPTY.Enum(),
		ContentText: &content,
		FooterText:  &footer,
		Buttons:     buttons,
		ContextInfo: ctx,
	}

	if header != nil {
		message.Header = header.(*waProto.ButtonsMessage_Text)
	}

	return message, nil

}
