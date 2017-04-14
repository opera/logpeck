package logpeck

import (
	"log"
)

type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat

	filter PeckFilter
}

func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) *PeckTask {
	var config *PeckTaskConfig = c
	var stat *PeckTaskStat
	if s == nil {
		stat = &PeckTaskStat{
			Name:    c.Name,
			LogPath: c.LogPath,
			Stop:    true,
		}
	} else {
		stat = s
	}
	filter := NewPeckFilter(config.FilterExpr)
	InitElasticSearchMapping(&config.ESConfig)

	task := &PeckTask{
		Config: *config,
		Stat:   *stat,
		filter: *filter,
	}
	log.Printf("[PeckTask] NewPeckTask %+v", task)
	return task
}

func (p *PeckTask) Start() {
	log.Printf("[PeckTask] Start")
	p.Stat.Stop = false
}

func (p *PeckTask) Stop() {
	p.Stat.Stop = true
}

func (p *PeckTask) IsStop() bool {
	return p.Stat.Stop
}

func (p *PeckTask) Process(content string) {
	if p.Stat.Stop {
		return
	}
	if p.filter.Drop(content) {
		return
	}
	SendToElasticSearch(&p.Config.ESConfig, content)
}
