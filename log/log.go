package log

import (
	"bytes"
	"fmt"
	"os"

	"github.com/morikuni/aec"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	NONE
)

var logLevel = DEBUG
var debugColor = aec.YellowF
var infoColor = aec.GreenF
var errorColor = aec.RedF

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func Debug(is ...interface{}) {
	if logLevel > DEBUG {
		return
	}
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, debugColor.Apply("[DEBUG] "))
	fmt.Fprintln(buf, is...)
	os.Stderr.Write(buf.Bytes())
}

func Info(is ...interface{}) {
	if logLevel > INFO {
		return
	}
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, infoColor.Apply("[INFO] "))
	fmt.Fprintln(buf, is...)
	os.Stderr.Write(buf.Bytes())
}

func Error(is ...interface{}) {
	if logLevel > ERROR {
		return
	}
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, errorColor.Apply("[ERROR] "))
	fmt.Fprintln(buf, is...)
	os.Stderr.Write(buf.Bytes())
}
