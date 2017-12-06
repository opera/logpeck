package logpeck

import (
	"sync"
	log "github.com/Sirupsen/logrus"
)

type InfluxDbConfig struct {
	Hosts     string          `json:"hosts"`
	Interval  int64           `json:"interval"`
	FieldName string          `json:"fieldName"`       //the column of measurement
	Tables map[string]Table   `json:"tables"`
}

type Table struct {
	Measurement  string          `json:"measurement"`
	Tags         []Tag           `json:"tags"`
	Aggregations []Aggregation   `json:"aggregations"`
	Time         string          `json:"time"`
}

type Tag struct {
	TagName   string           `json:"tagName"`
	Column    string           `json:"column"`
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
	buckets       map[string]map[string][]float64
	postTime      int64
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField) *Sender {
	sender := Sender{}
	sender.name = senderConfig.Name
	config := senderConfig.Config.(InfluxDbConfig)
	buckets :=make(map[string]map[string][]float64)
	postTime := int64(0)
	sender.senders = InfluxDbSender{
		config:    config,
		fields:    fields,
		postTime:  postTime,
		buckets:   buckets,
	}
	return &sender
}

func (p *InfluxDbSender)Send (now int) {

	log.Infof("--------------------------------------------------------%d\n",now)
	for k1,v1 := range p.buckets {
		log.Infof("%s",k1)
		for k2,v2 := range v1{
			log.Infof(" %s=%f",k2,v2)
		}
	}
	log.Infof("%v",p.buckets)
}