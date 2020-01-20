package graylogger

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/Devatoria/go-graylog"
	"github.com/tidwall/pretty"
)

// graylogTimeout declares the maximum amount of time
// a dial will wait for a connection to complete.
const graylogTimeout = 100 * time.Millisecond

// GraylogExtraFields represents the extra GELF data which will be sent into Graylog instance.
type GraylogExtraFields struct {
	Env      string `json:"log_env"`
	Level    string `json:"log_level"`
	Key      string `json:"log_key"`
	Value    string `json:"log_value"`
	Line     string `json:"track_line"`
	File     string `json:"track_file"`
	Function string `json:"track_function"`
}

// validateGraylogArguments checks that all obligatory parameters set,
// that are needed to send log messages to Graylog instance.
// The following values are must be set to send GELF messages to Graylog:
//   - Level (int) syslog level
//   - GraylogHost (string) the domain name of the Graylog instance
//   - GraylogPort (string) the port number of the Graylog instance
//   - GraylogProvider (string) the name of the service which sends the messages
//   - GraylogProtocol (Transport) TCP or UDP
func (g *GrayLogger) validateGraylogArguments(level int) bool {
	return g.isSetGraylogObligatoryFields() && g.checkHostIsAlive() && g.IsAllowedOutput()
}

// isSetGraylogObligatoryFields checks that all Graylog related obligatory fields are set .
func (g *GrayLogger) isSetGraylogObligatoryFields() bool {
	return g.level != 0 &&
		g.initData.GraylogHost != "" &&
		g.initData.GraylogPort != 0 &&
		g.initData.GraylogProvider != "" &&
		g.initData.GraylogProtocol != ""
}

// connect instantiates a new graylog connection using the given endpoint.
func (g *GrayLogger) connect() *GrayLogger {
	g.graylog, _ = graylog.NewGraylog(graylog.Endpoint{
		Transport: graylog.Transport(g.initData.GraylogProtocol),
		Address:   g.initData.GraylogHost,
		Port:      uint(g.initData.GraylogPort),
	})
	return g
}

// send iterates over key : value pairs
// and send them to Graylog instance one by one as a GELF message.
func (g *GrayLogger) send(level int, keysAndValues []interface{}) {
	tr := getTrackingInfo(3)
	for key, val := range keysAndValuesToMap(keysAndValues) {
		if g.graylog != nil && g.level >= level {
			_ = g.graylog.Send(graylog.Message{
				Version:      "1.1",
				Host:         g.initData.GraylogProvider,
				ShortMessage: prettifyKeyVal(keyValToSlice(key, cleanString(fmt.Sprint(val)))),
				FullMessage:  prettifyKeyVal(keyValToSlice(key, val)),
				Timestamp:    time.Now().Unix(),
				Level:        uint(level),
				Extra: createExtraFieldsMap(GraylogExtraFields{
					Env:      g.initData.LogEnv,
					Level:    logLevelToString(level),
					Key:      prettifyObject(key),
					Value:    prettifyObject(val),
					Line:     tr.Line,
					File:     tr.File,
					Function: tr.Function,
				}),
			})
			_ = g.graylog.Close()
		}
	}
}

// checkHostIsAlive validates Graylog host connection.
func (g *GrayLogger) checkHostIsAlive() bool {
	_, err := net.DialTimeout(
		fmt.Sprint(g.initData.GraylogProtocol),
		g.initData.GraylogHost+":"+fmt.Sprint(g.initData.GraylogPort),
		g.initData.GraylogTimeout)
	if err != nil {
		couldNotConnect := "could not connect to Graylog host with initialized data"
		errorMsg := []interface{}{couldNotConnect, g.initData}
		g.functions.Error.Println(g.formatLogLine(errorMsg...))
	}
	return err == nil
}

// cleanString removes whitespaces and newlines from a string.
func cleanString(text string) string {
	if json.Valid([]byte(text)) {
		return string(pretty.Ugly([]byte(text)))
	}
	return strings.Join(strings.Fields(strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))), " ")
}

// keysAndValuesToMap creates multiple key : value pairs
// if more than one are passed as an argument into logger function.
// After that, we can send these pairs into Graylog instance one by one.
func keysAndValuesToMap(keysAndValues []interface{}) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for i, v := range keysAndValues {
		if i%2 == 0 {
			m[v] = keysAndValues[i+1]
		}
	}
	return m
}

// prettifyKeyVal converts a string slice to a formatted string.
// For example:
//  Debug :: information
func prettifyKeyVal(keysAndValues []string) string {
	return strings.Join(keysAndValues, " :: ")
}

// createExtraFieldsMap makes a map from GraylogExtraFields struct,
// to send them into Graylog.
func createExtraFieldsMap(extra GraylogExtraFields) map[string]string {
	ret := make(map[string]string)
	v := reflect.ValueOf(extra)
	for i := 0; i < v.NumField(); i++ {
		ret[v.Type().Field(i).Tag.Get("json")] = v.Field(i).String()
	}
	return ret
}

// prettifyObject create a JSON string if it possible
// to make any structures human-readable in Graylog.
func prettifyObject(obj interface{}) string {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Struct:
		return makePrettifiedJSON(obj)
	case reflect.Slice:
		return makePrettifiedJSON(obj)
	case reflect.Map:
		return makePrettifiedJSON(obj)
	case reflect.String:
		return makePrettifiedJsonIfPossible(obj)
	default:
		return makePrettifiedJsonIfPossible(obj)
	}
}

// makePrettifiedJSON creates an indented JSON string.
func makePrettifiedJSON(obj interface{}) string {
	js, _ := json.MarshalIndent(obj, "", "  ")
	return string(js)
}

// makePrettifiedJsonIfPossible creates an indented JSON string if obj argument is a valid JSON.
func makePrettifiedJsonIfPossible(obj interface{}) string {
	e := fmt.Sprint(obj)
	if json.Valid([]byte(e)) {
		js := pretty.Pretty([]byte(e))
		return string(js)
	}
	return e
}
