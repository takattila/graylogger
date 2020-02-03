package graylogger_test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/takattila/graylogger"
)

func ExampleNew() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: true,
	})

	fmt.Println(g.GetInit())

	// Output: { 0   100ms test debug true}
}

func ExampleTracking() {
	track := graylogger.Tracking(1)

	fmt.Println(track.File)

	num, _ := strconv.Atoi(track.Line)
	fmt.Printf("%T\n", num)

	fmt.Println(track.Function)

	// Output:
	// output_example_test.go
	// int
	// graylogger_test.ExampleTracking
}

func ExampleGrayLogger_DiscardOutput() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: true,
	})

	g.DiscardOutput()

	fmt.Println(g.IsAllowedOutput())

	// Output: false
}

func ExampleGrayLogger_ResetLogger() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: true,
	})

	g.DiscardOutput()
	fmt.Println(g.IsAllowedOutput())

	g = g.ResetLogger()
	fmt.Println(g.IsAllowedOutput())

	// Output:
	// false
	// true
}

func ExampleGrayLogger_Debug() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")
	g.Debug("test", graylogger.LevelDebug)
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[DEBUG]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_Debug"))
	fmt.Println(strings.Contains(output, "[test :: debug]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_Info() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")
	g.Info("test", graylogger.LevelInfo)
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[INFO]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_Info"))
	fmt.Println(strings.Contains(output, "[test :: info]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_Warning() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")
	g.Warning("test", graylogger.LevelWarning)
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[WARNING]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_Warning"))
	fmt.Println(strings.Contains(output, "[test :: warning]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_LogWarningIfErr() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")

	err := fmt.Errorf("test :: %s", graylogger.LevelWarning)
	g.LogWarningIfErr(err)
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[WARNING]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_LogWarningIfErr"))
	fmt.Println(strings.Contains(output, "[graylogger_test.ExampleGrayLogger_LogWarningIfErr :: test :: warning]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_Error() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")
	g.Error("test", graylogger.LevelError)
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[ERROR]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_Error"))
	fmt.Println(strings.Contains(output, "[test :: error]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_LogErrorIfErr() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")

	err := fmt.Errorf("test :: %s", graylogger.LevelError)
	g.LogErrorIfErr(err)
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[ERROR]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_LogErrorIfErr"))
	fmt.Println(strings.Contains(output, "[graylogger_test.ExampleGrayLogger_LogErrorIfErr :: test :: error]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_ReturnWithError() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")

	ret := g.ReturnWithError("test", graylogger.LevelError)
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[ERROR]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_ReturnWithError"))
	fmt.Println(ret)

	// Output:
	// true
	// true
	// true
	// test :: error
}

func ExampleGrayLogger_Fatal() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	err := fmt.Errorf("test :: %s", graylogger.LevelError)
	g.Fatal(err)

	// [FATAL] 2020/02/03 14:43:37 [file: example.go line: 17 function: main.main] [main.main :: test :: error]
	// exit status 1
}

func ExampleGrayLogger_GetInit() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	init := g.GetInit()

	fmt.Println(init.LogEnv)
	fmt.Println(init.LogLevel)
	fmt.Println(init.LogColor)

	// Output:
	// test
	// debug
	// false
}

func ExampleGrayLogger_IsAllowedOutput() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: true,
	})

	fmt.Println(g.IsAllowedOutput())

	// Output: true
}

func ExampleGrayLogger_GetLogLevel() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: true,
	})

	levelNum, levelString := g.GetLogLevel()

	fmt.Println(levelNum)
	fmt.Println(levelString)

	// Output:
	// 7
	// debug
}

func ExampleGrayLogger_SaveOutput() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")
	g.Debug("test", "save_output")
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[DEBUG]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_SaveOutput"))
	fmt.Println(strings.Contains(output, "[test :: save_output]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_GetOutput() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")
	g.Debug("test", "get_output")
	g.SaveOutput()

	output := g.GetOutput()
	fmt.Println(strings.Contains(output, "[DEBUG]"))
	fmt.Println(strings.Contains(output, "file: output_example_test.go"))
	fmt.Println(strings.Contains(output, "function: graylogger_test.ExampleGrayLogger_GetOutput"))
	fmt.Println(strings.Contains(output, "[test :: get_output]"))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleGrayLogger_PrintOutput() {
	g := graylogger.New(graylogger.Init{
		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	})

	g.CaptureOutput("test.out")
	g.Debug("test", "print_output")
	g.SaveOutput()

	g.PrintOutput()

	// [DEBUG] 2020/02/03 14:44:37 [file: example.go line: 15 function: main.main] [test :: print_output]
}
