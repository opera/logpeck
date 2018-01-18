package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat

	filter     PeckFilter
	extractor  Extractor
	sender     Sender
	aggregator *Aggregator
}

type Sender interface {
	Send(map[string]interface{})
	Start() error
	Stop() error
}

func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) (*PeckTask, error) {
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
	extractor, err := NewExtractor(config.ExtractorConfig, config.Fields)
	if err != nil {
		return nil, err
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
		extractor:  extractor,
		sender:     sender,
		aggregator: aggregator,
	}
	log.Infof("[PeckTask] new peck task %#v", task)
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

func (p *PeckTask) Process(content string) {
	//log.Infof("sender%v",p.sender)
	if p.Stat.Stop {
		return
	}
	if p.filter.Drop(content) {
		return
	}
	if p.Config.SenderConfig.SenderName == "ElasticSearchConfig" {
		fields, _ := p.extractor.Extract(content)
		p.sender.Send(fields)
	} else if p.Config.SenderConfig.SenderName == "InfluxDbConfig" {
		fields, _ := p.extractor.Extract(content)
		timestamp := p.aggregator.Record(fields)
		deadline := p.aggregator.IsDeadline(timestamp)
		if deadline {
			aggregationResults := p.aggregator.Dump(timestamp)
			p.sender.Send(aggregationResults)
		}
	} else if p.Config.SenderConfig.SenderName == "KafkaConfig" {
		fields, _ := p.extractor.Extract(content)
		p.sender.Send(fields)
	}
}

func (p *PeckTask) ProcessTest(content string) (map[string]interface{}, error) {
	if p.filter.Drop(content) {
		var err error = errors.New("[peck_task]The line does not meet the rules ")
		s := make(map[string]interface{})
		return s, err
	}
	fields, _ := p.extractor.Extract(content)
	return fields, nil
}
