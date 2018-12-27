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
	Enable   bool               `json:Enable`
	Interval int64              `json:"Interval"`
	Options  []AggregatorOption `json:"Options"`
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
	config   AggregatorConfig
	buckets  map[string]map[string][]float64
	postTime int64
}

func NewAggregator(config *AggregatorConfig) *Aggregator {

	aggregator := &Aggregator{
		config:   *config,
		buckets:  make(map[string]map[string][]float64),
		postTime: 0,
	}
	return aggregator
}

func getSampleTime(ts int64, interval int64) int64 {
	return ts / interval
}

func (p *Aggregator) IsEnable() bool {
	return p.config.Enable
}

func (p *Aggregator) IsDeadline(timestamp int64) bool {
	interval := p.config.Interval
	nowTime := getSampleTime(timestamp, interval)
	if p.postTime != nowTime {
		return true
	}
	return false
}

func (p *Aggregator) Record(fields map[string]interface{}) int64 {
	var now int64
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
				log.Debugf("[Record] timestamp:%v can't use strconv.ParseInt", timestamp_tmp)
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
				log.Debugf("[Record] Fields[tag] format error: Fields[tag] must be a string")
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
	return now
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

func (p *Aggregator) Dump(timestamp int64) map[string]interface{} {
	fields := map[string]interface{}{}
	log.Debug("[Dump] bucket is", p.buckets)
	//now := strconv.FormatInt(timestamp, 10)
	for bucketName, bucketTag_value := range p.buckets {
		aggregations := []string{}
		for i := 0; i < len(p.config.Options); i++ {
			if p.config.Options[i].PreMeasurment+"_"+p.config.Options[i].Measurment+"_"+p.config.Options[i].Target == bucketName {
				aggregations = p.config.Options[i].Aggregations
				break
			}
		}
		for bucketTag, targetValue := range bucketTag_value {
			fields[bucketTag] = getAggregation(targetValue, aggregations)
		}
	}
	fields["timestamp"] = timestamp
	p.postTime = getSampleTime(timestamp, p.config.Interval)
	p.buckets = map[string]map[string][]float64{}
	log.Debug("[Dump] fields is", fields)
	return fields
}
