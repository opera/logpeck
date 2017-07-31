package logpeck

import (
	"fmt"
	sjson "github.com/bitly/go-simplejson"
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

func (p *PeckTask) ExtractFieldsFromPlain(content string) map[string]interface{} {
	if len(p.Config.Fields) == 0 {
		return map[string]interface{}{"Log": content}
	}
	fields := make(map[string]interface{})
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

func (p *PeckTask) ExtractFieldsFromJson(content string) map[string]interface{} {
	fields := make(map[string]interface{})
	jContent, err := sjson.NewJson([]byte(content))
	if err != nil {
		return map[string]interface{}{"Log": content, "Exception": err.Error()}
	}
	mContent, mErr := jContent.Map()
	if mErr != nil {
		return map[string]interface{}{"Log": content, "Exception": mErr.Error()}
	}
	if len(p.Config.Fields) == 0 {
		for k, v := range mContent {
			fields[k] = fmt.Sprintf("%v", v)
		}
	}
	for _, field := range p.Config.Fields {
		fields[field.Name] = fmt.Sprintf("%v", mContent[field.Name])
	}
	return fields
}

func (p *PeckTask) ExtractFields(content string) map[string]interface{} {
	if p.Config.LogFormat == "json" {
		return p.ExtractFieldsFromJson(content)
	} else {
		return p.ExtractFieldsFromPlain(content)
	}
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
