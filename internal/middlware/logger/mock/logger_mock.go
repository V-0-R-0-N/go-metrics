package mock

import (
	"bytes"
	"fmt"
)

type LoggerMock struct {
	buf bytes.Buffer
}

func (m *LoggerMock) Infoln(args ...interface{}) {
	for _, arg := range args {
		s := fmt.Sprintf("%v\n", arg)
		m.buf.Write([]byte(s))
	}
}

func (m *LoggerMock) Get() string {
	return m.buf.String()
}

func (m *LoggerMock) reset() {
	m.buf.Reset()
}

func (m *LoggerMock) Infow(msg string, keysAndValues ...interface{}) {

}
