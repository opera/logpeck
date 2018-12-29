package logpeck

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

const (
	exTypeLua  = "lua"
	exTypeJSON = "json"
	exTypeText = "text"
)

// Extractor .
type Extractor interface {
	Extract(content string) (map[string]interface{}, error)
	Close()
}

// NewExtractorConfig .
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
	switch strings.ToLower(name) {
	case exTypeLua:
		c.Config, err = NewLuaExtractorConfig(jbyte)
	case exTypeJSON:
		c.Config, err = NewJSONExtractorConfig(jbyte)
	case exTypeText:
		c.Config, err = NewTextExtractorConfig(jbyte)
	default:
		err = errors.New("extractor name error: " + c.Name)
	}
	c.Name = name
	log.Infof("[ExtractorConfig] Init finish %#v, %#v", c, err)
	return c, err
}

// NewExtractor .
func NewExtractor(c ExtractorConfig) (e Extractor, err error) {
	switch strings.ToLower(c.Name) {
	case exTypeLua:
		e, err = NewLuaExtractor(c.Config)
	case exTypeJSON:
		e, err = NewJSONExtractor(c.Config)
	case exTypeText:
		e, err = NewTextExtractor(c.Config)
	default:
		err = errors.New("extractor name error: " + c.Name)
	}
	return e, err
}
