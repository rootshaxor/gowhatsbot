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
		ContentText: &content,
		FooterText:  &footer,
		Buttons:     buttons,
		ContextInfo: ctx,
	}

	switch hd := header.(type) {
	case *waProto.ButtonsMessage_DocumentMessage:
		message.HeaderType = waProto.ButtonsMessage_DOCUMENT.Enum()
		message.Header = hd

	case *waProto.ButtonsMessage_ImageMessage:
		message.HeaderType = waProto.ButtonsMessage_IMAGE.Enum()
		message.Header = hd

	case *waProto.ButtonsMessage_VideoMessage:
		message.HeaderType = waProto.ButtonsMessage_VIDEO.Enum()
		message.Header = hd

	case *waProto.ButtonsMessage_Text:
		message.HeaderType = waProto.ButtonsMessage_TEXT.Enum()
		message.Header = hd

	case *waProto.ButtonsMessage_LocationMessage:
		message.HeaderType = waProto.ButtonsMessage_LOCATION.Enum()
		message.Header = hd

	}

	return message, nil

}
