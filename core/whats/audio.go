package whats

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-audio/wav"
	"github.com/tcolgate/mp3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

func NewAudioMessage(data []byte, seconds int, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {
	var message, err = NewAudio(data, seconds, ctx, client)
	return &waProto.Message{AudioMessage: message}, err

}

func NewAudio(data []byte, seconds int, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.AudioMessage, error) {

	var mimetype = http.DetectContentType(data)
	log.Println(mimetype)

	var ptt = strings.Contains(mimetype, "ogg")
	if ptt {
		mimetype = "audio/ogg; codecs=opus"
	}

	if up, err := client.Upload(context.Background(), data, whatsmeow.MediaAudio); err != nil {
		return nil, err
	} else {
		var message = &waProto.AudioMessage{
			Url:           &up.URL,
			Mimetype:      &mimetype,
			FileLength:    &up.FileLength,
			FileSha256:    up.FileSHA256,
			FileEncSha256: up.FileEncSHA256,
			MediaKey:      up.MediaKey,
			DirectPath:    &up.DirectPath,
			Ptt:           &ptt,
			ContextInfo:   ctx,
			Seconds:       proto.Uint32(uint32(seconds)),
		}

		switch mimetype {
		case "audio/mpeg":
			{
				var decoder = mp3.NewDecoder(bytes.NewReader(data))
				var frame mp3.Frame
				var total_sec float64
				for {
					if err := decoder.Decode(&frame, nil); err != nil {
						if err == io.EOF {
							break
						}
						log.Println(err)
					} else {

						total_sec += frame.Duration().Seconds()
					}
				}

				if total_sec > 0 {
					message.Seconds = proto.Uint32(uint32(total_sec))
				}
			}
		case "audio/wave":
			{
				var decoder = wav.NewDecoder(bytes.NewReader(data))
				if total_sec, err := decoder.Duration(); err == nil {
					message.Seconds = proto.Uint32(uint32(total_sec))
				}
			}

		}

		return message, err
	}

}

func NewAudioMessageFile(filename string, ctx *waProto.ContextInfo, client *whatsmeow.Client) (*waProto.Message, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	} else {
		if data, err := os.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return NewAudioMessage(data, 0, ctx, client)
		}
	}
}
