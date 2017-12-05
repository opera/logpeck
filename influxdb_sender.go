package logpeck

import (
	"sync"
)

type InfluxDbConfig struct {
	Hosts  string
	Tables []Table
}

type Table struct {
	post_interval int

	Tags     []string
	op       Filed
	cost     Filed
	status   Filed
	upstream Filed
	endpoint Filed
	req_len  Filed
	arti_cnt Filed
}

type Filed struct {
	name  string
	value string
}

type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(output *OutPutConfig, fields []PeckField) *Sender {
	sender := Sender{}
	sender.name = output.Name
	config := output.Config.(InfluxDbConfig)
	sender.sender = InfluxDbSender{
		config: config,
		fields: fields,
	}
	return &sender
}
