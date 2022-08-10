package media

import (
	"bytes"
	"io/fs"
	"log"
	"main/core/helper"
	"mime"
	"net/http"
	"os"
	"strconv"
)

type AudioInfo struct {
	Seconds  float64
	Size     int
	Mimetype string
	Bytes    []byte

	SampleRate string
	BitRate    string
	Channels   float64
}

func ConvertToOgg(data []byte) (AudioInfo, []byte, error) {
	var audio AudioInfo

	var mimetype = http.DetectContentType(data)
	var output_name = "ogg_outuput.ogg"
	var input_name = "ogg_input"
	if exts, err := mime.ExtensionsByType(mimetype); err == nil {
		input_name += exts[len(exts)-1]
	}

	if err := os.WriteFile(input_name, data, fs.ModeAppend); err != nil {
		return audio, nil, err
	}

	var args = []string{
		"-acodec libopus",
		"-b:a 64k",
		"-f opus",
		"-y",
	}
	if ouput_ff, err := Ffmpeg(input_name, output_name, true, args); err != nil && !bytes.ContainsAny(ouput_ff, "Error while decoding stream #0:0: Invalid data found when processing input") {

		return audio, ouput_ff, err
	} else {
		// log.Println(string(ouput_ff))
		if b, err := os.ReadFile(output_name); err != nil {
			return audio, nil, err
		} else {
			audio = AudioInfo{
				Size:     len(b),
				Mimetype: http.DetectContentType(b),
				Bytes:    b,
			}

			var args_probe = []string{
				"-v quiet",
				"-print_format json",
				"-show_streams",
			}

			if map_stream, output, err := FfProbe(output_name, true, args_probe); err != nil {
				log.Println(err, string(output))
				return audio, output, nil
			} else {

				if tmpvar, ok := map_stream["duration"]; ok {
					if tmpval, err := strconv.ParseFloat(tmpvar.(string), 32); err == nil {
						audio.Seconds = tmpval
					} else {
						log.Printf("duration %v", helper.GetType(tmpvar))
					}
				}

				if tmpvar, ok := map_stream["sample_rate"]; ok {
					if tmpval, ok := tmpvar.(string); ok {
						audio.SampleRate = tmpval
					} else {
						log.Printf("sample_rate %v", helper.GetType(tmpvar))
					}
				}

				if tmpvar, ok := map_stream["bit_rate"]; ok {
					if tmpval, ok := tmpvar.(string); ok {
						audio.BitRate = tmpval
					} else {
						log.Printf("bit_rate %v", helper.GetType(tmpvar))
					}
				}

				if tmpvar, ok := map_stream["channels"]; ok {
					if tmpval, ok := tmpvar.(float64); ok {
						audio.Channels = tmpval
					} else {
						log.Printf("channels %v", helper.GetType(tmpvar))
					}
				}
			}

			return audio, ouput_ff, nil
		}
	}
}
