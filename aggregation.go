package logpeck

import (
	"sort"
	"strconv"
)

type MeasurmentConfig struct {
	Tags         []string `json:"tags"`
	Aggregations []string `json:"aggregations"`
	Target       string   `json:"target"`
	Time         string   `json:"time"`
}

type Measurment struct {
	config   InfluxDbConfig
	buckets  map[string]map[string][]int
	postTime int64
}

func NewMeasurment(senderConfig *SenderConfig) *Measurment {
	config := senderConfig.Config.(InfluxDbConfig)
	//measurments := make(map[string]*Measurment)
	measurments := &Measurment{
		config:   config,
		buckets:  make(map[string]map[string][]int),
		postTime: 0,
	}
	return measurments
}

func getSampleTime(ts int64, interval int64) int64 {
	return ts / interval
}
func (p *Measurment) startSend(time int64) (bool, int64) {
	interval := p.config.Interval
	nowTime := getSampleTime(time, interval)
	if p.postTime != nowTime {
		return true, nowTime
	}
	return false, nowTime
}

func (p *Measurment) Recall(fields map[string]interface{}) int64 {
	//get sender
	//influxDbConfig := p.Config.SenderConfig.Config.(InfluxDbConfig)
	bucketName := fields[p.config.Name].(string)
	bucketTag := ""
	measurement := p.config.Measurments[bucketName]
	tags := measurement.Tags
	aggregations := measurement.Aggregations
	target := measurement.Target
	time := measurement.Time
	for i := 0; i < len(tags); i++ {
		bucketTag += "," + tags[i] + "=" + fields[tags[i]].(string)
	}
	int_bool := false
	for i := 0; i < len(aggregations); i++ {
		if aggregations[i] != "cnt" {
			int_bool = true
		}
	}
	aggValue := fields[target].(string)

	if _, ok := p.buckets[bucketName]; !ok {
		p.buckets[bucketName] = make(map[string][]int)
	}
	if int_bool == false {
		p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], 1)
	} else {
		aggValue, err := strconv.Atoi(aggValue)
		if err != nil {
			panic(aggValue)
		}
		p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], aggValue)
	}

	//get time
	now, err := strconv.ParseInt(fields[time].(string), 10, 64)
	if err != nil {
		panic(fields)
	}
	return now
}

func GetAggregation(v2 []int, aggregations []string) string {
	aggregation := ""
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
			if aggregations[i][0] == 'p' {
				proportion, err := strconv.Atoi(aggregations[i][1:])
				if err != nil {
					panic(aggregations[i])
				}
				percentile := v2[cnt*proportion/100-1]
				str := strconv.Itoa(percentile)
				aggregation += aggregations[i] + "=" + str
			}
		}
		if i < len(aggregations)-1 {
			aggregation += ","
		}
	}
	return aggregation
}

func (p *Measurment) Output(now int64) string {
	measurment := ""
	for bucketName, fields := range p.buckets {
		for bucketTag, value := range fields {
			aggregation := GetAggregation(value, p.config.Measurments[bucketName].Aggregations)
			now := strconv.FormatInt(now, 10)
			measurment += bucketName + bucketTag + " " + aggregation + " " + now + "\n"
		}
	}
	return measurment
}
