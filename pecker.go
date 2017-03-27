package logpeck

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Pecker struct {
	logTasks map[string]*LogTask
	mu       sync.Mutex
	db       *DB
}

func NewPecker(db *DB) (*Pecker, error) {
	pecker := &Pecker{
		logTasks: make(map[string]*LogTask),
		db:       db,
	}
	err := pecker.restorePeckTasks(db)
	if err != nil {
		return nil, err
	}
	return pecker, nil
}

func (p *Pecker) restorePeckTasks(db *DB) error {
	defer LogExecTime(time.Now(), "Restore PeckTaskConfig")
	configs, err := p.db.GetAllConfigs()
	if err != nil {
		return err
	}
	for i, config := range configs {
		p.AddPeckTask(&config)
		log.Printf("Restore PeckTask[%d] : %s", i, config)
	}
	return nil
}

func (p *Pecker) AddPeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		// Log file is not open. Open and tail it.
		var err error
		log_task, err = NewLogTask(log_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Add LogTask[%s] failed, err[%s].", log_path, err)
			return err
		}
	}

	if log_task.Exist(peck_conf) {
		return errors.New("Peck task already exist")
	}

	l_err := log_task.AddPeckTask(peck_conf)
	if l_err != nil {
		p.RemovePeckTask(peck_conf)
		return l_err
	}
	return nil
}

func (p *Pecker) UpdatePeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		// Log file is not open. Open and tail it.
		var err error
		log_task, err = NewLogTask(log_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Add LogTask[%s] failed, err[%s].", log_path, err)
			return err
		}
	}

	l_err := log_task.AddPeckTask(peck_conf)
	if l_err != nil {
		p.RemovePeckTask(peck_conf)
		return l_err
	}
	return nil
}

func (p *Pecker) RemovePeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return nil
	}
	log_task.RemovePeckTask(peck_conf)
	if log_task.Empty() {
		log_task.Close()
		delete(p.logTasks, log_path)
	}
	return nil
}

func (p *Pecker) StartPeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return nil
	}
	log_task.RemovePeckTask(peck_conf)
	if log_task.Empty() {
		log_task.Close()
		delete(p.logTasks, log_path)
	}
	return nil
}

func (p *Pecker) GetStat() *PeckerStat {
	p.mu.Lock()
	defer p.mu.Unlock()
	return nil
}
