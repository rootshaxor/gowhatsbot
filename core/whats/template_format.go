package whats

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func NewFormatTemplateMessage(title interface{}, content, footer, template_id string, buttons []*waProto.HydratedTemplateButton, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {

	if message, err := NewHydratedTemplate(title, content, footer, template_id, buttons, ctx, client); err != nil {
		return nil, err
	} else {
		return &waProto.Message{TemplateMessage: message}, err
	}

}

func NewFormatTemplate(title interface{}, content, footer, template_id string, buttons []*waProto.HydratedTemplateButton, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.TemplateMessage, error) {

	var message = &waProto.TemplateMessage{
		ContextInfo:      ctx,
		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{},
		Format:           &waProto.TemplateMessage_FourRowTemplate_{},
	}

	switch the_title := title.(type) {
	case *waProto.TemplateMessage_HydratedFourRowTemplate_HydratedTitleText:
		message.HydratedTemplate.Title = the_title

	case *waProto.TemplateMessage_HydratedFourRowTemplate_ImageMessage:
		message.HydratedTemplate.Title = the_title

	case *waProto.TemplateMessage_HydratedFourRowTemplate_DocumentMessage:
		message.HydratedTemplate.Title = the_title

	case *waProto.TemplateMessage_HydratedFourRowTemplate_VideoMessage:
		message.HydratedTemplate.Title = the_title

	case *waProto.TemplateMessage_HydratedFourRowTemplate_LocationMessage:
		message.HydratedTemplate.Title = the_title

	}

	return message, nil
}
