package log

import (
	"log"
	"os"
)

var stdoutLog, stderrLog *log.Logger

func init() {
	stdoutLog = log.New(os.Stdout, "", log.LstdFlags)
	stderrLog = log.New(os.Stderr, "", log.LstdFlags)
}

func Info(v ...any) {
	stdoutLog.Println(v)
}

func Error(v ...any) {
	stderrLog.Println(v)
}
