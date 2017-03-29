package logpeck

import (
	"bufio"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

type LogTask struct {
	LogPath string

	peckTasks map[string]*PeckTask
	file      *os.File
	mu        sync.Mutex
	stop      bool
	errMsg    string
}

func NewLogTask(path string) *LogTask {
	task := &LogTask{
		LogPath:   path,
		peckTasks: make(map[string]*PeckTask),
		file:      nil,
		stop:      true,
	}
	return task
}

func (p *LogTask) AddPeckTask(config *PeckTaskConfig) error {
	p.peckTasks[config.Name] = NewPeckTask(config)
	return nil
}

func (p *LogTask) UpdatePeckTask(conf *PeckTaskConfig) error {
	// NOT IMPLEMENT
	return nil
}

func (p *LogTask) RemovePeckTask(conf *PeckTaskConfig) error {
	// NOT IMPLEMENT
	return nil
}

func (p *LogTask) Exist(config *PeckTaskConfig) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.peckTasks[config.Name]
	return ok
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

func tailLog(f *os.File) string {
	return "hello logpeck"
}

func peckLogBG(p *LogTask) {
	log.Printf("[LogTask] Start peck log %s", p.LogPath)
	scanner := bufio.NewScanner(p.file)
	for scanner.Scan() {
		content := scanner.Text()
		{
			p.mu.Lock()
			defer p.mu.Unlock()
			for i, task := range p.peckTasks {
				// process log
				log.Printf("[LogTask] %d task[%s], content[%s]", i, task, content)
				task.Process(content)
			}
			if p.stop {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (p *LogTask) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.stop {
		return errors.New("LogTask already started")
	}
	log.Printf("[LogTask] Start LogTask on %s", p.LogPath)
	if p.file == nil {
		f, f_err := os.Open(p.LogPath)
		if f_err != nil {
			p.errMsg = f_err.Error()
			log.Printf("[LogTask] Log open failed, %s", f_err)
		} else {
			p.file = f
			p.file.Seek(0, 2)
		}
	} else {
		p.file.Seek(0, 2)
	}

	go peckLogBG(p)
	p.stop = false
	return nil
}

func (p *LogTask) Pause() error {
	// NOT IMPLEMENT
	return nil
}

func (p *LogTask) Close() error {
	// NOT IMPLEMENT
	return nil
}

func (p *LogTask) GetStat() *LogStat {
	// NOT IMPLEMENT
	return nil
}
