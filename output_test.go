package graylogger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testInit = Init{
	//GraylogHost:     "127.0.0.1",
	//GraylogPort:     12201,
	//GraylogProvider: "TestService",
	//GraylogProtocol: TransportUDP,

	LogEnv:   "test",
	LogLevel: LevelDebug,
	LogColor: false,
}

var testOutputFileName = "test.out"

type outputSuite struct {
	suite.Suite
}

func (s outputSuite) TestNew() {
	g := New(testInit)

	s.Equal(testInit.GraylogHost, g.initData.GraylogHost)
	s.Equal(testInit.GraylogPort, g.initData.GraylogPort)
	s.Equal(testInit.GraylogProtocol, g.initData.GraylogProtocol)
	s.Equal(testInit.GraylogProvider, g.initData.GraylogProvider)
	s.Equal(testInit.LogLevel, g.initData.LogLevel)
	s.Equal(testInit.LogEnv, g.initData.LogEnv)
	s.Equal(testInit.LogColor, g.initData.LogColor)
	s.Equal(true, g.IsAllowedOutput())
}

func (s outputSuite) TestNewValidateLogLevelFatal() {
	init := testInit
	init.GraylogHost = "127.0.0.1"
	init.GraylogPort = 12201
	init.GraylogProvider = "TestService"
	init.GraylogProtocol = TransportUDP
	init.LogLevel = "bad_log_level"

	ExpectedPanicText := "Fatal function called"

	panicFunc := func(int) {
		panic(ExpectedPanicText)
	}

	patch := monkey.Patch(os.Exit, panicFunc)
	defer patch.Unpatch()

	assert.PanicsWithValue(
		s.T(),
		ExpectedPanicText,
		func() {
			_ = New(init)

		},
		"Fatal function was not called")
}

func (s outputSuite) TestNewValidateTransportFatal() {
	init := testInit
	init.GraylogHost = "127.0.0.1"
	init.GraylogPort = 12201
	init.GraylogProvider = "TestService"
	init.GraylogProtocol = TransportUDP
	init.GraylogProtocol = "bad_protocol"

	ExpectedPanicText := "Fatal function called"

	panicFunc := func(int) {
		panic(ExpectedPanicText)
	}

	patch := monkey.Patch(os.Exit, panicFunc)
	defer patch.Unpatch()

	assert.PanicsWithValue(
		s.T(),
		ExpectedPanicText,
		func() {
			_ = New(init)

		},
		"Fatal function was not called")
}

func (s outputSuite) TestTracking() {
	t := Tracking(1)
	s.Equal("output_test.go", t.File)
	s.NotEqual(0, t.Line)
	s.Equal("graylogger.outputSuite.TestTracking", t.Function)

	func() {
		t := Tracking(1)
		s.Equal("graylogger.outputSuite.TestTracking.func1", t.Function)
	}()

	func() {
		t := Tracking(2)
		s.Equal("graylogger.outputSuite.TestTracking", t.Function)
	}()
}

func (s outputSuite) TestDiscardOutput() {
	g := New(testInit)

	g.DiscardOutput()

	g.CaptureOutput(testOutputFileName)
	g.Debug("test", LevelDebug)
	g.Info("test", LevelInfo)
	g.Warning("test", LevelWarning)
	g.Error("test", LevelError)
	g.SaveOutput()

	s.Equal("", g.GetOutput())

	resetTest(s)
}

func (s outputSuite) TestResetLogger() {
	g := New(testInit)

	g.DiscardOutput()

	g.CaptureOutput(testOutputFileName)
	g.Debug("test", LevelDebug)
	g.Info("test", LevelInfo)
	g.Warning("test", LevelWarning)
	g.Error("test", LevelError)
	g.SaveOutput()

	s.Equal("", g.GetOutput())

	g = g.ResetLogger()

	g.CaptureOutput(testOutputFileName)
	g.Debug("test", LevelDebug)
	g.Info("test", LevelInfo)
	g.Warning("test", LevelWarning)
	g.Error("test", LevelError)
	g.SaveOutput()

	s.Equal(true, strings.Contains(g.GetOutput(), "test :: debug"))
	s.Equal(true, strings.Contains(g.GetOutput(), "test :: info"))
	s.Equal(true, strings.Contains(g.GetOutput(), "test :: warning"))
	s.Equal(true, strings.Contains(g.GetOutput(), "test :: error"))

	resetTest(s)
}

func (s outputSuite) TestDebug() {
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	g.Debug("test", LevelDebug)
	g.SaveOutput()

	s.Equal(true, strings.Contains(g.GetOutput(), "test :: debug"))

	resetTest(s)
}

