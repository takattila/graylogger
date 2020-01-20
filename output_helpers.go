package graylogger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"runtime"
	"strings"

	c "github.com/bclicn/color"
)

const (
	colorRed    = "red"
	colorGreen  = "green"
	colorYellow = "yellow"
	colorBlue   = "blue"
	colorPurple = "purple"
	colorGray   = "gray"
)

// TrackInfo holds debug information about the caller's:
//   - File name
//   - Line number
//   - Function name
type TrackInfo struct {
	File     string
	Line     string
	Function string
}

// logLevelHandlers holds the available logging handler functions.
type logLevelHandlers struct {
	debug io.Writer
	info  io.Writer
	warn  io.Writer
	error io.Writer
	fatal io.Writer
}

// setLogLevelFunctions initializes the output of the logger functions.
func (i Init) setLogLevelFunctions(h logLevelHandlers) {
	functions.Debug = log.New(h.debug, i.colorOut(colorGreen, "[DEBUG] "), log.Ldate|log.Ltime)
	functions.Info = log.New(h.info, i.colorOut(colorBlue, "[INFO] "), log.Ldate|log.Ltime)
	functions.Warning = log.New(h.warn, i.colorOut(colorPurple, "[WARNING] "), log.Ldate|log.Ltime)
	functions.Error = log.New(h.error, i.colorOut(colorRed, "[ERROR] "), log.Ldate|log.Ltime)
	functions.Fatal = log.New(h.fatal, i.colorOut(colorYellow, "[FATAL] "), log.Ldate|log.Ltime)
}

// setLogLevelHandlers decides whether the output of the logger functions should be discarded or not.
func setLogLevelHandlers(logLevel int, w io.Writer) logLevelHandlers {
	debugHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard
	fatalHandle := ioutil.Discard

	if logLevel == levelDebugNum {
		debugHandle = w
		infoHandle = w
		warnHandle = w
		errorHandle = w
		fatalHandle = w
	}

	if logLevel == levelInfoNum {
		infoHandle = w
		warnHandle = w
		errorHandle = w
		fatalHandle = w
	}

	if logLevel == levelWarningNum {
		warnHandle = w
		errorHandle = w
		fatalHandle = w
	}

	if logLevel == levelErrorNum {
		errorHandle = w
		fatalHandle = w
	}

	if logLevel == levelFatalNum {
		fatalHandle = w
	}

	return logLevelHandlers{
		debug: debugHandle,
		info:  infoHandle,
		warn:  warnHandle,
		error: errorHandle,
		fatal: fatalHandle,
	}
}

// setOutput will set logger output to a given io.Writer
func (h logLevelHandlers) setOutput(g *GrayLogger) {
	g.functions.Debug.SetOutput(h.debug)
	g.functions.Info.SetOutput(h.info)
	g.functions.Warning.SetOutput(h.warn)
	g.functions.Error.SetOutput(h.error)
	g.functions.Fatal.SetOutput(h.fatal)
}

// validateTransport checks that given transport protocol is valid or not.
func (t Transport) validateTransport() error {
	switch t {
	case TransportTCP, TransportUDP:
		return nil
	}
	return fmt.Errorf("invalid transport protocol given: %s", t)
}

// validateLogLevel checks that given logging level is valid or not.
func (l LogLevel) validateLogLevel() error {
	switch l {
	case LevelDebug, LevelInfo, LevelWarning, LevelError, LevelFatal:
		return nil
	}
	return fmt.Errorf("invalid logging level given: %s", l)
}

// logLevelToInt converts the name of the log level to its integer value.
// For example:
//  debug -> 7
//  info -> 6
func logLevelToInt(levelStr LogLevel) int {
	switch strings.ToLower(fmt.Sprint(levelStr)) {
	case string(LevelDebug):
		return levelDebugNum
	case string(LevelInfo):
		return levelInfoNum
	case string(LevelWarning):
		return levelWarningNum
	case string(LevelError):
		return levelErrorNum
	case string(LevelFatal):
		return levelFatalNum
	default:
		return levelDebugNum
	}
}

// logLevelToString converts given integer log level value to its name.
// For example:
//  7 -> debug
//  6 -> info
func logLevelToString(levelNum int) string {
	switch levelNum {
	case levelDebugNum:
		return string(LevelDebug)
	case levelInfoNum:
		return string(LevelInfo)
	case levelWarningNum:
		return string(LevelWarning)
	case levelErrorNum:
		return string(LevelError)
	case levelFatalNum:
		return string(LevelFatal)
	default:
		return string(LevelDebug)
	}
}

// getTrackingInfo provides debug information about function invocations
// on the calling goroutine's stack:
//  - File (where Tracking was called)
//  - Line (where function was called)
//  - Function (name of the function)
func getTrackingInfo(depth int) TrackInfo {
	_, fileName, line, _ := runtime.Caller(depth + 1)
	return TrackInfo{
		File:     fetchNameFromPath(fileName),
		Line:     fmt.Sprintf("%d", line),
		Function: getFuncName(depth + 1),
	}
}

// getFuncName returns with the caller's function name
func getFuncName(depth int) string {
	pc, _, _, _ := runtime.Caller(depth + 1)
	me := runtime.FuncForPC(pc)
	if me == nil {
		return "unknown"
	}
	return fetchNameFromPath(me.Name())
}

// fetchNameFromPath extracts the name of a function from a path.
func fetchNameFromPath(fileName string) string {
	for i := len(fileName) - 1; i > 0; i-- {
		if fileName[i] == '/' {
			return fileName[i+1:]
		}
	}
	return fileName
}

// colorOut decides whether to have a color output or not.
// If Init.LogColor is set to false, the color output will be disabled.
func (i *Init) colorOut(color, text string) string {
	if i.LogColor {
		switch strings.ToLower(color) {
		case colorRed:
			return c.LightRed(text)
		case colorGreen:
			return c.Green(text)
		case colorYellow:
			return c.LightYellow(text)
		case colorBlue:
			return c.Blue(text)
		case colorPurple:
			return c.LightPurple(text)
		case colorGray:
			return c.DarkGray(text)
		}
	}
	return c.White(text)
}

// formatLogLine provide a formatted log line.
// For example:
//  [DEBUG] 2020/01/21 12:53:09 [file: example.go line: 39 function: main.main] [Debug :: 10]
func (g *GrayLogger) formatLogLine(keysAndValues ...interface{}) string {
	tr := getTrackingInfo(2)
	return fmt.Sprintf("%s [%s]",
		g.initData.colorOut(colorGray, fmt.Sprintf("[file: %s line: %s function: %s]",
			tr.File,
			tr.Line,
			tr.Function)),
		prettifyKeyVal(keyValToSlice(keysAndValues...)))
}

// keyValToSlice returns with a string slice by given keysAndValues argument.
// For example:
//  []string{"Debug", "information"}
func keyValToSlice(keysAndValues ...interface{}) (sl []string) {
	for _, v := range keysAndValues {
		sl = append(sl, fmt.Sprintf("%+v", v))
	}
	return sl
}
