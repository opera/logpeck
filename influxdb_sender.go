package logpeck

import (
	"sync"
	log "github.com/Sirupsen/logrus"
)

type InfluxDbConfig struct {
	Hosts  string             `json:"hosts"`
	Interval int              `json:"interval"`
	FieldName string          `json:"fieldName"`
	Tables map[string]Table   `json:"tables"`
}

type Table struct {
	Measurement string        `json:"measurement"`
	Tags     []Tag            `json:"tags"`
	Aggregations   []Aggregation    `json:"fields"`
	Time     int              `json:"time"`
}

type Tag struct {
	TagName   string              `json:"tagName"`
	Column string              `json:"column"`
}

type Aggregation struct {
    AggName     Tag           `json:"aggName"`
	Cnt         bool          `json:"cnt"`
	Sum         bool          `json:"sum"`
	Avg         bool          `json:"avg"`
	Min         bool          `json:"min"`
	Max         bool          `json:"max"`
	Percentile  bool          `json:"percentile"`
	Percentiles []string      `json:"percentiles"`
}

type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	buckets       map[string]map[string][]int
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField) *Sender {
	sender := Sender{}
	sender.name = senderConfig.Name
	config := senderConfig.Config.(InfluxDbConfig)
	sender.senders = InfluxDbSender{
		config: config,
		fields: fields,
	}
	return &sender
}

func (p *InfluxDbSender)Send () {
	log.Infof("---------/n ")
}