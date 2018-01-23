package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
)

type Sender interface {
	Send(map[string]interface{})
	Start() error
	Stop() error
}

func GetSenderConfig(j *sjson.Json) (senderConfig SenderConfig, err error) {
	cJson := j.Get("SenderConfig")
	if cJson.Interface() == nil {
		return senderConfig, nil
	}
	senderConfig.SenderName, err = cJson.Get("SenderName").String()
	if err != nil {
		log.Infof("[ParseConfig]err: %v", err)
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

	switch senderConfig.SenderName {
	case "ElasticsearchConfig":
		senderConfig.Config, err = NewElasticSearchSenderConfig(jbyte)
	case "InfluxDbConfig":
		senderConfig.Config, err = NewInfluxDbSenderConfig(jbyte)
	case "KafkaConfig":
		senderConfig.Config, err = NewKafkaSenderConfig(jbyte)
	default:
		err = errors.New("[GetSenderConfig]sender name error: " + senderConfig.SenderName)
	}

	return senderConfig, err
}

func NewSender(senderConfig *SenderConfig) (sender Sender, err error) {
	switch senderConfig.SenderName {
	case "ElasticsearchConfig":
		sender, err = NewElasticSearchSender(senderConfig)
	case "InfluxDbConfig":
		sender, err = NewInfluxDbSender(senderConfig)
	case "KafkaConfig":
		sender, err = NewKafkaSender(senderConfig)
	default:
		err = errors.New("[NewSender]sender name error: " + senderConfig.SenderName)
	}
	return sender, err
}
