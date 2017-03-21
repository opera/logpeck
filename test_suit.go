package logpeck

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

const kTestDBPath string = ".unittest.db"

func logExecTime(start time.Time, prefix string) {
	elapsed_ms := time.Since(start) / time.Millisecond
	fmt.Printf("Performance: %s cost %d ms.\n", prefix, elapsed_ms)
}

type MockLog struct {
	Path      string
	IsRunning bool
	stop      bool
	file      *os.File
	mu        sync.Mutex
}

func NewMockLog(path string) (*MockLog, error) {
	f, f_err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if f_err != nil {
		return nil, f_err
	}
	return &MockLog{Path: path, IsRunning: false, file: f, stop: false}, nil
}

func genLog() string {
	now := time.Now().String()
	rand_num := rand.Int()
	return fmt.Sprintf("%s mocklog %d .\n", now, rand_num)
}

func (p *MockLog) Run() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.IsRunning {
		return fmt.Errorf("log[%s] already running", p.Path)
	}

	p.IsRunning = true
	for !p.stop {
		p.file.WriteString(genLog())
		p.mu.Unlock()
		time.Sleep(9 * time.Millisecond)
		p.mu.Lock()
	}
	p.IsRunning = false
	p.stop = false
	return nil
}

func (p *MockLog) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stop = true
}

func (p *MockLog) Close() {
	p.Stop()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.file.Close()
}
