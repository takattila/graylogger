package graylogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type graylogSuite struct {
	suite.Suite
}

func (s graylogSuite) TestSendGELF() {
	init := testInit
	init.GraylogHost = "127.0.0.1"
	init.GraylogPort = 12201
	init.GraylogProvider = "TestService"
	init.GraylogProtocol = TransportUDP

	var response string
	go func() {
		resp, err := udpServer(init.GraylogPort)
		s.Equal(nil, err)
		response = resp
	}()

	time.Sleep(10 * time.Millisecond)

	g := New(init)
	g.Info("test", "info")

	for {
		if response != "" {
			obj := map[string]interface{}{}

			// Removing NUL characters from bytes
			responseBytes := bytes.Trim([]byte(response), "\x00")

			err := json.Unmarshal(responseBytes, &obj)
			s.Equal(nil, err)

			s.Equal("test", obj["_log_env"])
			s.Equal("test", obj["_log_key"])
			s.Equal("info", obj["_log_level"])
			s.Equal("info", obj["_log_value"])
			s.Equal("graylog_test.go", obj["_track_file"])
			s.Equal("graylogger.graylogSuite.TestSendGELF", obj["_track_function"])
			s.NotEqual("", obj["_track_line"])
			s.NotEqual("", obj["timestamp"])
			s.Equal("test :: info", obj["full_message"])
			s.Equal("test :: info", obj["short_message"])
			s.Equal("TestService", obj["host"])
			s.Equal("1.1", obj["version"])
			s.Equal(float64(6), obj["level"])

			break
		}
	}

}

func udpServer(port int) (string, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		return "", fmt.Errorf("net.ListenUDP: %s", err)
	}

	defer func() {
		_ = conn.Close()
	}()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	message := make([]byte, 262144)
	_, _, err = conn.ReadFromUDP(message[:])
	if err != nil {
		return "", fmt.Errorf("conn.ReadFromUDP: %s", err)
	}

	return string(message), nil
}

func TestGraylogSuite(t *testing.T) {
	suite.Run(t, new(graylogSuite))
}
