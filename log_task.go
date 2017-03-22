package logpeck

import (
	"os"
	"sync"
)

type LogTask struct {
	LogPath   string
	PeckTasks map[string]PeckTask
	IsRunning bool

	file *os.File
	stop bool
	mu   sync.Mutex
}

func NewLogTask(path string) (*LogTask, error) {
	f, f_err := os.Open(path)
	if f_err != nil {
		return nil, f_err
	}
	task := &LogTask{
		LogPath:   path,
		PeckTasks: make(map[string]PeckTask),
		file:      f,
	}
	return task, nil
}

func (p *LogTask) Run() {

}

func (p *LogTask) Pause() error {
	return nil
}

func (p *LogTask) Cancel() error {
	return nil
}