func (s outputSuite) TestInfo() {
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	g.Info("test", LevelInfo)
	g.SaveOutput()

	s.Equal(true, strings.Contains(g.GetOutput(), "test :: info"))

	resetTest(s)
}

func (s outputSuite) TestWarning() {
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	g.Warning("test", LevelWarning)
	g.SaveOutput()

	s.Equal(true, strings.Contains(g.GetOutput(), "test :: warning"))

	resetTest(s)
}

func (s outputSuite) TestLogWarningIfErr() {
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	err := fmt.Errorf("example error")
	g.LogWarningIfErr(err)
	g.SaveOutput()

	expected := "graylogger.outputSuite.TestLogWarningIfErr :: example error"
	s.Equal(true, strings.Contains(g.GetOutput(), expected))

	resetTest(s)
}

func (s outputSuite) TestError() {
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	g.Error("test", LevelError)
	g.SaveOutput()

	s.Equal(true, strings.Contains(g.GetOutput(), "test :: error"))

	resetTest(s)
}

func (s outputSuite) TestLogErrorIfErr() {
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	err := fmt.Errorf("example error")
	g.LogErrorIfErr(err)
	g.SaveOutput()

	expected := "graylogger.outputSuite.TestLogErrorIfErr :: example error"
	s.Equal(true, strings.Contains(g.GetOutput(), expected))

	resetTest(s)
}

func (s outputSuite) TestFatal() {
	init := testInit
	ExpectedPanicText := "Fatal function called"

	panicFunc := func(int) {
		panic(ExpectedPanicText)
	}

	patch := monkey.Patch(os.Exit, panicFunc)
	defer patch.Unpatch()

	assert.PanicsWithValue(
		s.T(),
		ExpectedPanicText,
		func() {
			g := New(init)
			g.CaptureOutput(testOutputFileName)

			err := fmt.Errorf("example fatal error")
			g.Fatal(err)

			g.SaveOutput()
			s.Equal(true, strings.Contains(g.GetOutput(), "example fatal error"))
		},
		"Fatal function was not called")

	resetTest(s)
}

func (s outputSuite) TestReturnWithError() {
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	err := g.ReturnWithError("example", "error")
	s.Equal("example :: error", fmt.Sprint(err))
	g.SaveOutput()

	expected := "example :: error"
	s.Equal(true, strings.Contains(g.GetOutput(), expected))

	resetTest(s)
}

func (s outputSuite) TestGetInit() {
	g := New(testInit)
	i := g.GetInit()

	s.Equal(LevelDebug, i.LogLevel)
	s.Equal("test", i.LogEnv)
	s.Equal(false, i.LogColor)
}

func (s outputSuite) TestIsAllowedOutput() {
	g := New(testInit)
	s.Equal(true, g.IsAllowedOutput())

	g.DiscardOutput()
	s.Equal(false, g.IsAllowedOutput())
}

func (s outputSuite) TestGetLogLevel() {
	g := New(testInit)
	LevelNum, levelString := g.GetLogLevel()

	s.Equal(7, LevelNum)
	s.Equal("debug", levelString)

	g.DiscardOutput()
	s.Equal(false, g.IsAllowedOutput())
}

func (s outputSuite) TestPrintOutput() {
	// 1. Save output
	g := New(testInit)

	g.CaptureOutput(testOutputFileName)
	g.Error("test", LevelError)
	g.SaveOutput()

	s.Equal(true, strings.Contains(g.GetOutput(), "test :: error"))

	// 2. Print saved output from file and capture printed text from stdOut
	testPrintOutputFileName := "capture_print_output.out"
	captureOutput(testPrintOutputFileName, func() {
		g.PrintOutput()
	})

	// 3. Captured output should match
	b, err := ioutil.ReadFile(testPrintOutputFileName)
	s.Equal(nil, err)

	s.Equal(true, strings.Contains(string(b), "test :: error"))

	err = os.Remove(testPrintOutputFileName)
	s.Equal(nil, err)

	resetTest(s)
}

func resetTest(s outputSuite) {
	err := os.Remove(testOutputFileName)
	s.Equal(nil, err)
}

func captureOutput(fileName string, f func()) {
	reader, writer, _ := os.Pipe()

	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()

	f()

	_ = writer.Close()
	_ = ioutil.WriteFile(fileName, []byte(<-out), os.ModePerm)
}

func TestOutputSuite(t *testing.T) {
	suite.Run(t, new(outputSuite))
}
