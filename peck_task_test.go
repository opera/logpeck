package logpeck

import (
	"log"
	"testing"
	"time"
)

func TestPeckTask(*testing.T) {
	log_name := ".test.log"

	// Mock a user log
	mock_log, m_err := NewMockLog(log_name)
	if m_err != nil {
		panic(m_err)
	}
	defer mock_log.Close()
	go mock_log.Run()

	time.Sleep(1001 * time.Millisecond)
	log.Print(mock_log)
}
