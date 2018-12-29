package logpeck

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
)

// LogTask .
type LogTask struct {
	LogPath string

	peckTasks map[string]*PeckTask
	tail      *tail.Tail
	stop      bool
	errMsg    string
}

// NewLogTask .
func NewLogTask(path string) *LogTask {
	task := &LogTask{
		LogPath:   path,
		peckTasks: make(map[string]*PeckTask),
		tail:      nil,
		stop:      true,
	}
	return task
}

// AddPeckTask .
func (p *LogTask) AddPeckTask(task *PeckTask) error {
	p.peckTasks[task.Config.Name] = task
	return nil
}

// UpdatePeckTask .
func (p *LogTask) UpdatePeckTask(task *PeckTask) error {
	if !task.IsStop() {
		if err := p.peckTasks[task.Config.Name].Stop(); err != nil {
			return err
		}
		p.peckTasks[task.Config.Name] = task
		if err := task.Start(); err != nil {
			return err
		}
	} else {
		p.peckTasks[task.Config.Name] = task
	}
	return nil
}

// RemovePeckTask .
func (p *LogTask) RemovePeckTask(config *PeckTaskConfig) error {
	if !p.peckTasks[config.Name].IsStop() {
		p.peckTasks[config.Name].Stop()
	}
	delete(p.peckTasks, config.Name)
	return nil
}

// StartPeckTask .
func (p *LogTask) StartPeckTask(config *PeckTaskConfig) error {
	if !p.Exist(config) {
		panic(config)
	}
	if p.peckTasks[config.Name].IsStop() {
		if err := p.peckTasks[config.Name].Start(); err != nil {
			return err
		}
	} else {
		panic(config)
	}
	return nil
}

// StopPeckTask .
func (p *LogTask) StopPeckTask(config *PeckTaskConfig) error {
	if !p.Exist(config) {
		panic(config)
	}
	if !p.peckTasks[config.Name].IsStop() {
		if err := p.peckTasks[config.Name].Stop(); err != nil {
			return err
		}
	} else {
		panic(config)
	}
	return nil
}

// Exist .
func (p *LogTask) Exist(config *PeckTaskConfig) bool {
	_, ok := p.peckTasks[config.Name]
	return ok
}

// Empty .
func (p *LogTask) Empty() bool {
	if len(p.peckTasks) == 0 {
		return true
	}
	return false
}

func peckLogBG(p *LogTask) {
	log.Infof("[LogTask %s] Start peck log", p.LogPath)
	for content := range p.tail.Lines {
		for name, task := range p.peckTasks {
			// process log
			log.Debugf("[LogTask %s] %s content[%s]", p.LogPath, name, content.Text)
			task.Process(content.Text)
		}
		if p.stop {
			break
		}
	}
}

// Start .
func (p *LogTask) Start() error {
	if !p.stop {
		return errors.New("LogTask already started")
	}
	log.Infof("[LogTask %s] Start LogTask", p.LogPath)
	if p.tail == nil {
		tailConf := tail.Config{
			ReOpen: true,
			Poll:   true,
			Follow: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: 2,
			},
		}
		p.tail, _ = tail.TailFile(p.LogPath, tailConf)
	}

	go peckLogBG(p)
	p.stop = false
	return nil
}

// Stop .
func (p *LogTask) Stop() error {
	if p.stop {
		return errors.New("LogTask already stopped")
	}
	log.Infof(" [LogTask %s] Stop LogTask", p.LogPath)
	p.stop = true
	p.tail.Stop()
	p.tail = nil
	return nil
}

// IsStop .
func (p *LogTask) IsStop() bool {
	return p.stop
}

// Close .
func (p *LogTask) Close() error {
	// NOT IMPLEMENT
	return nil
}

// GetStat .
func (p *LogTask) GetStat() *LogStat {
	// NOT IMPLEMENT
	return nil
}
