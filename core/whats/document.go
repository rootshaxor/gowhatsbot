package whats

import (
	"context"
	"net/http"
	"os"
	"path"

	waProto "go.mau.fi/whatsmeow/binary/proto"

	"go.mau.fi/whatsmeow"
	"google.golang.org/protobuf/proto"
)

func NewDocumentMessage(data []byte, filename string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {
	var message, err = NewDocument(data, filename, ctx, client)
	return &waProto.Message{DocumentMessage: message}, err
}

func NewDocument(data []byte, filename string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.DocumentMessage, error) {
	if up, err := client.Upload(context.Background(), data, whatsmeow.MediaDocument); err != nil {
		return nil, err
	} else {

		message := &waProto.DocumentMessage{
			Url:           &up.URL,
			Mimetype:      proto.String(http.DetectContentType(data)),
			Title:         &filename,
			FileSha256:    up.FileSHA256,
			FileLength:    &up.FileLength,
			MediaKey:      up.MediaKey,
			FileName:      &filename,
			FileEncSha256: up.FileEncSHA256,
			DirectPath:    &up.DirectPath,
			ContextInfo:   ctx,
		}

		return message, err
	}
}

func NewDocumentMessageFile(filename string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		if data, err := os.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return NewDocumentMessage(data, path.Base(filename), ctx, client)
		}
	}
}
