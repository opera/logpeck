package logpeck

import (
	"encoding/json"
	"errors"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

// TextExtractorConfig .
type TextExtractorConfig struct {
	Delimiters string
	Fields     []PeckField
}

// TextExtractor .
type TextExtractor struct {
	config *TextExtractorConfig
	fields map[string]int
}

// NewTextExtractorConfig .
func NewTextExtractorConfig(configStr []byte) (TextExtractorConfig, error) {
	c := TextExtractorConfig{}
	err := json.Unmarshal(configStr, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

// NewTextExtractor .
func NewTextExtractor(config interface{}) (TextExtractor, error) {
	c, ok := config.(TextExtractorConfig)
	e := TextExtractor{
		config: &c,
		fields: make(map[string]int),
	}
	if !ok {
		return e, errors.New("TextExtractor config error")
	}
	log.Info(c.Fields)
	for _, f := range c.Fields {
		if f.Value[0] != '$' {
			return e, errors.New("field format error: " + f.Value)
		}
		pos, err := strconv.Atoi(f.Value[1:])
		if err != nil {
			return e, errors.New("field format error: " + f.Value)
		}
		e.fields[f.Name] = pos
	}
	log.Infof("[TextExtractor] Init extractor finished %#v", e)
	return e, nil
}

// Extract .
func (te TextExtractor) Extract(content string) (map[string]interface{}, error) {
	if len(te.fields) == 0 {
		return map[string]interface{}{"_Log": content}, nil
	}
	fields := make(map[string]interface{})
	arr := SplitString(content, te.config.Delimiters)
	for k, v := range te.fields {
		if len(arr) < v {
			continue
		}
		fields[k] = arr[v-1]
	}
	return fields, nil
}

// Close .
func (te TextExtractor) Close() {
}
