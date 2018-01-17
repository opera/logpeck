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

func NewExtractorConfig(configStr string) (*ExtractorConfig, error) {
	c := ExtractorConfig{}
	j, err := sjson.NewJson([]byte(configStr))
	name, err := j.Get("Name").String()
	cJ := j.Get("Config")
	if err != nil || name == "" {
		return nil, errors.New("extractor unmarshal error")
	}
	jbyte, err := cJ.MarshalJSON()
	if err == nil {
		return nil, err
	}
	switch name {
	case "lua":
		c.Config, err = NewLuaExtractorConfig(jbyte)
	case "json":
		err = errors.New("not support json")
	case "text":
		err = errors.New("not support text")
	default:
		err = errors.New("extractor name error: " + c.Name)
	}
	log.Infof("[ExtractorConfig] Init finish %#v, %#v", c, err)
	return &c, err
}

func NewExtractor(configStr string) (*Extractor, error) {
	return nil, nil
}
