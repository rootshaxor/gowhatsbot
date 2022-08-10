package logs

import (
	"fmt"
	"log"
	"path"
	"runtime"
)

// Get runtime code filename
func CodeFilename(skip int) string {
	_, filename, line, ok := runtime.Caller(skip)
	if ok {
		return fmt.Sprintf(`%s:%d`, path.Base(filename), line)
	} else {
		return ""
	}
}

func Code(me ...interface{}) {
	var mein = []interface{}{
		CodeFilename(2),
	}
	mein = append(mein, me...)

	log.Println(mein...)
}

func Text(me ...interface{}) {
	log.Println(me...)
}
