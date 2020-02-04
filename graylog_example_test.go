package graylogger_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/takattila/graylogger"
)

func ExampleGrayLogger_SendGELF() {
	//  UDP server initialization ...
	udpServer := func(port int, response chan string) {
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("0.0.0.0"),
		})
		if err != nil {
			log.Fatalf("net.ListenUDP: %s\n", err)
		}

		defer func() {
			_ = conn.Close()
		}()

		message := make([]byte, 262144)
		_, _, err = conn.ReadFromUDP(message[:])
		if err != nil {
			log.Fatalf("conn.ReadFromUDP: %s\n", err)
		}

		response <- string(message)
	}

	// Setup Init structure ...
	init := graylogger.Init{
		GraylogHost:     "127.0.0.1",
		GraylogPort:     12201,
		GraylogProvider: "TestService",
		GraylogProtocol: graylogger.TransportUDP,

		LogEnv:   "test",
		LogLevel: graylogger.LevelDebug,
		LogColor: false,
	}

	// Starting UDP server ...
	response := make(chan string)
	go udpServer(init.GraylogPort, response)
	time.Sleep(100 * time.Millisecond)

	// GrayLogger initialization ...
	g := graylogger.New(init)

	g.CaptureOutput("test.out")
	g.Info("test", "gelf_message")
	g.SaveOutput()

	// Removing NUL characters from bytes
	responseBytes := bytes.Trim([]byte(<-response), "\x00")

	// Parsing JSON response ...
	obj := map[string]interface{}{}
	err := json.Unmarshal(responseBytes, &obj)
	if err != nil {
		log.Fatalf("json.Unmarshal: %v\n", err)
	}

	fmt.Println(obj["_log_env"])
	fmt.Println(obj["_log_key"])
	fmt.Println(obj["_log_level"])
	fmt.Println(obj["_log_value"])
	fmt.Println(obj["_track_file"])
	fmt.Println(obj["_track_function"])
	fmt.Println(obj["full_message"])
	fmt.Println(obj["short_message"])
	fmt.Println(obj["host"])
	fmt.Println(obj["version"])
	fmt.Println(obj["level"])

	// Output:
	// test
	// test
	// info
	// gelf_message
	// graylog_example_test.go
	// graylogger_test.ExampleGrayLogger_SendGELF
	// test :: gelf_message
	// test :: gelf_message
	// TestService
	// 1.1
	// 6
}
