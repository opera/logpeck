package logpeck

import (
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

// AggregatorConfig .
type AggregatorConfig struct {
	Enable   bool               `json:"Enable"`
	Interval int64              `json:"Interval"`
	Options  []AggregatorOption `json:"Options"`
}

// AggregatorOption .
type AggregatorOption struct {
	PreMeasurment string   `json:"PreMeasurment"`
	Measurment    string   `json:"Measurment"`
	Target        string   `json:"Target"`
	Tags          []string `json:"Tags"`
	Aggregations  []string `json:"Aggregations"`
	Timestamp     string   `json:"Timestamp"`
}

// Aggregator .
type Aggregator struct {
	config     AggregatorConfig
	buckets    map[string]map[string][]float64
	recordTime int64
	mu         sync.Mutex
}

// NewAggregator create aggregator
func NewAggregator(config *AggregatorConfig) *Aggregator {
	aggregator := &Aggregator{
		config:     *config,
		buckets:    make(map[string]map[string][]float64),
		recordTime: 0,
	}
	return aggregator
}

// IsEnable return true if enable
func (p *Aggregator) IsEnable() bool {
	return p.config.Enable
}

// Record fields
func (p *Aggregator) Record(fields map[string]interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.recordTime = time.Now().Unix()
	for i := 0; i < len(p.config.Options); i++ {
		tags := p.config.Options[i].Tags
		target := p.config.Options[i].Target
		timestamp := p.config.Options[i].Timestamp
		bucketName := p.config.Options[i].PreMeasurment + "_" + p.config.Options[i].Measurment + "_" + target
		bucketTag := ""
		if p.config.Options[i].PreMeasurment != "" {
			bucketTag += p.config.Options[i].PreMeasurment + "_"
		}
		if p.config.Options[i].Measurment == "_default" {
			bucketTag += target
		} else {
			measurment, ok := fields[p.config.Options[i].Measurment].(string)
			if !ok {
				log.Debug("[Record] Fields[measurment] format error: Fields[measurment] must be a string")
				continue
			}
			bucketTag += measurment + "_" + target
		}

		//get time
		var err error
		timestampTmp, ok := fields[timestamp].(string)
		if ok {
			ts, err := strconv.ParseInt(timestampTmp, 10, 64)
			if err == nil {
				p.recordTime = ts
			} else {
				log.Debugf("[Record] timestamp:%v can't use strconv.ParseInt", timestampTmp)
			}
		}

		if target == "" {
			log.Debug("[Record] Target is error: Target is null")
			return
		}
		for i := 0; i < len(tags); i++ {
			tagsTmp, ok := fields[tags[i]].(string)
			if !ok {
				log.Debugf("[Record] Fields[tag] format error: Fields[tag] must be a string")
			} else {
				bucketTag += "," + tags[i] + "=" + tagsTmp
			}
		}

		aggValue, ok := fields[target].(string)
		if !ok {
			log.Debugf("[Record] Fields[aggValue] format error: %v", fields[target])
			return
		}
		if _, ok := p.buckets[bucketName]; !ok {
			p.buckets[bucketName] = make(map[string][]float64)
		}
		aggValueFloat64, err := strconv.ParseFloat(aggValue, 64)
		if err != nil {
			log.Debugf("[Record] target:%v can't use strconv.ParseFloat", aggValue)
			p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], -1)
		} else {
			p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], aggValueFloat64)
		}
	}
	return
}

func quickSort(values []float64, left, right int64) {
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

func getAggregation(targetValue []float64, aggregations []string) map[string]float64 {
	aggregationResults := map[string]float64{}
	cnt := int64(len(targetValue))
	avg := float64(0)
	sum := float64(0)
	min := float64(0)
	max := float64(0)
	if cnt > 0 {
		min = targetValue[0]
		max = targetValue[0]
	}
	quickSort(targetValue, int64(0), int64(len(targetValue)-1))
	for _, value := range targetValue {
		sum += value
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	avg = sum / float64(cnt)
	for i := 0; i < len(aggregations); i++ {
		switch aggregations[i] {
		case "cnt":
			aggregationResults["cnt"] = float64(len(targetValue))
		case "sum":
			aggregationResults["sum"] = sum
		case "avg":
			aggregationResults["avg"] = avg
		case "min":
			aggregationResults["min"] = min
		case "max":
			aggregationResults["max"] = max
		default:
			if aggregations[i][0] == 'p' {
				proportion, err := strconv.ParseInt(aggregations[i][1:], 10, 64)
				if err != nil {
					panic(aggregations[i])
				}
				index := cnt*proportion/100 - 1
				if cnt*proportion/100-1 < 0 {
					index = 0
				}
				percentile := targetValue[index]
				aggregationResults[aggregations[i]] = percentile
			}
		}
	}
	return aggregationResults
}

// Dump aggregation result
func (p *Aggregator) Dump() map[string]interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	fields := map[string]interface{}{}
	log.Debug("[Dump] bucket is", p.buckets)
	//now := strconv.FormatInt(timestamp, 10)
	for bucketName, bucketTagValue := range p.buckets {
		aggregations := []string{}
		for i := 0; i < len(p.config.Options); i++ {
			if p.config.Options[i].PreMeasurment+"_"+p.config.Options[i].Measurment+"_"+p.config.Options[i].Target == bucketName {
				aggregations = p.config.Options[i].Aggregations
				break
			}
		}
		for bucketTag, targetValue := range bucketTagValue {
			fields[bucketTag] = getAggregation(targetValue, aggregations)
		}
	}
	fields["timestamp"] = p.recordTime
	p.buckets = map[string]map[string][]float64{}
	log.Debug("[Dump] fields is", fields)
	return fields
}
