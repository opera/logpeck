package logpeck

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

const (
	senderTypeES       = "elasticsearch"
	senderTypeKafka    = "kafka"
	senderTypeInfluxDb = "influxdb"
)

// Sender .
type Sender interface {
	Send(map[string]interface{})
	Start() error
	Stop() error
}

// GetSenderConfig .
func GetSenderConfig(j *sjson.Json) (senderConfig SenderConfig, err error) {
	cJson := j.Get("Sender")
	if cJson.Interface() == nil {
		return senderConfig, nil
	}
	senderConfig.Name, err = cJson.Get("Name").String()
	if err != nil {
		log.Infof("[GetSenderConfig]err: %v", err)
		return senderConfig, err
	}
	cJson = cJson.Get("Config")
	if cJson.Interface() == nil {
		return senderConfig, nil
	}
	jbyte, err := cJson.MarshalJSON()
	if err != nil {
		return senderConfig, err
	}

	switch strings.ToLower(senderConfig.Name) {
	case senderTypeES:
		senderConfig.Config, err = NewElasticSearchSenderConfig(jbyte)
	case senderTypeInfluxDb:
		senderConfig.Config, err = NewInfluxDbSenderConfig(jbyte)
	case senderTypeKafka:
		senderConfig.Config, err = NewKafkaSenderConfig(jbyte)
	default:
		err = errors.New("[GetSenderConfig]sender name error: " + senderConfig.Name)
	}

	return senderConfig, err
}

// NewSender .
func NewSender(senderConfig *SenderConfig) (sender Sender, err error) {
	switch strings.ToLower(senderConfig.Name) {
	case senderTypeES:
		sender, err = NewElasticSearchSender(senderConfig)
	case senderTypeInfluxDb:
		sender, err = NewInfluxDbSender(senderConfig)
	case senderTypeKafka:
		sender, err = NewKafkaSender(senderConfig)
	default:
		err = errors.New("[NewSender]sender name error: " + senderConfig.Name)
	}
	return sender, err
}
