package whats

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

func NewProduct(header interface{}, body, footer string, ctx *waProto.ContextInfo) (*waProto.ProductMessage, error) {
	var message = &waProto.ProductMessage{
		Product:          &waProto.ProductMessage_ProductSnapshot{},
		BusinessOwnerJid: proto.String("owner@jid"),
		Catalog:          &waProto.ProductMessage_CatalogSnapshot{},
		Body:             &body,
		Footer:           &footer,
		ContextInfo:      ctx,
	}

	return message, nil
}

func NewProductMessage(header interface{}, body, footer string, ctx *waProto.ContextInfo) (*waProto.Message, error) {

	var message, err = NewProduct(header, body, footer, ctx)

	return &waProto.Message{ProductMessage: message}, err
}
