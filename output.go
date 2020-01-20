// graylogger was made to provide easy-to-analyze log messages.
// For sending GELF messages it uses the: https://github.com/Devatoria/go-graylog package.
package graylogger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Devatoria/go-graylog"
)

var functions Functions

// Init initializes the logger instance
type Init struct {
	GraylogHost     string        // Host name of the Graylog server
	GraylogPort     int           // Port number of the Graylog server
	GraylogProvider string        // The Name of the service which generates the log messages or sends logs into Graylog
	GraylogProtocol Transport     // The name of the transport protocol: the way we send GELF messages (TransportTCP or TransportUDP)
	GraylogTimeout  time.Duration // Optional, it declares the maximum amount of time a dial will wait for a connection to complete.

	LogEnv   string   // Environment of the service: dev / test / prod
	LogLevel LogLevel // It can be: debug, info, warning, error
	LogColor bool     // Turn on/off colored log output. Color output can be useful during development, but it is recommended to turn it off in production environment.
}

type (
	// Transport declares the the way we send GELF messages (TransportTCP or TransportUDP).
	Transport string

	// LogLevel defines the log levels that can be entered.
	LogLevel string
)

// Functions provide a different kind of logging writers which controlled by log level.
type Functions struct {
	Debug   *log.Logger // Designates fine-grained informational events that are most useful to debug an application.
	Info    *log.Logger // Designates informational messages that highlight the progress of the application at coarse-grained level.
	Warning *log.Logger // Designates potentially harmful situations.
	Error   *log.Logger // Designates error events that might still allow the application to continue running.
	Fatal   *log.Logger // Designates very severe error events that will presumably lead the application to abort.
}

// GrayLogger holds the needed data to use the functions of this package.
//  - initData -> the data with which the package was initialized
//  - functions -> logger functions: Debug, Info, Warning, Error, Fatal
//  - level -> log level converted to integer
//  - fileName -> set by CaptureOutput() function, provides the filename where output can be saved
//  - fileOpen ->  set by CaptureOutput() function, holds *os.File
//  - graylog -> set by connect() function, it represents an established graylog connection
type GrayLogger struct {
	initData  Init
	functions Functions
	level     int
	fileName  string
	fileOpen  *os.File
	graylog   *graylog.Graylog
}

const (
	// TransportTCP  is connection-oriented, and a connection between client
	// and server is established (passive open) before data can be sent.
	TransportTCP Transport = "tcp"

	// TransportUDP is suitable for purposes where error checking
	// and correction are either not necessary or are performed in the application
	TransportUDP Transport = "udp"

	// LevelDebug logs everything
	LevelDebug    LogLevel = "debug"
	levelDebugNum int      = 7

	// LevelInfo logs Info, Warnings and Errors
	LevelInfo    LogLevel = "info"
	levelInfoNum int      = 6

	// LevelWarning logs Warning and Errors
	LevelWarning    LogLevel = "warning"
	levelWarningNum int      = 4

	// LevelError logs just Errors
	LevelError    LogLevel = "error"
	levelErrorNum int      = 3

	// LevelFatal logs just Fatal errors and aborts the application
	LevelFatal    LogLevel = "fatal"
	levelFatalNum int      = 2
)

// New configures the logging writers.
func New(init Init) *GrayLogger {
	logLevel := logLevelToInt(init.LogLevel)
	init.setLogLevelFunctions(setLogLevelHandlers(logLevel, os.Stdout))

	l := &GrayLogger{
		initData:  init,
		functions: functions,
		level:     logLevel,
	}

	if err := l.initData.LogLevel.validateLogLevel(); err != nil && l.isSetGraylogObligatoryFields() {
		l.Fatal(err)
	}

	if err := l.initData.GraylogProtocol.validateTransport(); err != nil && l.isSetGraylogObligatoryFields() {
		l.Fatal(err)
	}

	if l.initData.GraylogTimeout == 0 {
		l.initData.GraylogTimeout = graylogTimeout
	}

	return l
}

// Tracking provides debug information about function invocations
// on the calling goroutine's stack:
//  - File (where Tracking was called)
//  - Line (where function was called)
//  - Function (name of the function)
func Tracking(depth int) TrackInfo {
	return getTrackingInfo(depth)
}

// DiscardOutput discards all level outputs and prevents sending messages to Graylog as well.
func (g *GrayLogger) DiscardOutput() {
	g.functions.Debug.SetOutput(ioutil.Discard)
	g.functions.Info.SetOutput(ioutil.Discard)
	g.functions.Warning.SetOutput(ioutil.Discard)
	g.functions.Error.SetOutput(ioutil.Discard)
	g.functions.Fatal.SetOutput(ioutil.Discard)

	g.level = 0
}

