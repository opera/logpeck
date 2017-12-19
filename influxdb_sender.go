package logpeck

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
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
	Aggregators map[string]AggregatorConfig `json:"aggregators"`
}

type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField) *InfluxDbSender {
	config := senderConfig.Config.(InfluxDbConfig)
	sender := InfluxDbSender{
		config: config,
		fields: fields,
	}
	return &sender
}

func toInfluxdbLine(fields map[string]interface{}) string {
	lines := ""
	timestamp := fields["timestamp"].(int64)
	for k, v := range fields {
		if k == "timestamp" {
			continue
		}
		aggregationResults := v.(map[string]int)
		lines = k + " "
		for aggregation, result := range aggregationResults {
			lines += aggregation + "=" + strconv.Itoa(result) + ","

		}
		length := len(lines)
		lines = lines[0:length-1] + " " + strconv.FormatInt(timestamp, 10) + "\n"
	}
	return lines
}

func (p *InfluxDbSender) Send(fields map[string]interface{}) {
	lines := toInfluxdbLine(fields)
	log.Infof("%s", lines)
	//p.measurments.MeasurmentRecall(fields)
}
