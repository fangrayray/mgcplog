package mgcplog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&MGCPJSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// LogrusLogger : logrus wrapper
type LogrusLogger struct {
	Config *LogConfiguration
}

// Info : log with info level
func (ll *LogrusLogger) Info(msg string, sid string) {
	err := ll.logWithLogLevel(Info, msg, sid)
	if err != nil {
		log2Console(err)
	}
}

// Warn : log with warn level
func (ll *LogrusLogger) Warn(msg string, sid string) {
	err := ll.logWithLogLevel(Warn, msg, sid)
	if err != nil {
		log2Console(err)
	}
}

// Error : log with error level
func (ll *LogrusLogger) Error(msg string, sid string) {
	err := ll.logWithLogLevel(Error, msg, sid)
	if err != nil {
		log2Console(err)
	}
}

// Panic : log with panic level
func (ll *LogrusLogger) Panic(msg string, sid string) {
	err := ll.logWithLogLevel(Panic, msg, sid)
	if err != nil {
		log2Console(err)
	}
}

// Fatal : log with fatal level
func (ll *LogrusLogger) Fatal(msg string, sid string) {
	err := ll.logWithLogLevel(Fatal, msg, sid)
	if err != nil {
		log2Console(err)
	}
}

func (ll *LogrusLogger) logWithLogLevel(lv LogLevel, msg string, sid string) error {
	fName, fLine, fcName := trace()

	le := &LogEntity{
		Timestamp:    time.Now().Format("2006-01-02T15:04:05Z"),
		Level:        lv,
		Service:      ll.Config.ServiceName,
		SessionID:    sid,
		FunctionName: fcName,
		FileName:     fmt.Sprintf("%s:%d", fName, fLine),
		Message:      msg,
	}
	err := ll.log2File(le)
	if err != nil {
		return err
	}
	return nil
}

func trace() (string, int, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
	return frame.File, frame.Line, frame.Function

}

func (ll *LogrusLogger) log2File(le *LogEntity) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	var f *os.File
	var err error

	if ll.Config.LogFile != "" {
		f, err = os.OpenFile(ll.Config.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		defer f.Close()
		mw := io.MultiWriter(f, os.Stdout)
		log.SetOutput(mw)
	}
	fields := &log.Fields{
		"time":          le.Timestamp,
		"level":         le.Level,
		"release_name":  le.Service,
		"request_id":    le.SessionID,
		"file_name":     le.FileName,
		"function_name": le.FunctionName,
		"message":       le.Message,
	}
	entry := log.WithFields(*fields)
	switch le.Level {
	case Info:
		entry.Info(le.Message)
	case Warn:
		entry.Warn(le.Message)
	case Error:
		entry.Error(le.Message)
	case Fatal:
		entry.Fatal(le.Message)
	case Panic:
		entry.Panic(le.Message)
	default:
		return fmt.Errorf("not supported log level")
	}
	return nil
}

func log2Console(err error) {
	log.SetOutput(os.Stdout)
	log.Error(fmt.Sprintf("failed to log, %v", err))
}

// MGCPJSONFormatter : MGCP JSONFormater
type MGCPJSONFormatter struct {
}

// Format :
func (f *MGCPJSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	// Note this doesn't include Time, Level and Message which are available on
	// the Entry. Consult `godoc` on information about those fields or read the
	// source of the official loggers.
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
