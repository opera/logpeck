package logpeck

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
)

func TestTailLog(*testing.T) {
	defer LogExecTime(time.Now(), "TestTailLog")
	logName := ".test.log"

	// Mock a user log
	mockLog, err := NewMockLog(logName)
	if err != nil {
		panic(err)
	}
	defer mockLog.Close()
	go func() {
		time.Sleep(200 * time.Millisecond)
		mockLog.Run()
	}()

	conf := tail.Config{ReOpen: true, Poll: true, Follow: true}
	t, _ := tail.TailFile(logName, conf)
	cnt := 0
	for line := range t.Lines {
		log.Infof("[" + line.Text + "]")
		if cnt > 5 {
			break
		}
		cnt++
		time.Sleep(100 * time.Millisecond)
	}
}
