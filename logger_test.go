package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	log := NewFileLogger("debug", "./log", true)
	log.Info("info message %s", "i am hello world")
	log.Debug("debug message")
	log.Warning("warning message")
	log.Error("error message")
}

func BenchmarkLogger(b *testing.B) {
	log := NewFileLogger("debug", "./log", true)
	for i := 0; i <= b.N; i++ {
		log.Info("info message")
	}
}
