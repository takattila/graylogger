package graylogger

// SendGELF sends GELF messages into Graylog instance.
// If the Graylog host is unreachable, it writes an error message to stdOut.
func (g *GrayLogger) SendGELF(level int, keysAndValues ...interface{}) {
	if g.validateGraylogArguments(level) {
		g.connect().send(level, keysAndValues)
	}
}
