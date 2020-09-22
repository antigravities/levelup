package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	// LogFile represents the log file.
	file *os.File
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

	msg := fmt.Sprintf("[%s] [%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05 MST"), level, ctmp[len(ctmp)-1], message)

	fmt.Printf(msg)
	file.Write([]byte(msg))
}

// LogOpen opens the log file for writing.
func LogOpen() {
	var err error

	file, err = os.OpenFile("server.log", os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		panic(fmt.Errorf("could not open log files: %v", err))
	}
}

// LogClose closes the log file.
func LogClose() {
	file.Close()
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
