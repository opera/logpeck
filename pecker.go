package logpeck

import (
	"sync"
)

type Pecker struct {
	LogTasks map[string]LogTask

	mu sync.Mutex
}

func NewPecker() (*Pecker, error) {
	pecker := &Pecker{
		LogTasks: make(map[string]LogTask),
	}
	return pecker, nil
}

func (p *Pecker) AddPeckTask(peck_conf PeckTaskConfig) error {
	return nil
}

func (p *Pecker) RomovePeckTask(peck_conf PeckTaskConfig) error {
	return nil
}

func (p *Pecker) GetStat() *PeckerStat {
	return nil
}
