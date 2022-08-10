package whats

import (
	"context"
	"main/core/media"
	"net/http"
	"os"
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"

	"github.com/disintegration/imaging"
	"go.mau.fi/whatsmeow"
	"google.golang.org/protobuf/proto"
)

func NewImageMessage(data []byte, caption string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {
	var message, err = NewImage(data, caption, ctx, client)
	return &waProto.Message{ImageMessage: message}, err
}

func NewImage(data []byte, caption string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.ImageMessage, error) {
	var mimetype = http.DetectContentType(data)
	if i, err := media.ImageReOrient(data); err != nil {
		return nil, err
	} else {
		i = media.ImageResize(i, 1024, 1024)
		var f imaging.Format = imaging.JPEG
		switch mimetype {
		case "image/gif":
			f = imaging.GIF
		case "image/png":
			f = imaging.PNG
		case "image/webp":
			f = imaging.PNG
		case "image/bmp":
			f = imaging.BMP
		case "image/jpg":
			f = imaging.JPEG
		}

		data, _ = media.ImageToByte(i, f)
	}

	if up, err := client.Upload(context.Background(), data, whatsmeow.MediaImage); err != nil {
		return nil, err
	} else {

		var message = &waProto.ImageMessage{
			Url:           &up.URL,
			Mimetype:      proto.String(mimetype),
			Caption:       &caption,
			FileSha256:    up.FileSHA256,
			FileEncSha256: up.FileEncSHA256,
			FileLength:    &up.FileLength,
			MediaKey:      up.MediaKey,
			DirectPath:    &up.DirectPath,
			ContextInfo:   ctx,
		}

		if strings.HasPrefix(message.GetMimetype(), "image/") {
			if image_info, err := media.ByteToImage(data); err != nil {
				return message, err
			} else {
				message.Width = proto.Uint32(uint32(image_info.Bounds().Max.X))
				message.Height = proto.Uint32(uint32(image_info.Bounds().Max.Y))

				var thumb_w, thumb_h = 100, 100

				img_thumb := media.ImageResize(image_info, thumb_w, thumb_h)

				if data_thumb, err := media.ImageToByte(img_thumb, imaging.JPEG); err != nil {
					return message, err
				} else {

					if thumbnail_up, err := client.Upload(context.Background(), data_thumb, whatsmeow.MediaImage); err != nil {
						return message, err
					} else {
						// message.Caption = proto.String(fmt.Sprintf("%s\nSize : %dK, %dK \nOX, OY : %d, %d \nMX, MY : %d, %d \nPX, PY : %d, %d",
						// 	caption,
						// 	len(data)/1024, len(data_thumb)/1024,
						// 	image_info.Bounds().Max.X, image_info.Bounds().Max.Y,
						// 	img_thumb.Bounds().Max.X, img_thumb.Bounds().Max.Y,
						// 	*message.Width, *message.Height))

						message.ThumbnailDirectPath = &thumbnail_up.DirectPath
						message.ThumbnailEncSha256 = thumbnail_up.FileEncSHA256
						message.ThumbnailSha256 = thumbnail_up.FileSHA256
						message.JpegThumbnail = data_thumb
					}
				}

			}
		}

		return message, err
	}
}

func NewImageMessageFile(filename, caption string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {

	if _, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		if data, err := os.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return NewImageMessage(data, caption, ctx, client)
		}
	}
}
