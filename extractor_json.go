package logpeck

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

// JSONExtractorConfig .
type JSONExtractorConfig struct {
	Fields []PeckField
}

// JSONExtractor .
type JSONExtractor struct {
	config *JSONExtractorConfig
	fields map[string]bool
}

// NewJSONExtractorConfig .
func NewJSONExtractorConfig(configStr []byte) (JSONExtractorConfig, error) {
	c := JSONExtractorConfig{}
	err := json.Unmarshal(configStr, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

// NewJSONExtractor .
func NewJSONExtractor(config interface{}) (JSONExtractor, error) {
	c, ok := config.(JSONExtractorConfig)
	if !ok {
		return JSONExtractor{}, errors.New("JSONExtractor config error")
	}
	e := JSONExtractor{
		config: &c,
		fields: make(map[string]bool),
	}
	for _, f := range c.Fields {
		e.fields[f.Name] = true
	}
	log.Infof("[JSONExtractor] Init extractor finished %#v", e)
	return e, nil
}

// Extract .
func (je JSONExtractor) Extract(content string) (map[string]interface{}, error) {
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
func (je JSONExtractor) Close() {
}
