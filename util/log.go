package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func log(level string, message string) {
	caller := ""

	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		caller = "(unknown)"
	} else {
		caller = runtime.FuncForPC(pc).Name()
	}

	ctmp := strings.Split(caller, "/")

	fmt.Printf("[%s] [%s] %s\n", level, ctmp[len(ctmp)-1], message)
}

// Info logs a message with this log level
func Info(message string) {
	log("INFO", message)
}

// Debug logs a message with this log level.
// The DEBUG environment variable must be non-empty in order for messages to show.
func Debug(message string) {
	if os.Getenv("DEBUG") == "" {
		return
	}

	log("DEBUG", message)
}

// Warn logs a message with this log level.
func Warn(message string) {
	log("WARN", message)
}
