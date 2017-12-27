package logpeck

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
	"strconv"
)

type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat

	filter     PeckFilter
	fields     map[string]bool
	sender     Sender
	aggregator *Aggregator
}

type Sender interface {
	Send(map[string]interface{})
}

func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) (*PeckTask, error) {
	err := c.Check()
	if err != nil {
		log.Infof("[PeckTask] config check failed: %s", err)
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
	filter := NewPeckFilter(config.Keywords)
	var sender Sender
	aggregator := &Aggregator{}
	if c.SenderConfig.SenderName == "ElasticSearchConfig" {
		sender = NewElasticSearchSender(&c.SenderConfig, c.Fields)
	}

	if c.SenderConfig.SenderName == "InfluxDbConfig" {
		sender = NewInfluxDbSender(&c.SenderConfig, c.Fields)
		interval := c.SenderConfig.Config.(InfluxDbConfig).Interval
		aggregatorConfigs := c.SenderConfig.Config.(InfluxDbConfig).AggregatorConfigs
		aggregator = NewAggregator(interval, &aggregatorConfigs)
	}

	task := &PeckTask{
		Config:     *config,
		Stat:       *stat,
		filter:     *filter,
		sender:     sender,
		aggregator: aggregator,
	}
	return task, nil
}

func (p *PeckTask) Start() {
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

func FormatJsonValue(iValue interface{}) interface{} {
	if value, ok := iValue.([]*sjson.Json); ok {
		var valueArray []interface{}
		for _, e := range value {
			valueArray = append(valueArray, FormatJsonValue(e))
		}
		return valueArray
	} else if value, ok := iValue.(*sjson.Json); ok {
		m, _ := value.Map()
		ret := sjson.New()
		for k, v := range m {
			ret.Set(k, fmt.Sprint("%v", v))
		}
		return ret
	} else {
		return iValue
	}
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
		return mContent
	}
	for _, field := range p.Config.Fields {
		key := SplitString(field.Name, ".")
		value := ""
		length := len(key)
		tmp := mContent
		for i := 0; i < length; i++ {
			if i == length-1 {
				value = tmp[key[i]].(string)
				break
			}
			tmp = tmp[key[i]].(map[string]interface{})
		}
		fields[field.Name] = value
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
	//log.Infof("sender%v",p.sender)
	if p.Stat.Stop {
		return
	}
	if p.filter.Drop(content) {
		return
	}
	if p.Config.SenderConfig.SenderName == "ElasticSearchConfig" {
		fields := p.ExtractFields(content)
		p.sender.Send(fields)
	} else if p.Config.SenderConfig.SenderName == "InfluxDbConfig" {
		fields := p.ExtractFields(content)
		timestamp := p.aggregator.Record(fields)
		deadline := p.aggregator.IsDeadline(timestamp)
		if deadline {
			aggregationResults := p.aggregator.Dump(timestamp)
			p.sender.Send(aggregationResults)
		}
	}
}

func (p *PeckTask) ProcessTest(content string) (map[string]interface{}, error) {
	if p.filter.Drop(content) {
		var err error = errors.New("[peck_task]The line does not meet the rules ")
		s := make(map[string]interface{})
		return s, err
	}
	fields := p.ExtractFields(content)
	return fields, nil
}
