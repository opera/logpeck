package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

type ExtractorConfig struct {
	Name   string
	Config interface{}
}

type Extractor interface {
	Extract(content string) (map[string]interface{}, error)
	Close()
}

func NewExtractorConfig(configStr string) (ExtractorConfig, error) {
	c := ExtractorConfig{}
	j, err := sjson.NewJson([]byte(configStr))
	name, err := j.Get("Name").String()
	cJ := j.Get("Config")
	if err != nil || name == "" {
		return c, errors.New("extractor unmarshal error")
	}
	jbyte, err := cJ.MarshalJSON()
	if err == nil {
		return c, err
	}
	switch name {
	case "lua":
		c.Config, err = NewLuaExtractorConfig(jbyte)
	case "json":
		c.Config, err = NewJsonExtractorConfig(jbyte)
	case "text":
		c.Config, err = NewTextExtractorConfig(jbyte)
	default:
		err = errors.New("extractor name error: " + c.Name)
	}
	log.Infof("[ExtractorConfig] Init finish %#v, %#v", c, err)
	return c, err
}

func NewExtractor(configStr string, fields []PeckField) (Extractor, error) {
	c, err := NewExtractorConfig(configStr)
	if err != nil {
		return nil, err
	}
	var e Extractor
	switch c.Name {
	case "lua":
		e, err = NewLuaExtractor(c)
	case "json":
		e, err = NewJsonExtractor(c, fields)
	case "text":
		e, err = NewTextExtractor(c, fields)
	default:
		err = errors.New("extractor name error: " + c.Name)
	}
	return e, nil
}
