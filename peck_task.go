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
	name   string
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

func getSampleTime(ts int,interval int64) int64 {
	return int64(ts)/interval
}

func (p *PeckTask) ExtractElasticSearchFieldsFromPlain(content string) map[string]interface{} {
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

func (p *PeckTask) ExtractInfluxDbFieldsFromPlain(content string)  {

	influxDbConfig :=p.Config.SenderConfig.Config.(InfluxDbConfig)
	if influxDbConfig.FieldName[0] != '$'{
		panic(influxDbConfig.FieldName)
	}
	arr := SplitString(content, p.Config.Delimiters)
	pos, err := strconv.Atoi(influxDbConfig.FieldName[1:])
	if err != nil {
		panic(influxDbConfig.FieldName)
	}
	log.Infof("pos:%v",pos)
	if len(arr) < pos {

	}
	filename := arr[pos-1]                                       //the column of measurement
	measurement := influxDbConfig.Tables[filename].Measurement   //the value of measurement
	tags := influxDbConfig.Tables[filename].Tags                 //the value of tags
	aggregations := influxDbConfig.Tables[filename].Aggregations //the value of fields
	s := measurement
	for i := 0; i < len(tags); i++ {
		if tags[i].Column[0] != '$'{
			panic(tags[i].Column)
		}
		pos, err := strconv.Atoi(tags[i].Column[1:])
		if err != nil {
			panic(influxDbConfig.FieldName)
		}
		if len(arr) < pos {

		}
		tagValue := arr[pos-1]
		s+=","+tags[i].TagName+"="+tagValue
	}
	sender := p.sender.senders.(InfluxDbSender)
	for i := 0; i < len(aggregations) ; i++ {
		if aggregations[i].AggName.Column[0] != '$'{
			panic(aggregations[i].AggName.Column)
		}
		aggName := aggregations[i].AggName.TagName
		pos, err := strconv.Atoi(aggregations[i].AggName.Column[1:])
		if err != nil {
			panic(influxDbConfig.FieldName)
		}
		if len(arr) < pos {

		}
		aggValue,err:=strconv.ParseFloat(arr[pos-1],32/64)
		if _,ok:=sender.buckets[s];!ok{
			sender.buckets[s]=make(map[string][]float64)
		}
		if _,ok:=sender.buckets[s][aggName];!ok{
			//sender.buckets[s][aggName]=[]int{}
		}
		if aggregations[i].Cnt == true {
			sender.buckets[s][aggName]=append(sender.buckets[s][aggName],1)
		} else{
		sender.buckets[s][aggName]=append(sender.buckets[s][aggName],aggValue)
		}
	}
	p.sender.senders = sender

	interval := influxDbConfig.Interval
	if influxDbConfig.Tables[filename].Time[0] != '$'{
		panic(influxDbConfig.FieldName)
	}
	pos, err = strconv.Atoi(influxDbConfig.Tables[filename].Time[1:])
	if err != nil {
		panic(influxDbConfig.Tables[filename].Time)
	}
	if len(arr) < pos {

	}
	now,err:=strconv.Atoi(arr[pos-1])
	nowTime := getSampleTime(now,interval)

	if sender.postTime != nowTime {
		sender.Send(now)
		sender.postTime = nowTime
		sender.buckets=map[string]map[string][]float64{}
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

func (p *PeckTask) ExtractElasticSearchFieldsFromJson(content string) map[string]interface{} {
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
		return p.ExtractElasticSearchFieldsFromJson(content)
	} else {
		return p.ExtractElasticSearchFieldsFromPlain(content)
	}
}

func (p *PeckTask) ExtractInfluxDbFields(content string) {
	if p.Config.LogFormat == "json" {
		 p.ExtractInfluxDbFieldsFromPlain(content)
	} else {
		 p.ExtractInfluxDbFieldsFromPlain(content)
	}
}


func (p *PeckTask) Process(content string) {
	log.Infof("Process-----------")
	log.Infof("sender%v",p.sender)
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
