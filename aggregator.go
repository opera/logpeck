package logpeck

import (
	"github.com/Sirupsen/logrus"
	"strconv"
)

type AggregatorConfig struct {
	Tags         []string `json:"tags"`
	Aggregations []string `json:"aggregations"`
	Target       string   `json:"target"`
	Time         string   `json:"time"`
}

type Aggregator struct {
	interval          int64
	name              string
	aggregatorConfigs map[string]AggregatorConfig
	buckets           map[string]map[string][]int64
	postTime          int64
}

func NewAggregator(interval int64, name string, aggregators *map[string]AggregatorConfig) *Aggregator {
	aggregator := &Aggregator{
		interval:          interval,
		name:              name,
		aggregatorConfigs: *aggregators,
		buckets:           make(map[string]map[string][]int64),
		postTime:          0,
	}
	return aggregator
}

func getSampleTime(ts int64, interval int64) int64 {
	return ts / interval
}

func (p *Aggregator) IsDeadline(timestamp int64) bool {
	interval := p.interval
	nowTime := getSampleTime(timestamp, interval)
	if p.postTime != nowTime {
		return true
	}
	return false
}

func (p *Aggregator) Record(fields map[string]interface{}) int64 {
	//get sender
	//influxDbConfig := p.Config.SenderConfig.Config.(InfluxDbConfig)
	bucketName := fields[p.name].(string)
	bucketTag := ""
	aggregatorConfig := p.aggregatorConfigs[bucketName]
	tags := aggregatorConfig.Tags
	aggregations := aggregatorConfig.Aggregations
	target := aggregatorConfig.Target
	time := aggregatorConfig.Time
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

	//get time
	now, err := strconv.ParseInt(fields[time].(string), 10, 64)
	if err != nil {
		logrus.Debug("[Record] timestamp:%v can't use strconv.ParseInt", fields[time].(string))
		return now
	}

	if _, ok := p.buckets[bucketName]; !ok {
		p.buckets[bucketName] = make(map[string][]int64)
	}
	if int_bool == false {
		p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], 1)
	} else {
		aggValue, err := strconv.ParseInt(aggValue, 10, 64)
		if err != nil {
			logrus.Debug("[Record] target:%v can't use strconv.ParseInt", aggValue)
			return now
		}
		p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], aggValue)
	}
	return now
}

func quickSort(values []int64, left, right int) {

	temp := values[left]
	p := left
	i, j := left, right

	for i <= j {
		for j >= p && values[j] >= temp {
			j--
		}
		if j >= p {
			values[p] = values[j]
			p = j
		}

		for i <= p && values[i] <= temp {
			i++
		}
		if i <= p {
			values[p] = values[i]
			p = i
		}

	}

	values[p] = temp

	if p-left > 1 {
		quickSort(values, left, p-1)
	}
	if right-p > 1 {
		quickSort(values, p+1, right)
	}
}

func getAggregation(targetValue []int64, aggregations []string) map[string]int64 {
	aggregationResults := map[string]int64{}
	cnt := int64(len(targetValue))
	avg := int64(0)
	sum := int64(0)
	quickSort(targetValue, 0, len(targetValue)-1)
	for _, value := range targetValue {
		sum += value
	}
	avg = sum / cnt
	for i := 0; i < len(aggregations); i++ {
		switch aggregations[i] {
		case "cnt":
			aggregationResults["cnt"] = int64(len(targetValue))
		case "avg":
			aggregationResults["avg"] = avg
		default:
			if aggregations[i][0] == 'p' {
				proportion, err := strconv.ParseInt(aggregations[i][1:], 10, 64)
				if err != nil {
					panic(aggregations[i])
				}
				percentile := targetValue[cnt*proportion/100-1]
				aggregationResults[aggregations[i]] = percentile
			}
		}
	}
	return aggregationResults
}

func (p *Aggregator) Dump(timestamp int64) map[string]interface{} {
	fields := map[string]interface{}{}
	//now := strconv.FormatInt(timestamp, 10)
	for bucketName, bucketTag_value := range p.buckets {
		for bucketTag, targetValue := range bucketTag_value {
			aggregations := p.aggregatorConfigs[bucketName].Aggregations
			fields[bucketName+bucketTag] = getAggregation(targetValue, aggregations)
		}
	}
	fields["timestamp"] = timestamp
	p.postTime = getSampleTime(timestamp, p.interval)
	p.buckets = map[string]map[string][]int64{}
	return fields
}
