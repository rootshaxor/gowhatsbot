package validators

import "go.mau.fi/whatsmeow/binary/proto"

// check if the message can be downloaded
func IsDownloadable(e *proto.Message) bool {
	return IsDocument(e) || IsAudio(e) || IsGraphics(e)
}

// check if the message is a document
func IsDocument(e *proto.Message) bool {
	return e.DocumentMessage != nil
}

// check if the message is a voice
func IsAudio(e *proto.Message) bool {
	return e.AudioMessage != nil
}

// check if the message contains sound media
func IsAudioAble(e *proto.Message) bool {
	return IsAudio(e) || IsVideo(e)
}

// check if the message is a video
func IsVideo(e *proto.Message) bool {
	return e.VideoMessage != nil
}

// check if the message is a graphics media
func IsGraphics(e *proto.Message) bool {
	return e.ImageMessage != nil || e.VideoMessage != nil || e.StickerMessage != nil
}
