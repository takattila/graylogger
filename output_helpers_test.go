package graylogger

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type outputHelpersSuite struct {
	suite.Suite
}

func (s outputHelpersSuite) TestSetLogLevelHandlers() {
	// All levels are enabled
	lh := setLogLevelHandlers(levelDebugNum, os.Stdout)
	s.NotEqual("0", fmt.Sprint(lh.debug))
	s.NotEqual("0", fmt.Sprint(lh.info))
	s.NotEqual("0", fmt.Sprint(lh.warn))
	s.NotEqual("0", fmt.Sprint(lh.error))
	s.NotEqual("0", fmt.Sprint(lh.fatal))

	// The debug level is disabled
	lh = setLogLevelHandlers(levelInfoNum, os.Stdout)
	s.Equal("0", fmt.Sprint(lh.debug))
	s.NotEqual("0", fmt.Sprint(lh.info))
	s.NotEqual("0", fmt.Sprint(lh.warn))
	s.NotEqual("0", fmt.Sprint(lh.error))
	s.NotEqual("0", fmt.Sprint(lh.fatal))

	// The debug, info levels are disabled
	lh = setLogLevelHandlers(levelWarningNum, os.Stdout)
	s.Equal("0", fmt.Sprint(lh.debug))
	s.Equal("0", fmt.Sprint(lh.info))
	s.NotEqual("0", fmt.Sprint(lh.warn))
	s.NotEqual("0", fmt.Sprint(lh.error))
	s.NotEqual("0", fmt.Sprint(lh.fatal))

	// The debug, info, warning levels are disabled
	lh = setLogLevelHandlers(levelErrorNum, os.Stdout)
	s.Equal("0", fmt.Sprint(lh.debug))
	s.Equal("0", fmt.Sprint(lh.info))
	s.Equal("0", fmt.Sprint(lh.warn))
	s.NotEqual("0", fmt.Sprint(lh.error))
	s.NotEqual("0", fmt.Sprint(lh.fatal))

	// All levels are disabled, except: fatal
	lh = setLogLevelHandlers(levelFatalNum, os.Stdout)
	s.Equal("0", fmt.Sprint(lh.debug))
	s.Equal("0", fmt.Sprint(lh.info))
	s.Equal("0", fmt.Sprint(lh.warn))
	s.Equal("0", fmt.Sprint(lh.error))
	s.NotEqual("0", fmt.Sprint(lh.fatal))
}

func (s outputHelpersSuite) TestLogLevelToInt() {
	num := logLevelToInt(LevelDebug)
	s.Equal(7, num)

	num = logLevelToInt(LevelInfo)
	s.Equal(6, num)

	num = logLevelToInt(LevelWarning)
	s.Equal(4, num)

	num = logLevelToInt(LevelError)
	s.Equal(3, num)

	num = logLevelToInt(LevelFatal)
	s.Equal(2, num)
}

func (s outputHelpersSuite) TestLogLevelToString() {
	str := logLevelToString(levelDebugNum)
	s.Equal("debug", str)

	str = logLevelToString(levelInfoNum)
	s.Equal("info", str)

	str = logLevelToString(levelWarningNum)
	s.Equal("warning", str)

	str = logLevelToString(levelErrorNum)
	s.Equal("error", str)

	str = logLevelToString(levelFatalNum)
	s.Equal("fatal", str)

	str = logLevelToString(0)
	s.Equal("debug", str)
}

func (s outputHelpersSuite) TestGetFuncName() {
	fn := getFuncName(0)
	s.Equal("graylogger.outputHelpersSuite.TestGetFuncName", fn)

	unknown := getFuncName(100)
	s.Equal("unknown", unknown)
}

func (s outputHelpersSuite) TestFetchNameFromPath() {
	path := "/path/fo/testFuncName"

	funcName := fetchNameFromPath(path)
	s.Equal("testFuncName", funcName)

	path = "testFuncName"

	funcName = fetchNameFromPath(path)
	s.Equal("testFuncName", funcName)
}

func (s outputHelpersSuite) TestColorOut() {
	testInit.LogColor = true

	colorText := testInit.colorOut(colorRed, colorRed)
	s.Equal(`"\x1b[0;91mred\x1b[0m"`, fmt.Sprintf("%q", colorText))

	colorText = testInit.colorOut(colorGreen, colorGreen)
	s.Equal(`"\x1b[0;32mgreen\x1b[0m"`, fmt.Sprintf("%q", colorText))

	colorText = testInit.colorOut(colorYellow, colorYellow)
	s.Equal(`"\x1b[0;93myellow\x1b[0m"`, fmt.Sprintf("%q", colorText))

	colorText = testInit.colorOut(colorBlue, colorBlue)
	s.Equal(`"\x1b[0;34mblue\x1b[0m"`, fmt.Sprintf("%q", colorText))

	colorText = testInit.colorOut(colorPurple, colorPurple)
	s.Equal(`"\x1b[0;95mpurple\x1b[0m"`, fmt.Sprintf("%q", colorText))

	colorText = testInit.colorOut(colorGray, colorGray)
	s.Equal(`"\x1b[0;90mgray\x1b[0m"`, fmt.Sprintf("%q", colorText))

	testInit.LogColor = false
}

func TestOutputHelpersSuite(t *testing.T) {
	suite.Run(t, new(outputHelpersSuite))
}
