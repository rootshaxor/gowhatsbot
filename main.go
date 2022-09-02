package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"main/plugins"
	"os"
	"path"
	"reflect"
	"runtime"
)

var GoWhatsClient IGoWhatsBot
var WhatsConfig Config

func main() {
	generatePlugins()

	if cfg, err := os.ReadFile("gowhatsbot.json"); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(cfg, &WhatsConfig); err != nil {
			panic(err)
		}
	}

	GoWhatsClient = NewGoWhatsBot(
		WhatsConfig,
		eventHandler,
	)

	GoWhatsClient.Run()
	GoWhatsClient.Stop()
}

func eventHandler(e interface{}) {
	plugins.PluginExecutor(e, GoWhatsClient.GetClient())

}

type Pkg struct{}

func generatePlugins() {

	if _, filename, _, ok := runtime.Caller(0); ok {
		var dirname = "plugins"
		var current_dir = path.Dir(filename)
		var plugin_dir = fmt.Sprintf(`%s/%s`, current_dir, dirname)
		var autload_name = fmt.Sprintf(`%s/%s`, current_dir, "autoload.go")
		var packs []string
		var autload_byte []byte
		if a, err := os.ReadFile(autload_name); err == nil {
			autload_byte = a
		}

		if entries, err := os.ReadDir(plugin_dir); err == nil {
			for _, entry := range entries {
				var package_path = fmt.Sprintf(`"%s/plugins/%s"`, reflect.TypeOf(Pkg{}).PkgPath(), entry.Name())
				if entry.IsDir() && entry.Name() != dirname {
					packs = append(packs, package_path)

				}
			}

			if len(packs) > 0 {
				var buff_s bytes.Buffer

				buff_s.WriteString(fmt.Sprintf("package %s\n\n", reflect.TypeOf(Pkg{}).PkgPath()))
				buff_s.WriteString("import (\n")
				for _, p := range packs {
					buff_s.WriteString("\t _ " + p + "\n")
				}

				buff_s.WriteString(")\n")
				if buff_s.Len() > 0 && len(autload_byte) != buff_s.Len() {
					if err := os.WriteFile(autload_name, buff_s.Bytes(), fs.ModeAppend); err != nil {
						log.Println(err)
					}
					log.Println(path.Base(autload_name), "has been updated. Restarting App required.")
				}
			}
		}
	}
}
