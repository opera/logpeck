package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
	"strings"
)

const (
	SenderTypeES       = "elasticsearch"
	SenderTypeKafka    = "kafka"
	SenderTypeInfluxDb = "influxdb"
)

type Sender interface {
	Send(map[string]interface{})
	Start() error
	Stop() error
}

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
	case SenderTypeES:
		senderConfig.Config, err = NewElasticSearchSenderConfig(jbyte)
	case SenderTypeInfluxDb:
		senderConfig.Config, err = NewInfluxDbSenderConfig(jbyte)
	case SenderTypeKafka:
		senderConfig.Config, err = NewKafkaSenderConfig(jbyte)
	default:
		err = errors.New("[GetSenderConfig]sender name error: " + senderConfig.Name)
	}

	return senderConfig, err
}

func NewSender(senderConfig *SenderConfig) (sender Sender, err error) {
	switch strings.ToLower(senderConfig.Name) {
	case SenderTypeES:
		sender, err = NewElasticSearchSender(senderConfig)
	case SenderTypeInfluxDb:
		sender, err = NewInfluxDbSender(senderConfig)
	case SenderTypeKafka:
		sender, err = NewKafkaSender(senderConfig)
	default:
		err = errors.New("[NewSender]sender name error: " + senderConfig.Name)
	}
	return sender, err
}
