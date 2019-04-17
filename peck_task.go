package logpeck

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
)

// PeckTask .
type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat

	filter     PeckFilter
	extractor  Extractor
	sender     Sender
	aggregator *Aggregator
	postTime   int64
}

// NewPeckTask .
func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) (*PeckTask, error) {
	config := c
	var stat *PeckTaskStat
	if s == nil {
		stat = &PeckTaskStat{
			Name: c.Name,
			Stop: true,
		}
	} else {
		stat = s
	}
	extractor, err := NewExtractor(config.Extractor)
	if err != nil {
		return nil, err
	}
	filter := NewPeckFilter(config.Keywords)
	//var sender Sender
	sender, err := NewSender(&config.Sender)
	if err != nil {
		return nil, err
	}
	aggregator := NewAggregator(&config.Aggregator)
	task := &PeckTask{
		Config:     *config,
		Stat:       *stat,
		filter:     *filter,
		extractor:  extractor,
		sender:     sender,
		aggregator: aggregator,
	}
	if aggregator.IsEnable() {
		go task.tryDumpAggragator()
	}
	log.Infof("[PeckTask] new peck task %#v", task)
	return task, nil
}

// Start .
func (p *PeckTask) Start() error {
	p.Stat.Stop = false
	if err := p.sender.Start(); err != nil {
		return err
	}
	return nil
}

// Stop .
func (p *PeckTask) Stop() error {
	p.Stat.Stop = true
	if err := p.sender.Stop(); err != nil {
		return err
	}
	return nil
}

// IsStop .
func (p *PeckTask) IsStop() bool {
	return p.Stat.Stop
}

// Process .
func (p *PeckTask) Process(content string) {
	//log.Infof("sender%v",p.sender)
	if p.Stat.Stop {
		return
	}
	if p.filter.Drop(content) {
		return
	}

	fields, _ := p.extractor.Extract(content)
	if p.aggregator.IsEnable() {
		p.aggregator.Record(fields)
	} else {
		p.sender.Send(fields)
	}
}

// ProcessTest .
func (p *PeckTask) ProcessTest(content string) (map[string]interface{}, error) {
	if p.filter.Drop(content) {
		return map[string]interface{}{}, errors.New("Discarded")
	}
	fields, err := p.extractor.Extract(content)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return fields, nil
}

func (p *PeckTask) tryDumpAggragator() {
	for {
		now := time.Now().Unix()
		interval := p.aggregator.config.Interval
		if now/interval != p.postTime/interval {
			fields := p.aggregator.Dump()
			if len(fields) > 0 {
				p.sender.Send(fields)
			}
			p.postTime = now
		}
		time.Sleep(100 * time.Millisecond)
	}

}