// ResetLogger allows StdOut with the initialized log level and allows sending messages to Graylog as well.
func (g *GrayLogger) ResetLogger() *GrayLogger {
	return New(g.initData)
}

// Debug writes Info to stdOut and sends GELF message to Graylog.
func (g *GrayLogger) Debug(keysAndValues ...interface{}) {
	g.functions.Debug.Println(g.formatLogLine(keysAndValues...))
	g.SendGELF(levelDebugNum, keysAndValues...)
}

// Info writes Info to stdOut and sends GELF message to Graylog.
func (g *GrayLogger) Info(keysAndValues ...interface{}) {
	g.functions.Info.Println(g.formatLogLine(keysAndValues...))
	g.SendGELF(levelInfoNum, keysAndValues...)
}

// Warning writes Warning to stdOut and sends GELF message to Graylog.
func (g *GrayLogger) Warning(keysAndValues ...interface{}) {
	g.functions.Warning.Println(g.formatLogLine(keysAndValues...))
	g.SendGELF(levelWarningNum, keysAndValues...)
}

// LogWarningIfErr only writes Warning to stdOut, if err doesn't nil.
// It also sends GELF message to Graylog, if it possible.
func (g *GrayLogger) LogWarningIfErr(err error) {
	if err != nil {
		g.functions.Warning.Println(g.formatLogLine(getTrackingInfo(1).Function, err))
		g.SendGELF(levelWarningNum, getTrackingInfo(1).Function, err)
	}
}

// Error writes Error to stdOut and sends GELF message to Graylog.
func (g *GrayLogger) Error(keysAndValues ...interface{}) {
	g.functions.Error.Println(g.formatLogLine(keysAndValues...))
	g.SendGELF(levelErrorNum, keysAndValues...)
}

// LogErrorIfErr only writes Error to stdOut, if err doesn't nil.
// It also sends GELF message to Graylog, if it possible.
func (g *GrayLogger) LogErrorIfErr(err error) {
	if err != nil {
		g.functions.Error.Println(g.formatLogLine(getTrackingInfo(1).Function, err))
		g.SendGELF(levelErrorNum, getTrackingInfo(1).Function, err)
	}
}

// Fatal writes Error to stdOut and exit with exit code 1, if err doesn't nil.
// It also sends GELF message to Graylog, if it possible.
func (g *GrayLogger) Fatal(err error) {
	if err != nil {
		g.functions.Fatal.Println(g.formatLogLine(getTrackingInfo(1).Function, err))
		g.SendGELF(levelFatalNum, getTrackingInfo(1).Function, err)
		os.Exit(1)
	}
}

// ReturnWithError writes Error to stdOut and returns with that error message at the same time.
// It also sends GELF message to Graylog, if it possible.
func (g *GrayLogger) ReturnWithError(keysAndValues ...interface{}) error {
	g.functions.Error.Println(g.formatLogLine(keysAndValues...))
	g.SendGELF(levelErrorNum, keysAndValues...)
	return fmt.Errorf(prettifyKeyVal(keyValToSlice(keysAndValues...)))
}

// GetInit returns with the initial logger data.
func (g *GrayLogger) GetInit() Init {
	return g.initData
}

// IsAllowedOutput tells that StdOut is allowed or discarded on all log levels.
func (g *GrayLogger) IsAllowedOutput() bool {
	return strings.Contains(fmt.Sprint(
		g.functions.Debug.Writer(),
		g.functions.Info.Writer(),
		g.functions.Warning.Writer(),
		g.functions.Error.Writer(),
		g.functions.Fatal.Writer(),
	), "x")
}

// GetLogLevel returns with the set log level (int, string).
func (g *GrayLogger) GetLogLevel() (int, string) {
	return logLevelToInt(g.GetInit().LogLevel), fmt.Sprint(g.GetInit().LogLevel)
}

// CaptureOutput will redirect logger output to a given fileName.
func (g *GrayLogger) CaptureOutput(fileName string) *GrayLogger {
	g.fileName = fileName
	_ = os.Remove(fileName)

	f, _ := os.OpenFile(g.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	setLogLevelHandlers(g.level, f).setOutput(g)

	return g
}

// SaveOutput will close the file, what we set in CaptureOutput() function
// and re-set the output of all logger functions .
func (g *GrayLogger) SaveOutput() {
	_ = g.fileOpen.Close()
	g.initData.setLogLevelFunctions(setLogLevelHandlers(g.level, os.Stdout))
}

// GetOutput reads the file content what we set in the CaptureOutput() function.
func (g *GrayLogger) GetOutput() string {
	b, _ := ioutil.ReadFile(g.fileName)
	return string(b)
}

// PrintOutput prints out the file content what we set in the CaptureOutput() function.
func (g *GrayLogger) PrintOutput() {
	fmt.Println(strings.TrimSuffix(g.GetOutput(), "\n"))
}
