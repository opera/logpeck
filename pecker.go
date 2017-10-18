package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

type Pecker struct {
	logTasks map[string]*LogTask
	nametopath map[string] string
	mu       sync.Mutex
	db       *DB
	stop     bool
}

func NewPecker(db *DB) (*Pecker, error) {
	pecker := &Pecker{
		logTasks: make(map[string]*LogTask),
		nametopath :make(map[string] string),
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
		stat, _ := p.db.GetStat(config.LogPath, config.Name)
		p.AddPeckTask(&config, stat)
		log.Infof("[Pecker] Restore PeckTask[%d] : %s", i, config)
	}
	return nil
}

func (p *Pecker) AddPeckTask(config *PeckTaskConfig, stat *PeckTaskStat) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log.Infof("[Pecker] AddPeckTask %s", *config)
	log_path := config.LogPath
	log_task, ok := p.logTasks[log_path]
	if !ok {
		if _, ok2 := p.nametopath[config.Name];ok2{
			return errors.New("Peck task name already exist")
		} else{
		    log_task = NewLogTask(log_path)
		    p.logTasks[log_path] = log_task
		    p.nametopath[config.Name]=log_path
		}
 	}

	if log_task.Exist(config) {
		return errors.New("Peck task already exist")
	}

	task, err := NewPeckTask(config, stat)
	if err != nil {
		return err
	}

	{
		err1 := db.SaveConfig(&task.Config)
		err2 := db.SaveStat(&task.Stat)
		if err1 != nil || err2 != nil {
			panic(err1.Error() + " " + err2.Error())
		}
	}
	err = log_task.AddPeckTask(task)
	if err != nil {
		// AddPeckTask must be successful
		panic(err)
	}
	return nil
}

func (p *Pecker) UpdatePeckTask(config *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log.Infof("[Pecker] UpdatePeckTask %s", *config)
	if _, ok := p.nametopath[config.Name];!ok{
		return errors.New("Peck task name not exist")
	}
	log_path := p.nametopath[config.Name]
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return errors.New("Peck task not exist")
	}

	if !log_task.Exist(config) {
		return errors.New("Peck task not exist")
	}

	task, err := NewPeckTask(config, nil)
	if err != nil {
		return err
	}

	{
		err := db.SaveConfig(&task.Config)
		if err != nil {
			panic(err.Error())
		}
	}
	err = log_task.UpdatePeckTask(task)
	if err != nil {
		// UpdatePeckTask must be successful
		panic(err)
	}
	return nil
}

func (p *Pecker) RemovePeckTask(config *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.nametopath[config.Name];!ok{
		return errors.New("Peck task name not exist")
	}
	log_path := p.nametopath[config.Name]
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return errors.New("Task not exist")
	}

	{
		log.Infof("[Pecker] Remove PeckTask try clean db: %s", config)
		err1 := db.RemoveConfig(log_path, config.Name)
		err2 := db.RemoveStat(log_path, config.Name)
		if err1 != nil || err2 != nil {
			panic(err1.Error() + " " + err2.Error())
		}
	}

	log_task.RemovePeckTask(config)
	delete(p.nametopath,config.Name)
	if log_task.Empty() {
		log_task.Close()
		delete(p.logTasks, log_path)
	}
	log.Infof("[Pecker] Remove PeckTask finish: %s", config)
	return nil
}

func (p *Pecker) ListPeckTask() ([]PeckTaskConfig, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	configs, err := p.db.GetAllConfigs()
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (p *Pecker) StartPeckTask(config *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	log_path := p.nametopath[config.Name]
	log_task, ok := p.logTasks[log_path]
	if !ok {
		return errors.New("Task not exist")
	}

	{
		// Try update peck task stat in boltdb
		// stat, err := db.GetStat(config.LogPath, config.Name)
		stat, err := db.GetStat(log_path, config.Name)
		if err != nil {
			return err
		}
		if !stat.Stop {
			return errors.New("Task already started")
		}
		stat.Stop = false
		err = db .SaveStat(stat)
	}
	if log_task.IsStop() {
		log_task.Start()
	}

	return log_task.StartPeckTask(config)
}

func (p *Pecker) StopPeckTask(config *PeckTaskConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.nametopath[config.Name];!ok{
		return errors.New("Peck task name not exist")
	}
	log_path := p.nametopath[config.Name]
	log_task, ok := p.logTasks[log_path]
	if !ok {
			return errors.New("Task not exist")
	}

	{
		// Try update peck task stat in boltdb
		stat, err := db.GetStat(log_path, config.Name)
		if err != nil {
			return err
		}
		if stat.Stop {
			return errors.New("Task already stopped")
		}
		stat.Stop = true
		err = db.SaveStat(stat)
	}

	return log_task.StopPeckTask(config)
}

func (p *Pecker) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.stop {
		return errors.New("Pecker already started")
	}
	for path, logTask := range p.logTasks {
		log.Infof("[Pecker] Start LogTask %s", path)
		logTask.Start()
	}
	return nil
}

func (p *Pecker) GetStat() *PeckerStat {
	p.mu.Lock()
	defer p.mu.Unlock()
	return nil
}

