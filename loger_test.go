package logger

import (
	"testing"
)

func TestLoger(t *testing.T) {
	log := NewFileLogger("debug", "./log")
	log.Info("info message")
	log.Error("error message")
}

func BenchmarkLoger(b *testing.B) {
	log := NewFileLogger("debug", "./log")
	for i := 0; i <= b.N; i++ {
		log.Info("info message")
		log.Error("error message")
	}
}
