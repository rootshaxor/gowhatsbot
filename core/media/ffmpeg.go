package media

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func runCmd(cmdline string, combine bool) ([]byte, error) {
	var command = exec.Command("/bin/sh", "-c", cmdline)

	if combine {
		return command.CombinedOutput()
	} else {
		return command.Output()
	}

}

func Ffmpeg(input_name, output_name string, combine bool, argumenst []string) ([]byte, error) {
	var ffmpeg_cmd = []string{
		"ffmpeg",
		"-i",
		input_name,
	}
	ffmpeg_cmd = append(ffmpeg_cmd, argumenst...)
	ffmpeg_cmd = append(ffmpeg_cmd, output_name)

	var ffmpeg_string = strings.Join(ffmpeg_cmd, " ")

	if output_b, err := runCmd(ffmpeg_string, combine); err != nil {
		return output_b, err
	} else {
		return output_b, nil

	}

}

func FfProbe(input_name string, combine bool, argumenst []string) (map[string]interface{}, []byte, error) {

	argumenst = append([]string{"ffprobe"}, argumenst...)
	argumenst = append(argumenst, input_name)

	var ffprobe_string = strings.Join(argumenst, " ")
	if o, err := runCmd(ffprobe_string, combine); err != nil {
		return nil, o, err
	} else {
		var tmp_f map[string][]map[string]interface{}

		if err := json.Unmarshal(o, &tmp_f); err != nil {
			return nil, o, err
		} else {
			if len(tmp_f["streams"]) > 0 {
				var streams = tmp_f["streams"][len(tmp_f["streams"])-1]
				return streams, o, nil
			} else {
				return nil, o, errors.New("can't extract stream probe")

			}
		}
	}
}
