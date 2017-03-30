package logpeck

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

type LogTask struct {
	LogPath string

	peckTasks map[string]*PeckTask
	file      *os.File
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

func (p *LogTask) AddPeckTask(task *PeckTask) error {
	p.peckTasks[task.Config.Name] = task
	return nil
}

func (p *LogTask) UpdatePeckTask(conf *PeckTaskConfig) error {
	// NOT IMPLEMENT
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

func tailLog(f *os.File) string {
	return "hello logpeck"
}

func peckLogBG(p *LogTask) {
	log.Printf("[LogTask %s] Start peck log", p.LogPath)
	reader := bufio.NewReader(p.file)
	for {
		line, _, err := reader.ReadLine()
		content := string(line[:])
		if err == io.EOF {
			log.Printf("[LogTask %s] Log finished", p.LogPath)
		}
		log.Printf("[LogTask %s] Log found, content[%s]", p.LogPath, content)
		for name, task := range p.peckTasks {
			// process log
			log.Printf("[LogTask %s] %s content[%s]", p.LogPath, name, content)
			task.Process(content)
		}
		if p.stop {
			break
		}
		log.Printf("[LogTask %s] Sleep for a while", p.LogPath)
		time.Sleep(10 * time.Millisecond)
	}
}

func (p *LogTask) Start() error {
	if !p.stop {
		return errors.New("LogTask already started")
	}
	log.Printf(" [LogTask %s] Start LogTask", p.LogPath)
	if p.file == nil {
		f, f_err := os.Open(p.LogPath)
		if f_err != nil {
			p.errMsg = f_err.Error()
			log.Printf("[LogTask %s] Log open failed, %s", p.LogPath, f_err)
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
