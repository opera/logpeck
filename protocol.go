package logpeck

import (
	"encoding/json"
	"errors"
	"fmt"
	sjson "github.com/bitly/go-simplejson"
)

type PeckTaskConfig struct {
	Name             string
	LogPath          string
	ExtractorConfig  ExtractorConfig
	SenderConfig     SenderConfig
	AggregatorConfig AggregatorConfig

	Keywords string
	Fields   []PeckField
	Test     TestModule
}

type PeckField struct {
	Name  string
	Value string
}

type SenderConfig struct {
	SenderName string
	Config     interface{}
}

type PeckTaskStat struct {
	Name        string
	LogPath     string
	LinesPerSec int64
	BytesPerSec int64
	LinesTotal  int64
	BytesTotal  int64
	Stop        bool
}

type Stat struct {
	Name        string
	LinesPerSec int64
	BytesPerSec int64
	LinesTotal  int64
	BytesTotal  int64
}

type LogStat struct {
	LogPath         string
	PeckTaskConfigs []PeckTaskConfig
	PeckTaskStats   []PeckTaskStat
}

type PeckerStat struct {
	Name     string
	Stat     Stat
	LogStats []LogStat
}

type TestModule struct {
	TestNum int
	Timeout int
}

func GetString(j *sjson.Json, key string, required bool) (string, error) {
	valJson := j.Get(key)

	if valJson.Interface() == nil {
		if required {
			return "", errors.New("Parse error: need field " + key)
		} else {
			return "", nil
		}
	}
	return valJson.String()
}

func GetStringArray(j *sjson.Json, key string) ([]string, error) {
	valJson := j.Get(key)

	if valJson.Interface() == nil {
		return []string{""}, errors.New("Parse error: need field " + key)
	}
	return valJson.StringArray()
}

func GetMarshalString(j *sjson.Json, name string) (string, bool) {
	cJson := j.Get(name)
	if cJson.Interface() == nil {
		return "", false
	}
	jbyte, err := cJson.MarshalJSON()
	if err != nil {
		return "", false
	}
	return string(jbyte), true
}

func (p *PeckTaskConfig) Unmarshal(jsonStr []byte) (e error) {
	j, e := sjson.NewJson(jsonStr)
	if e != nil {
		return e
	}

	// Parse "Name", required
	p.Name, e = GetString(j, "Name", true)
	if e != nil {
		return e
	}

	// Parse "LogPath", optional
	p.LogPath, e = GetString(j, "LogPath", false)
	if e != nil {
		return e
	}

	// Parse "ExtractorConfig", optional
	eConfStr, ok := GetMarshalString(j, "ExtractorConfig")
	if ok {
		p.ExtractorConfig, e = NewExtractorConfig(eConfStr)
		if e != nil {
			return e
		}
	}

	// Parse "SenderConfig", optional
	p.SenderConfig, e = GetSenderConfig(j)
	if e != nil {
		return e
	}

	//Parse "aggregatorConfig", optional
	aggregatorConfig := j.Get("AggregatorConfig")
	jbyte, e := aggregatorConfig.MarshalJSON()
	if e != nil {
		return e
	}
	e = json.Unmarshal(jbyte, &p.AggregatorConfig)
	if e != nil {
		return e
	}

	// Parse "FilterExpr", optional
	p.Keywords, e = GetString(j, "Keywords", false)
	if e != nil {
		return e
	}

	testJ := j.Get("Test")
	if e != nil {
		p.Test.TestNum = 1
		p.Test.Timeout = 1
	}
	// Parse "TestNum", optional
	val, e := testJ.Get("TestNum").Int()
	if e != nil {
		p.Test.TestNum = 1
	}
	p.Test.TestNum = val

	// Parse "Time", optional
	time, e := testJ.Get("Timeout").Int()
	if e != nil {
		p.Test.Timeout = 1
	}
	p.Test.Timeout = time

	// Parse "Fields", optional
	if fields, e := j.Get("Fields").Array(); e == nil {
		fmt.Println(len(fields))
		for _, field := range fields {
			var f PeckField
			if name, ok := field.(map[string]interface{})["Name"]; ok {
				if f.Name, ok = name.(string); !ok {
					return errors.New("Fields format error: Name must be a string")
				}
			} else {
				return errors.New("Fields error: need Name")
			}
			p.Fields = append(p.Fields, f)
		}
	}

	return nil
}
