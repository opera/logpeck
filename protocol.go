package logpeck

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

type PeckTaskConfig struct {
	Name         string
	LogPath      string
	SenderConfig SenderConfig

	LogFormat  string
	Keywords   string
	Fields     []PeckField
	Delimiters string
	Test       TestModule
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

func ParseConfig(j *sjson.Json) (senderConfig SenderConfig, e error) {
	cJson := j.Get("SenderConfig")
	if cJson.Interface() == nil {
		return senderConfig, nil
	}
	senderConfig.SenderName, e = cJson.Get("SenderName").String()
	if e != nil {
		log.Infof("[ParseConfig]err: %v", e)
		return
	}
	if senderConfig.SenderName == "ElasticSearchConfig" {
		elasticSearchConfig := ElasticSearchConfig{}
		cJson := cJson.Get("Config")
		if cJson.Interface() == nil {
			return senderConfig, nil
		}

		jbyte, err := cJson.MarshalJSON()
		if err != nil {
			return
		}
		err = json.Unmarshal(jbyte, &elasticSearchConfig)
		if err != nil {
			return
		}
		log.Infof("[ParseConfig]ElasticSearchConfig: %v", elasticSearchConfig)
		senderConfig.Config = elasticSearchConfig
	}
	if senderConfig.SenderName == "InfluxDbConfig" {
		influxDbConfig := InfluxDbConfig{}
		cJson := cJson.Get("Config")
		if cJson.Interface() == nil {
			return senderConfig, nil
		}

		jbyte, err := cJson.MarshalJSON()
		if err != nil {
			return
		}
		err = json.Unmarshal(jbyte, &influxDbConfig)
		if err != nil {
			return
		}
		log.Infof("[ParseConfig]influxDbConfig: %v", influxDbConfig)
		senderConfig.Config = influxDbConfig
	}
	return senderConfig, nil
}

func (p *PeckTaskConfig) Unmarshal(jsonStr []byte) (e error) {
	j, je := sjson.NewJson(jsonStr)
	if je != nil {
		return je
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

	// Parse "SenderConfig", optional
	p.SenderConfig, e = ParseConfig(j)
	if e != nil {
		return e
	}

	// Parse "LogFormat", optional
	p.LogFormat, e = GetString(j, "LogFormat", false)
	if e != nil {
		return e
	}

	// Parse "FilterExpr", optional
	p.Keywords, e = GetString(j, "Keywords", false)
	if e != nil {
		return e
	}

	// Parse "Delimiters", optional
	p.Delimiters, e = GetString(j, "Delimiters", false)
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
			if val, ok := field.(map[string]interface{})["Value"]; ok {
				if f.Value, ok = val.(string); !ok {
					return errors.New("Fields format error: Value must be a string")
				}
			}
			p.Fields = append(p.Fields, f)
		}
	}

	return nil
}
