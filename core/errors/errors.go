package errors

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

type ErrorCode int

const (
	CodePlugin    ErrorCode = 1
	CodeValidator ErrorCode = 2
	CodeCore      ErrorCode = 2
)

// Create new error template
func NewGoWhatsBot(s ErrorCode, n string, c ...string) error {

	_, filename, line, ok := runtime.Caller(int(s))
	if ok {
		return fmt.Errorf("%s:%d %s: %s", path.Base(filename), line, n, strings.Join(c, " "))
	} else {
		return fmt.Errorf("%s: %s", n, c)
	}
}

var GoWhatsBotCommand = NewGoWhatsBot(1, "Command", "Command not found")

// Error template for plugin
func NewPlugin(n string, c ...string) error {
	return NewGoWhatsBot(2, n, c...)
}

// Error template for core
func NewCore(n string, c ...string) error {
	return NewGoWhatsBot(2, n, c...)
}
