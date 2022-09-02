package media

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imageorient"
	"github.com/disintegration/imaging"
	// "github.com/harukasan/go-libwebp/webp"
	// "github.com/tidbyt/go-libwebp/webp"
)

type MediaInfo struct {
	Size         int
	Bytes        []byte
	Seconds      int
	Mimetype     string
	JpegThumbail []byte
	Width        int
	Height       int
}

type StreamsInfo struct {
	Streams []map[string]any `json:"streams"`
}

func GifToMP4(data []byte) (MediaInfo, error) {
	var tmpdir = "tmp"
	var outname = path.Join(tmpdir, "out.mp4")
	var inname = path.Join(tmpdir, "out.gif")
	var res MediaInfo

	if _, err := os.Stat("./tmp"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("./tmp", os.ModeAppend)
	}

	if err := os.WriteFile(inname, data, fs.ModeAppend); err != nil {
		return res, err
	} else {

		ffmpeg_cmd := `ffmpeg -i %s -movflags faststart -pix_fmt "yuv420p" -vf "scale=trunc(iw/2)*2:trunc(ih/2)*2" %s -y`
		ffmpeg_cmd = fmt.Sprintf(ffmpeg_cmd, inname, outname)

		if _, err := runCmd(ffmpeg_cmd, true); err != nil {
			return res, err
		} else {

			if res.Bytes, err = os.ReadFile(outname); err != nil {
				return res, err
			}

			res.Size = len(res.Bytes)
			res.Mimetype = http.DetectContentType(res.Bytes)

			cmd_streams := `ffprobe -v quiet -print_format json -show_streams ` + outname
			if streamsb, err := runCmd(cmd_streams, true); err != nil {
				return res, err
			} else {
				var streams StreamsInfo
				if err := json.Unmarshal(streamsb, &streams); err == nil {
					res.Height = int(streams.Streams[0]["height"].(float64))
					res.Width = int(streams.Streams[0]["width"].(float64))

					if flt, err := strconv.ParseFloat(streams.Streams[0]["duration"].(string), 64); err == nil {
						res.Seconds = int(flt)
					}
				}
			}

			thumbname := path.Join(tmpdir, "out.png")
			cmd_thumbnail := fmt.Sprintf("ffmpeg -i %s -ss 00:00:01.000 -vframes 1 -y %s", outname, thumbname)
			if _, err := runCmd(cmd_thumbnail, true); err != nil {
				return res, err
			} else {
				if tmb, err := os.ReadFile(thumbname); err != nil {
					return res, err
				} else {
					res.JpegThumbail = tmb
				}
			}

			os.Remove(thumbname)
			os.Remove(inname)
			os.Remove(outname)
		}

	}

	return res, nil
}

func ImageToByte(img image.Image, f imaging.Format) ([]byte, error) {
	var b = bytes.NewBuffer(nil)

	err := imaging.Encode(b, img, f)

	return b.Bytes(), err
}

func ByteToImage(b []byte) (image.Image, error) {
	return imaging.Decode(bytes.NewReader(b))
}

func ImageResize(i image.Image, width, height int) image.Image {
	return imaging.Fit(i, width, height, imaging.ResampleFilter{})
}

func ImageReOrient(b []byte) (image.Image, error) {

	if x, _, err := imageorient.Decode(bytes.NewReader(b)); err != nil {
		return x, err
	} else {
		return x, err
	}
}

func ByteToWebp(b []byte) (image.Image, error) {

	if ii, err := jpeg.Decode(bytes.NewBuffer(b)); err == nil {
		if bb, err := ImageToByte(ii, imaging.JPEG); err == nil {
			b = bb
		}
	}

	if i, err := webp.Decode(bytes.NewReader(b)); err != nil {
		return nil, err
	} else {

		return ImageToWebp(i)
	}
}

func ImageToWebp(img image.Image) (image.Image, error) {
	var b bytes.Buffer
	if err := webp.Encode(&b, img, &webp.Options{Lossless: true}); err != nil {
		return nil, err
	}

	return ByteToImage(b.Bytes())
}

func ByteToWebpByte(b []byte) ([]byte, error) {

	if bb, err := ByteImageToByte(b); err == nil {
		b = bb
	}

	if wp, err := webp.Decode(bytes.NewReader(b)); err != nil {
		return nil, err
	} else {

		var buff bytes.Buffer

		if err := webp.Encode(&buff, wp, &webp.Options{Lossless: true}); err != nil {
			return nil, err
		} else {

			log.Println(http.DetectContentType(buff.Bytes()))
			return buff.Bytes(), nil
		}

	}
}

func WebpToByte(w image.Image) ([]byte, error) {
	var b bytes.Buffer

	if err := webp.Encode(&b, w, &webp.Options{Lossless: true}); err != nil {
		return b.Bytes(), err
	}

	return b.Bytes(), nil
}

func ByteImageToByte(b []byte) ([]byte, error) {
	if i, err := ByteToImage(b); err != nil {
		return nil, err
	} else {
		i = ImageResize(i, 512, 512)

		return WebpToByte(i)
	}
}

func ByteToWebpAnimationByte(b []byte) ([]byte, error) {

	var tempname_input = "temp_in"
	var tempname_output = "temp_out.webp"
	var cmd_ffmpeg = []string{
		"ffmpeg",
		"-t 5",
		"-i " + tempname_input,
		"-filter:v",
		"fps=fps=10",
		fmt.Sprintf(`-vf "%s"`, "format=rgba,setsar=1"),
		"-compression_level 0",
		// "-q:v 30",
		"-loop 0",
		"-preset picture",
		"-an",
		"-vsync 0",
		"-s 512:512",

		"-y",
		tempname_output,
	}

	for _, filename := range []string{tempname_input, tempname_output} {
		if _, err := os.Stat(filename); err == nil {
			os.Remove(filename)
		}

	}

	if err := os.WriteFile(tempname_input, b, fs.ModeAppend); err != nil {
		return nil, err
	}

	var cmd_ffmpeg_str = strings.Join(cmd_ffmpeg, " ")
	if resp, err := runCmd(cmd_ffmpeg_str, true); err != nil {
		return resp, err
	} else {
		if _, err := os.Stat(tempname_output); err != nil {
			return nil, err
		} else {
			return os.ReadFile(tempname_output)

		}
	}
}
