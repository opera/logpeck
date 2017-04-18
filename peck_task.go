package logpeck

import (
	"log"
	"strconv"
)

type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat

	filter PeckFilter
	fields map[string]bool
}

func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) (*PeckTask, error) {
	err := c.Check()
	if err != nil {
		log.Printf("[PeckTask] config check failed: %s", err)
		return nil, err
	}
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
	fields := make(map[string]bool)
	for _, v := range config.Fields {
		fields[v.Name] = true
	}
	filter := NewPeckFilter(config.FilterExpr)
	InitElasticSearchMapping(config)

	task := &PeckTask{
		Config: *config,
		Stat:   *stat,
		filter: *filter,
	}
	log.Printf("[PeckTask] NewPeckTask %+v", task)
	return task, nil
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

func (p *PeckTask) ExtractFields(content string) map[string]string {
	if len(p.Config.Fields) == 0 {
		return map[string]string{"Log": content}
	}
	fields := make(map[string]string)
	arr := SplitString(content, p.Config.Delimiters)
	for _, field := range p.Config.Fields {
		if field.Value[0] != '$' {
			panic(field)
		}
		pos, err := strconv.Atoi(field.Value[1:])
		if err != nil {
			panic(field)
		}
		if len(arr) < pos {
			continue
		}
		fields[field.Name] = arr[pos-1]
	}
	return fields
}

func (p *PeckTask) Process(content string) {
	if p.Stat.Stop {
		return
	}
	if p.filter.Drop(content) {
		return
	}
	fields := p.ExtractFields(content)
	SendToElasticSearch(&p.Config.ESConfig, fields)
}
