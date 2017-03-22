package logpeck

import (
	"os"
	"sync"
)

type LogTask struct {
	LogPath   string
	IsRunning bool

	peckTasks map[string]PeckTask
	file      *os.File
	stop      bool
	mu        sync.Mutex
}

func NewLogTask(path string) (*LogTask, error) {
	f, f_err := os.Open(path)
	if f_err != nil {
		return nil, f_err
	}
	task := &LogTask{
		LogPath:   path,
		peckTasks: make(map[string]PeckTask),
		file:      f,
	}
	return task, nil
}

func (p *LogTask) AddPeckTask(conf *PeckTaskConfig) error {
	return nil
}

func (p *LogTask) RemovePeckTask(conf *PeckTaskConfig) error {
	return nil
}

func (p *LogTask) Empty() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.peckTasks) == 0 {
		return true
	} else {
		return false
	}
}

func (p *LogTask) Run() {

}

func (p *LogTask) Pause() error {
	return nil
}

func (p *LogTask) Close() error {
	return nil
}

func (p *LogTask) GetStat() *LogStat {
	return nil
}
