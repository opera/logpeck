package logpeck

import (
	"encoding/json"
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
	Start() error
	Stop() error
}

func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) (*PeckTask, error) {
	if c.LogFormat == "text" {
		err := c.Check()
		if err != nil {
			log.Infof("[PeckTask] config check failed: %s", err)
			return nil, err
		}
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
	} else if c.SenderConfig.SenderName == "InfluxDbConfig" {
		sender = NewInfluxDbSender(&c.SenderConfig, c.Fields)
		interval := c.SenderConfig.Config.(InfluxDbConfig).Interval
		aggregatorConfigs := c.SenderConfig.Config.(InfluxDbConfig).AggregatorConfigs
		aggregator = NewAggregator(interval, &aggregatorConfigs)
	} else if c.SenderConfig.SenderName == "KafkaConfig" {
		sender = NewKafkaSender(&c.SenderConfig, c.Fields)
	} else {
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

func (p *PeckTask) Start() error {
	p.Stat.Stop = false
	if err := p.sender.Start(); err != nil {
		return err
	}
	return nil
}

func (p *PeckTask) Stop() error {
	p.Stat.Stop = true
	if err := p.sender.Stop(); err != nil {
		return err
	}
	return nil
}

func (p *PeckTask) IsStop() bool {
	return p.Stat.Stop
}

func (p *PeckTask) ExtractFieldsFromPlain(content string) map[string]interface{} {
	if len(p.Config.Fields) == 0 {
		return map[string]interface{}{"_Log": content}
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
		return map[string]interface{}{"_Log": content, "_Exception": err.Error()}
	}
	mContent, mErr := jContent.Map()
	if mErr != nil {
		return map[string]interface{}{"_Log": content, "_Exception": mErr.Error()}
	}
	if len(p.Config.Fields) == 0 {
		return map[string]interface{}{"_Log": content}
	}
	for _, field := range p.Config.Fields {
		key := SplitString(field.Name, ".")
		value := ""
		length := len(key)
		tmp := mContent
		for i := 0; i < length; i++ {
			if i == length-1 {
				if v, ok := tmp[key[i]].(string); ok {
					value = v
				} else if v, ok := tmp[key[i]].(json.Number); ok {
					value = v.String()
				} else {
					value = fmt.Sprintf("unknown type %v", tmp[key[i]])
				}
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
	} else if p.Config.SenderConfig.SenderName == "KafkaConfig" {
		fields := p.ExtractFields(content)
		p.sender.Send(fields)
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
