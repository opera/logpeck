package logpeck

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

// JsonExtractorConfig .
type JsonExtractorConfig struct {
	Fields []PeckField
}

// JsonExtractor .
type JsonExtractor struct {
	config *JsonExtractorConfig
	fields map[string]bool
}

// NewJsonExtractorConfig .
func NewJsonExtractorConfig(configStr []byte) (JsonExtractorConfig, error) {
	c := JsonExtractorConfig{}
	err := json.Unmarshal(configStr, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

// NewJsonExtractor .
func NewJsonExtractor(config interface{}) (JsonExtractor, error) {
	c, ok := config.(JsonExtractorConfig)
	if !ok {
		return JsonExtractor{}, errors.New("JsonExtractor config error")
	}
	e := JsonExtractor{
		config: &c,
		fields: make(map[string]bool),
	}
	for _, f := range c.Fields {
		e.fields[f.Name] = true
	}
	log.Infof("[JsonExtractor] Init extractor finished %#v", e)
	return e, nil
}

// Extract .
func (je JsonExtractor) Extract(content string) (map[string]interface{}, error) {
	fields := make(map[string]interface{})
	jContent, err := sjson.NewJson([]byte(content))
	if err != nil {
		return nil, err
	}
	mContent, err := jContent.Map()
	if err != nil {
		return nil, errors.New("Log is not json format")
	}
	if len(je.fields) == 0 {
		return map[string]interface{}{"_Log": content}, nil
	}
	for field := range je.fields {
		key := SplitString(field, ".")
		value := ""
		length := len(key)
		tmp := mContent
		for i := 0; i < length; i++ {
			if i == length-1 {
				if v, ok := tmp[key[i]].(string); ok {
					value = v
				} else if v, ok := tmp[key[i]].(json.Number); ok {
					value = v.String()
				} else {
					value = fmt.Sprintf("unknown type %v", tmp[key[i]])
				}
				break
			}
			tmp = tmp[key[i]].(map[string]interface{})
		}
		fields[field] = value
	}
	return fields, nil
}

// Close .
func (je JsonExtractor) Close() {
}
