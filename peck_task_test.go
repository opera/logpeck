package logpeck

import (
	"fmt"
	"testing"
	"time"
)

func TestPeckTask(*testing.T) {
	logname := ".test.log"
	mock_log, m_err := NewMockLog(logname)
	if m_err != nil {
		panic(m_err)
	}
	go mock_log.Run()
	time.Sleep(1001 * time.Millisecond)
	fmt.Print(mock_log)
	mock_log.Close()
}
