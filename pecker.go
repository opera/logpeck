package logpeck

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Pecker struct {
	logTasks map[string]*LogTask
	mu       sync.Mutex
	db       *DB
	stop     bool
}

func NewPecker(db *DB) (*Pecker, error) {
	pecker := &Pecker{
		logTasks: make(map[string]*LogTask),
		db:       db,
		stop:     true,
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
		log.Printf("[Pecker] Restore PeckTask[%d] : %s", i, config)
	}
	return nil
}

func (p *Pecker) AddPeckTask(config *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log.Printf("[Pecker] AddPeckTask %s", *config)
	log_path := config.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		log_task = NewLogTask(log_path)
		p.logTasks[log_path] = log_task
	}

	if log_task.Exist(config) {
		return errors.New("Peck task already exist")
	}

	task := NewPeckTask(config)

	{
		err1 := db.SaveConfig(&task.Config)
		err2 := db.SaveStat(&task.Stat)
		if err1 != nil || err2 != nil {
			panic(err1.Error() + " " + err2.Error())
		}
	}
	err := log_task.AddPeckTask(task)
	if err != nil {
		// AddPeckTask must be successful
		panic(err)
	}
	return nil
}

func (p *Pecker) UpdatePeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		log.Printf("[Pecker] Failed to UpdatePeckTask, PeckTask not exist")
		return fmt.Errorf("PeckTask not exist")
	}

	log_task.UpdatePeckTask(peck_conf)
	return nil
}

func (p *Pecker) RemovePeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return errors.New("Task Not Exist")
	}

	{
		log.Printf("[Pecker] Remove PeckTask try clean db: %s", peck_conf)
		err1 := db.RemoveConfig(peck_conf.LogPath, peck_conf.Name)
		err2 := db.RemoveStat(peck_conf.LogPath, peck_conf.Name)
		if err1 != nil || err2 != nil {
			panic(err1.Error() + " " + err2.Error())
		}
	}

	log_task.RemovePeckTask(peck_conf)
	if log_task.Empty() {
		log_task.Close()
		delete(p.logTasks, log_path)
	}
	log.Printf("[Pecker] Remove PeckTask finish: %s", peck_conf)
	return nil
}

func (p *Pecker) StartPeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return errors.New("Task not exist")
	}

	{
		// Try update peck task stat in boltdb
		stat, err := db.GetStat(peck_conf.LogPath, peck_conf.Name)
		if err != nil {
			return err
		}
		if !stat.Stop {
			return errors.New("Task already started")
		}
		stat.Stop = false
		err = db.SaveStat(stat)
	}

	return log_task.StartPeckTask(peck_conf)
}

func (p *Pecker) StopPeckTask(peck_conf *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := peck_conf.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return errors.New("Task not exist")
	}

	{
		// Try update peck task stat in boltdb
		stat, err := db.GetStat(peck_conf.LogPath, peck_conf.Name)
		if err != nil {
			return err
		}
		if stat.Stop {
			return errors.New("Task already stopped")
		}
		stat.Stop = true
		err = db.SaveStat(stat)
	}

	return log_task.StopPeckTask(peck_conf)
}

func (p *Pecker) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.stop {
		return errors.New("Pecker already started")
	}
	for path, logTask := range p.logTasks {
		log.Printf("[Pecker] Start LogTask %s", path)
		logTask.Start()
	}
	return nil
}

func (p *Pecker) GetStat() *PeckerStat {
	p.mu.Lock()
	defer p.mu.Unlock()
	return nil
}
