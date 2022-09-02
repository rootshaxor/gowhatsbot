package whats

import (
	"context"
	"os"

	waProto "go.mau.fi/whatsmeow/binary/proto"

	"go.mau.fi/whatsmeow"
	"google.golang.org/protobuf/proto"
)

func NewStickerMessage(data []byte, animated bool, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {
	var message, err = NewSticker(data, animated, ctx, client)
	return &waProto.Message{StickerMessage: message}, err
}

func NewSticker(data []byte, animated bool, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.StickerMessage, error) {

	var mimetype = "image/webp"
	if up, err := client.Upload(context.Background(), data, whatsmeow.MediaImage); err != nil {
		return nil, err
	} else {

		var message = &waProto.StickerMessage{
			Url:           &up.URL,
			FileSha256:    up.FileSHA256,
			FileEncSha256: up.FileEncSHA256,
			MediaKey:      up.MediaKey,
			Mimetype:      proto.String(mimetype),
			DirectPath:    &up.DirectPath,
			FileLength:    proto.Uint64(uint64(len(data))),
			IsAnimated:    &animated,
			ContextInfo:   ctx,
		}

		return message, err
	}
}

func NewStickerMessageFile(filename string, animated bool, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {

	if _, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		if data, err := os.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return NewStickerMessage(data, animated, ctx, client)
		}
	}
}
