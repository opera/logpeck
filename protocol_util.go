package logpeck

import (
	"encoding/json"
	"errors"

	sjson "github.com/bitly/go-simplejson"
)

// GetString from Json struct
func GetString(j *sjson.Json, key string, required bool) (string, error) {
	valJ := j.Get(key)

	if valJ.Interface() == nil {
		if required {
			return "", errors.New("Parse error: need field " + key)
		}
		return "", nil
	}
	return valJ.String()
}

// GetStringArray .
func GetStringArray(j *sjson.Json, key string) ([]string, error) {
	valJ := j.Get(key)

	if valJ.Interface() == nil {
		return []string{""}, errors.New("Parse error: need field " + key)
	}
	return valJ.StringArray()
}

// GetMarshalString .
func GetMarshalString(j *sjson.Json, name string) (string, bool) {
	cJ := j.Get(name)
	if cJ.Interface() == nil {
		return "", false
	}
	jbyte, err := cJ.MarshalJSON()
	if err != nil {
		return "", false
	}
	return string(jbyte), true
}

// Unmarshal .
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
	eConfStr, ok := GetMarshalString(j, "Extractor")
	if ok {
		p.Extractor, e = NewExtractorConfig(eConfStr)
		if e != nil {
			return e
		}
	}

	// Parse "SenderConfig", optional
	p.Sender, e = GetSenderConfig(j)
	if e != nil {
		return e
	}

	//Parse "aggregatorConfig", optional
	aggregatorConfig := j.Get("Aggregator")
	jbyte, e := aggregatorConfig.MarshalJSON()
	if e != nil {
		return e
	}
	e = json.Unmarshal(jbyte, &p.Aggregator)
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

	/*
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
				if value, ok := field.(map[string]interface{})["Value"]; ok {
					if f.Value, ok = value.(string); !ok {
						return errors.New("Fields format error: Name must be a string")
					}
				}else {
					return errors.New("Fields error: need Value")
				}
				p.Fields = append(p.Fields, f)
			}
		}*/

	return nil
}
