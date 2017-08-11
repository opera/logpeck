package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
	"time"
)

type LogTask struct {
	LogPath string

	peckTasks map[string]*PeckTask
	tail      *tail.Tail
	stop      bool
	errMsg    string
}

func NewLogTask(path string) *LogTask {
	task := &LogTask{
		LogPath:   path,
		peckTasks: make(map[string]*PeckTask),
		tail:      nil,
		stop:      true,
	}
	return task
}

func (p *LogTask) AddPeckTask(task *PeckTask) error {
	p.peckTasks[task.Config.Name] = task
	return nil
}

func (p *LogTask) UpdatePeckTask(task *PeckTask) error {
	task.Stat = p.peckTasks[task.Config.Name].Stat
	p.peckTasks[task.Config.Name] = task
	return nil
}

func (p *LogTask) RemovePeckTask(config *PeckTaskConfig) error {
	delete(p.peckTasks, config.Name)
	return nil
}

func (p *LogTask) StartPeckTask(config *PeckTaskConfig) error {
	if !p.Exist(config) {
		panic(config)
	}
	if p.peckTasks[config.Name].IsStop() {
		p.peckTasks[config.Name].Start()
	} else {
		panic(config)
	}
	return nil
}

func (p *LogTask) StopPeckTask(config *PeckTaskConfig) error {
	if !p.Exist(config) {
		panic(config)
	}
	if !p.peckTasks[config.Name].IsStop() {
		p.peckTasks[config.Name].Stop()
	} else {
		panic(config)
	}
	return nil
}

func (p *LogTask) Exist(config *PeckTaskConfig) bool {
	_, ok := p.peckTasks[config.Name]
	return ok
}

func (p *LogTask) Empty() bool {
	if len(p.peckTasks) == 0 {
		return true
	} else {
		return false
	}
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
		time.Sleep(10 * time.Millisecond)
	}
}

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

func (p *LogTask) IsStop() bool {
	return p.stop
}

func (p *LogTask) Close() error {
	// NOT IMPLEMENT
	return nil
}

func (p *LogTask) GetStat() *LogStat {
	// NOT IMPLEMENT
	return nil
}
