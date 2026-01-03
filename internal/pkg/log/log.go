package log

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"lrcsnc/internal/pkg/types"
)

// logMessage is the struct
// used to communicate inside the logger module
type logMessage struct {
	Type       types.LogLevelType
	ModulePath string
	Msg        string
}

var logger *os.File
var enabled = false
var destination = os.ExpandEnv("$HOME/.local/state/lrcsnc/log")
var level = types.LogLevelInfo
var msgChan = make(chan logMessage)

// Init initializes the logger.
func Init() {
	var err error
	enabled = true

	// Check for the file to be writeable, otherwise panic
	if err := os.MkdirAll(path.Dir(destination), 0744); err != nil {
		Fatal("log/Init", fmt.Sprintf("Log directory could not be created. Either disable logs before launch, or change the directory and try again.\n\nMore:\n%s", err.Error()))
	}
	logger, err = os.Create(os.ExpandEnv(destination))
	if err != nil {
		Fatal("log/Init", fmt.Sprintf("Log file could not be created. Either disable logs before launch, or change the file path and try again.\n\nMore:\n%s", err.Error()))
	}

	// A listener for messages to log
	go func() {
		for msg := range msgChan {
			logger.WriteString(time.Now().Format(time.DateTime) + " [" + msg.ModulePath + "] " + strings.ToUpper(string(msg.Type)) + ": " + msg.Msg + "\n")
			logger.Sync()
		}
	}()
}

func SetDestination(path string) {
	destination = os.ExpandEnv(path)
}

func SetLevel(l types.LogLevelType) {
	level = l
}

// Debug is a log function specifically for detailed debug information.
func Debug(modulePath string, message string) {
	if !enabled || level.ToInt() < types.LogLevelDebug.ToInt() {
		return
	}

	msgChan <- logMessage{types.LogLevelDebug, modulePath, message}
}

// Info is a general-use log function for stuff like start/stop, config updates, etc.
func Info(modulePath string, message string) {
	if !enabled || level.ToInt() < types.LogLevelInfo.ToInt() {
		return
	}

	msgChan <- logMessage{types.LogLevelInfo, modulePath, message}
}

// Warn is a log function specifically for warnings.
func Warn(modulePath string, message string) {
	if !enabled || level.ToInt() < types.LogLevelWarn.ToInt() {
		return
	}

	msgChan <- logMessage{types.LogLevelWarn, modulePath, message}
}

// Error is a log function specifically for non-fatal errors.
func Error(modulePath string, message string) {
	if !enabled || level.ToInt() < types.LogLevelError.ToInt() {
		return
	}

	msgChan <- logMessage{types.LogLevelError, modulePath, message}
}

// Fatal is a log function specifically for fatal errors.
// Immediately exits the app.
// Doesn't write to file: instead uses stdout and stderr to ensure the message gets to user;
// thus it doesn't require for logger to be initialized.
func Fatal(modulePath string, message string) {
	fmt.Fprintln(os.Stdout, "["+modulePath+"] "+strings.ToUpper(string(types.LogLevelFatal))+": "+message)
	// An attempt at graceful shutdown
	//   (I still did not learn contexts...)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
