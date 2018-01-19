package logpeck

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
	"time"
)

var FormatTime map[string]string = map[string]string{
	"ANSIC":       "Mon Jan _2 15:04:05 2006",
	"UnixDate":    "Mon Jan _2 15:04:05 MST 2006",
	"RubyDate":    "Mon Jan 02 15:04:05 -0700 2006",
	"RFC822":      "02 Jan 06 15:04 MST",
	"RFC822Z":     "02 Jan 06 15:04 -0700", // RFC822 with numeric zone
	"RFC850":      "Monday, 02-Jan-06 15:04:05 MST",
	"RFC1123":     "Mon, 02 Jan 2006 15:04:05 MST",
	"RFC1123Z":    "Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123 with numeric zone
	"RFC3339":     "2006-01-02T15:04:05Z07:00",
	"RFC3339Nano": "2006-01-02T15:04:05.999999999Z07:00",
	"Kitchen":     "3:04PM",
	// Handy time stamps.
	"Stamp":      "Jan _2 15:04:05",
	"StampMilli": "Jan _2 15:04:05.000",
	"StampMicro": "Jan _2 15:04:05.000000",
	"StampNano":  "Jan _2 15:04:05.000000000",
}

type AggregatorConfig struct {
	Enable            bool               `json:Enable`
	Interval          int64              `json:"Interval"`
	AggregatorOptions []AggregatorOption `json:"AggregatorOptions"`
}

type AggregatorOption struct {
	PreMeasurment string   `json:"PreMeasurment"`
	Measurment    string   `json:"Measurment"`
	Target        string   `json:"Target"`
	Tags          []string `json:"Tags"`
	Aggregations  []string `json:"Aggregations"`
	Timestamp     string   `json:"Timestamp"`
}

type Aggregator struct {
	aggregatorConfig AggregatorConfig
	buckets          map[string]map[string][]int64
	postTime         int64
}

func NewAggregator(aggregatorConfig *AggregatorConfig) *Aggregator {

	aggregator := &Aggregator{
		aggregatorConfig: *aggregatorConfig,
		buckets:          make(map[string]map[string][]int64),
		postTime:         0,
	}
	return aggregator
}

func getSampleTime(ts int64, interval int64) int64 {
	return ts / interval
}

func (p *Aggregator) IsEnable() bool {
	return p.aggregatorConfig.Enable
}

func (p *Aggregator) IsDeadline(timestamp int64) bool {
	interval := p.aggregatorConfig.Interval
	nowTime := getSampleTime(timestamp, interval)
	if p.postTime != nowTime {
		return true
	}
	return false
}

func (p *Aggregator) Record(fields map[string]interface{}) int64 {
	var now int64
	for i := 0; i < len(p.aggregatorConfig.AggregatorOptions); i++ {
		tags := p.aggregatorConfig.AggregatorOptions[i].Tags
		target := p.aggregatorConfig.AggregatorOptions[i].Target
		timestamp := p.aggregatorConfig.AggregatorOptions[i].Timestamp
		bucketName := p.aggregatorConfig.AggregatorOptions[i].PreMeasurment + "_" + p.aggregatorConfig.AggregatorOptions[i].Measurment + "_" + target
		bucketTag := ""
		if p.aggregatorConfig.AggregatorOptions[i].PreMeasurment != "" {
			bucketTag += p.aggregatorConfig.AggregatorOptions[i].PreMeasurment + "_"
		}
		if p.aggregatorConfig.AggregatorOptions[i].Measurment == "_default" {
			bucketTag += target
		} else {
			measurment, ok := fields[p.aggregatorConfig.AggregatorOptions[i].Measurment].(string)
			if !ok {
				log.Debug("[Record] Fields[measurment] format error: Fields[measurment] must be a string")
				now = time.Now().Unix()
				continue
			}
			bucketTag += measurment + "_" + target
		}

		//get time
		var err error
		timestamp_tmp, ok := fields[timestamp].(string)
		if !ok {
			now = time.Now().Unix()
		} else {
			now, err = strconv.ParseInt(timestamp_tmp, 10, 64)
			if err != nil {
				log.Debug("[Record] timestamp:%v can't use strconv.ParseInt", timestamp_tmp)
				now = time.Now().Unix()
			}
		}

		if target == "" {
			log.Error("[Record] Target is error: Target is null")
			return time.Now().Unix()
		}
		for i := 0; i < len(tags); i++ {
			tags_tmp, ok := fields[tags[i]].(string)
			if !ok {
				log.Debug("[Record] Fields[tag] format error: Fields[tag] must be a string")
			} else {
				bucketTag += "," + tags[i] + "=" + tags_tmp
			}
		}

		aggValue, ok := fields[target].(string)
		if !ok {
			log.Error("[Record] Fields[aggValue] format error: Fields[aggValue] must be a string")
			return now
		}
		if _, ok := p.buckets[bucketName]; !ok {
			p.buckets[bucketName] = make(map[string][]int64)
		}
		aggValueInt, err := strconv.ParseInt(aggValue, 10, 64)
		if err != nil {
			log.Debug("[Record] target:%v can't use strconv.ParseInt", aggValue)
			p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], 1)
		} else {
			p.buckets[bucketName][bucketTag] = append(p.buckets[bucketName][bucketTag], aggValueInt)
		}
	}
	return now
}

func quickSort(values []int64, left, right int64) {
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
	min := int64(0)
	max := int64(0)
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
	avg = sum / cnt
	for i := 0; i < len(aggregations); i++ {
		switch aggregations[i] {
		case "cnt":
			aggregationResults["cnt"] = int64(len(targetValue))
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

func (p *Aggregator) Dump(timestamp int64) map[string]interface{} {
	fields := map[string]interface{}{}
	log.Debug("[Dump] bucket is : %v", p.buckets)
	//now := strconv.FormatInt(timestamp, 10)
	for bucketName, bucketTag_value := range p.buckets {
		aggregations := []string{}
		for i := 0; i < len(p.aggregatorConfig.AggregatorOptions); i++ {
			if p.aggregatorConfig.AggregatorOptions[i].PreMeasurment+"_"+p.aggregatorConfig.AggregatorOptions[i].Measurment+"_"+p.aggregatorConfig.AggregatorOptions[i].Target == bucketName {
				aggregations = p.aggregatorConfig.AggregatorOptions[i].Aggregations
				break
			}
		}
		for bucketTag, targetValue := range bucketTag_value {
			fields[bucketTag] = getAggregation(targetValue, aggregations)
		}
	}
	fields["timestamp"] = timestamp
	p.postTime = getSampleTime(timestamp, p.aggregatorConfig.Interval)
	p.buckets = map[string]map[string][]int64{}
	log.Debug("[Dump] fields is : %v", fields)
	return fields
}
