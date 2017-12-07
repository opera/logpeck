package logpeck

import (
	log "github.com/Sirupsen/logrus"
	"sort"
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
	Hosts       string                `json:"hosts"`
	Interval    int64                 `json:"interval"`
	Name        string                `json:"name"`
	Measurments map[string]Measurment `json:"measurments"`
}

type Measurment struct {
	Tags         []string `json:"tags"`
	Aggregations []string `json:"aggregations"`
	Target       string   `json:"target"`
	Time         string   `json:"time"`
}

/*
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
*/
type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	buckets       map[string]map[string][]int
	postTime      int64
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField) *Sender {
	sender := Sender{}
	sender.name = senderConfig.Name
	config := senderConfig.Config.(InfluxDbConfig)
	buckets := make(map[string]map[string][]int)
	postTime := int64(0)
	sender.senders = InfluxDbSender{
		config:   config,
		fields:   fields,
		postTime: postTime,
		buckets:  buckets,
	}
	return &sender
}

func (p *InfluxDbSender) Send(now int, aggregations []string) {

	for k1, v1 := range p.buckets {
		for k2, v2 := range v1 {
			aggregation := " "
			cnt := len(v2)
			avg := 0
			sum := 0
			sort.Ints(v2)
			for _, value := range v2 {
				sum += value
				avg = sum / cnt
			}
			for i := 0; i < len(aggregations); i++ {
				switch aggregations[i] {
				case "cnt":
					str := strconv.Itoa(cnt)
					aggregation += "cnt=" + str
				case "avg":
					str := strconv.Itoa(avg)
					aggregation += "avg=" + str
				default:
					if k2[0] == 'p' {
						proportion, err := strconv.Atoi(k2[1:])
						if err != nil {
							panic(k2)
						}
						percentile := v2[cnt*proportion/100-1]
						str := strconv.Itoa(percentile)
						aggregation += k2 + "=" + str
					}
				}
				if i < len(aggregations)-1 {
					aggregation += ","
				}
				log.Infof("----------------------------")
				log.Infof("%s%s %d", k1, aggregation, now)
			}
		}
	}
}
