package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	//AddOption(WithPath("./log1"))
	AddOption(WithLevelStr("debug"), WithPath("./log1"), WithAsync(true))
	Info("info message %s %s", "i am hello world", "aa")
	AddOption(WithLevelStr("debug"), WithPath("./log2"))
	Info("info message %s %s", "i am hello world", "aa")
	//log.Debug("debug message")
	//log.Warning("warning message")
	//log.Error("error message")
}

func BenchmarkLogger(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		Info("info message")
	}
}
