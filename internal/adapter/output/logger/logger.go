package logger

import (
	"fmt"
	"os"
	"reflect"
)

var isDebug = os.Getenv("DEBUG") == "true"

func Debug(i interface{}, text string) {
	if isDebug {
		fmt.Printf("[Debug][%s] - %s\n", reflect.TypeOf(i).Name(), text)
	}
}

func Info(i interface{}, text string) {
	fmt.Printf("[Info][%s] - %s\n", reflect.TypeOf(i).Name(), text)
}

func Warn(i interface{}, text string) {
	fmt.Printf("[Warn][%s] - %s\n", reflect.TypeOf(i).Name(), text)
}

func Error(i interface{}, text string, err error) {
	var e string
	if err == nil {
		e = "nil"
	} else {
		e = err.Error()
	}

	fmt.Printf("[Error][%s] - %s: %s\n", reflect.TypeOf(i).Name(), text, e)
}
