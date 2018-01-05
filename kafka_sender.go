package logpeck

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

type KafkaConfig struct {
	Hosts []string `json:"Hosts"`
	Topic string   `json:"Topic"`

	MaxMessageBytes int                           `json:"MaxMessageBytes"`
	RequiredAcks    sarama.RequiredAcks           `json:"RequiredAcks"`
	Timeout         time.Duration                 `json:"Timeout"`
	Compression     sarama.CompressionCodec       `json:"Compression"`
	Partitioner     sarama.PartitionerConstructor `json:"Partitioner"`
	Return          KafkaReturn                   `json:"Return"`
	Flush           KafkaFlush                    `json:"Flush"`
	Retry           KafkaRetry                    `json:"Retry"`

	Interval          int64              `json:"Interval"`
	AggregatorConfigs []AggregatorConfig `json:"AggregatorConfigs"`
}

type KafkaReturn struct {
	Successes bool `json:"ReturnSuccesses"`
	Errors    bool `json:"ReturnErrors"`
}

type KafkaFlush struct {
	Bytes       int           `json:"FlushBytes"`
	Messages    int           `json:"FlushMessages"`
	Frequency   time.Duration `json:"FlushFrequency"`
	MaxMessages int           `json:"FlushMaxMessages"`
}

type KafkaRetry struct {
	Max     int           `json:"RetryMax"`
	Backoff time.Duration `json:"RetryBackoff"`
}

type KafkaSender struct {
	config        KafkaConfig
	fields        []PeckField
	mu            sync.Mutex
	lastIndexName string
}

func NewKafkaSender(senderConfig *SenderConfig, fields []PeckField) *KafkaSender {
	config := senderConfig.Config.(KafkaConfig)
	sender := KafkaSender{
		config: config,
		fields: fields,
	}
	return &sender
}
func (p *KafkaSender) Send(fields map[string]interface{}) {
	log.Infof("[KafkaSender.send]%v", fields)
	config := sarama.NewConfig()
	log.Infof("[Kafka.Send] config=%v", config)
	/*
		config.Producer.MaxMessageBytes = p.config.MaxMessageBytes
		config.Producer.RequiredAcks = p.config.RequiredAcks
		config.Producer.Timeout = p.config.Timeout
		config.Producer.Compression = p.config.Compression
		config.Producer.Partitioner = p.config.Partitioner
		config.Producer.Return = p.config.Return
		config.Producer.Flush= p.config.Flush
		config.Producer.Retry = p.config.Retry
	*/
	producer, err := sarama.NewSyncProducer(p.config.Hosts, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	log.Infof("%v", fields)
	msg := &sarama.ProducerMessage{
		Topic:     p.config.Topic,
		Partition: int32(-1),
		Key:       sarama.StringEncoder("key"),
	}
	value, err := json.Marshal(fields)
	if err != nil {
		panic(err)
	}
	msg.Value = sarama.ByteEncoder(value)
	paritition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Send Message Fail")
	}

	log.Infof("Partion = %d, offset = %d, value = %v \n", paritition, offset, fields)
	//p.measurments.MeasurmentRecall(fields)
}
