package logpeck

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
	"testing"
	"time"
)

func TestTailLog(*testing.T) {
	defer LogExecTime(time.Now(), "TestTailLog")
	logName := ".test.log"

	// Mock a user log
	mock_log, m_err := NewMockLog(logName)
	if m_err != nil {
		panic(m_err)
	}
	defer mock_log.Close()
	go func() {
		time.Sleep(200 * time.Millisecond)
		mock_log.Run()
	}()

	conf := tail.Config{ReOpen: true, Poll: true, Follow: true}
	t, _ := tail.TailFile(logName, conf)
	cnt := 0
	for line := range t.Lines {
		log.Infof("[" + line.Text + "]")
		if cnt > 5 {
			break
		}
		cnt += 1
		time.Sleep(100 * time.Millisecond)
	}
}
