package logpeck

import (
	log "github.com/Sirupsen/logrus"
	"sync"
)

/*
[{
"name":"module",
"tags":[
	"upstream",
	"downstream"
],
"aggr":["cnt","p99","avg"],
"target":"cost"
}]
*/

type InfluxDbConfig struct {
	Hosts       string                      `json:"hosts"`
	Interval    int64                       `json:"interval"`
	Name        string                      `json:"name"`
	Measurments map[string]MeasurmentConfig `json:"measurments"`
}

type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField) *InfluxDbSender {
	config := senderConfig.Config.(InfluxDbConfig)
	//measurments := make(map[string]*Measurment)
	/*measurments := &Measurment{
		config:   config,
		buckets: make(map[string]map[string][]int),
		postTime: 0,
	}*/
	//postTime := int64(0)
	sender := InfluxDbSender{
		config: config,
		fields: fields,
	}
	return &sender
}

func (p *InfluxDbSender) Send(fields string) {
	log.Infof("%s", fields)
	//p.measurments.MeasurmentRecall(fields)

}
