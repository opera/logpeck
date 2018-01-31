package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

const (
	ExTypeLua  = "Lua"
	ExTypeJson = "Json"
	ExTypeText = "Text"
)

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
		return c, nil
	}
	jbyte, err := cJ.MarshalJSON()
	if err != nil {
		return c, err
	}
	switch name {
	case ExTypeLua:
		c.Config, err = NewLuaExtractorConfig(jbyte)
	case ExTypeJson:
		c.Config, err = NewJsonExtractorConfig(jbyte)
	case ExTypeText:
		c.Config, err = NewTextExtractorConfig(jbyte)
	default:
		err = errors.New("extractor name error: " + c.Name)
	}
	c.Name = name
	log.Infof("[ExtractorConfig] Init finish %#v, %#v", c, err)
	return c, err
}

func NewExtractor(c ExtractorConfig) (e Extractor, err error) {
	switch c.Name {
	case ExTypeLua:
		e, err = NewLuaExtractor(c.Config)
	case ExTypeJson:
		e, err = NewJsonExtractor(c.Config)
	case ExTypeText:
		e, err = NewTextExtractor(c.Config)
	default:
		err = errors.New("extractor name error: " + c.Name)
	}
	return e, err
}
