package logpeck

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// MockLog .
type MockLog struct {
	Path      string
	IsRunning bool

	stop bool
	file *os.File
	mu   sync.Mutex
}

// NewMockLog .
func NewMockLog(path string) (*MockLog, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &MockLog{Path: path, IsRunning: false, file: f, stop: false}, nil
}

func genLog() string {
	now := time.Now().String()
	randNum := rand.Intn(65536)
	return fmt.Sprintf("%s mocklog %d .\n", now, randNum)
}

// Run .
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
		time.Sleep(1027 * time.Millisecond)
		p.mu.Lock()
	}
	p.IsRunning = false
	p.stop = false
	return nil
}

// Stop .
func (p *MockLog) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stop = true
}

// Close .
func (p *MockLog) Close() {
	p.Stop()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.file.Close()
}
