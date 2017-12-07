package logpeck

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
	"strconv"
)

type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat

	filter PeckFilter
	fields map[string]bool
	sender Sender
}

type Sender struct {
	name    string
	senders interface{}
}

func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) (*PeckTask, error) {
	err := c.Check()
	if err != nil {
		log.Infof("[PeckTask] config check failed: %s", err)
		return nil, err
	}
	var config *PeckTaskConfig = c
	var stat *PeckTaskStat
	if s == nil {
		stat = &PeckTaskStat{
			Name:    c.Name,
			LogPath: c.LogPath,
			Stop:    true,
		}
	} else {
		stat = s
	}
	fields := make(map[string]bool)
	for _, v := range config.Fields {
		fields[v.Name] = true
	}
	filter := NewPeckFilter(config.FilterExpr)
	sender := &Sender{}
	if c.SenderConfig.Name == "ElasticSearchConfig" {
		sender = NewElasticSearchSender(&c.SenderConfig, c.Fields)
	}

	if c.SenderConfig.Name == "InfluxDbConfig" {
		sender = NewInfluxDbSender(&c.SenderConfig, c.Fields)
	}

	task := &PeckTask{
		Config: *config,
		Stat:   *stat,
		filter: *filter,
		sender: *sender,
	}
	return task, nil
}

func (p *PeckTask) Start() {
	p.Stat.Stop = false
}

func (p *PeckTask) Stop() {
	p.Stat.Stop = true
}

func (p *PeckTask) IsStop() bool {
	return p.Stat.Stop
}

func getSampleTime(ts int, interval int64) int64 {
	return int64(ts) / interval
}

func (p *PeckTask) ExtractFieldsFromPlain(content string) map[string]interface{} {
	if len(p.Config.Fields) == 0 {
		return map[string]interface{}{"Log": content}
	}
	fields := make(map[string]interface{})
	arr := SplitString(content, p.Config.Delimiters)
	for _, field := range p.Config.Fields {
		if field.Value[0] != '$' {
			panic(field)
		}
		pos, err := strconv.Atoi(field.Value[1:])
		if err != nil {
			panic(field)
		}
		if len(arr) < pos {
			continue
		}
		fields[field.Name] = arr[pos-1]
	}
	return fields
}

func (p *PeckTask) ExtractInfluxDbFieldsFromPlain(fields map[string]interface{}) {
	//get sender
	influxDbConfig := p.Config.SenderConfig.Config.(InfluxDbConfig)
	s := fields[influxDbConfig.Name].(string)
	measurement := influxDbConfig.Measurments[s]
	tags := measurement.Tags
	aggregations := measurement.Aggregations
	target := measurement.Target
	time := measurement.Time
	for i := 0; i < len(tags); i++ {
		s += "," + tags[i] + "=" + fields[tags[i]].(string)
	}
	sender := p.sender.senders.(InfluxDbSender)
	int_bool := false
	for i := 0; i < len(aggregations); i++ {
		if aggregations[i] != "cnt" {
			int_bool = true
		}
	}
	aggValue := fields[target].(string)

	if _, ok := sender.buckets[s]; !ok {
		sender.buckets[s] = make(map[string][]int)
	}
	if _, ok := sender.buckets[s][target]; !ok {
		//sender.buckets[s][aggName]=[]int{}
	}
	if int_bool == false {
		sender.buckets[s][target] = append(sender.buckets[s][target], 1)
	} else {
		aggValue, err := strconv.Atoi(aggValue)
		if err != nil {
			panic(aggValue)
		}
		sender.buckets[s][target] = append(sender.buckets[s][target], aggValue)
	}
	p.sender.senders = sender

	//get time
	interval := influxDbConfig.Interval
	now, err := strconv.Atoi(fields[time].(string))
	if err != nil {
		panic(fields)
	}
	nowTime := getSampleTime(now, interval)
	if sender.postTime != nowTime {
		sender.Send(now, aggregations)
		sender.postTime = nowTime
		sender.buckets = map[string]map[string][]int{}
		p.sender.senders = sender
	}

}

func FormatJsonValue(iValue interface{}) interface{} {
	if value, ok := iValue.([]*sjson.Json); ok {
		var valueArray []interface{}
		for _, e := range value {
			valueArray = append(valueArray, FormatJsonValue(e))
		}
		return valueArray
	} else if value, ok := iValue.(*sjson.Json); ok {
		m, _ := value.Map()
		ret := sjson.New()
		for k, v := range m {
			ret.Set(k, fmt.Sprint("%v", v))
		}
		return ret
	} else {
		return iValue
	}
}

func (p *PeckTask) ExtractFieldsFromJson(content string) map[string]interface{} {
	fields := make(map[string]interface{})
	jContent, err := sjson.NewJson([]byte(content))
	if err != nil {
		return map[string]interface{}{"Log": content, "Exception": err.Error()}
	}
	mContent, mErr := jContent.Map()
	if mErr != nil {
		return map[string]interface{}{"Log": content, "Exception": mErr.Error()}
	}
	if len(p.Config.Fields) == 0 {
		return mContent
	}
	for _, field := range p.Config.Fields {
		fields[field.Name] = mContent[field.Name]
	}
	return fields
}

func (p *PeckTask) ExtractElasticSearchFields(content string) map[string]interface{} {
	if p.Config.LogFormat == "json" {
		return p.ExtractFieldsFromJson(content)
	} else {
		return p.ExtractFieldsFromPlain(content)
	}
}

func (p *PeckTask) ExtractInfluxDbFields(content string) {
	fields := map[string]interface{}{}
	if p.Config.LogFormat == "json" {
		fields = p.ExtractFieldsFromJson(content)
	} else {
		fields = p.ExtractFieldsFromPlain(content)
	}
	p.ExtractInfluxDbFieldsFromPlain(fields)
}

func (p *PeckTask) Process(content string) {
	//log.Infof("sender%v",p.sender)
	if p.Stat.Stop {
		return
	}
	if p.filter.Drop(content) {
		return
	}
	if p.sender.name == "ElasticSearchConfig" {
		fields := p.ExtractElasticSearchFields(content)
		sender := p.sender.senders.(ElasticSearchSender)
		sender.Send(fields)
	}
	if p.sender.name == "InfluxDbConfig" {
		p.ExtractInfluxDbFields(content)
	}
}

func (p *PeckTask) ProcessTest(content string) (map[string]interface{}, error) {
	if p.filter.Drop(content) {
		var err error = errors.New("[peck_task]The line does not meet the rules ")
		s := make(map[string]interface{})
		return s, err
	}
	fields := p.ExtractElasticSearchFields(content)
	return fields, nil
}
