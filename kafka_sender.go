package logpeck

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

type KafkaConfig struct {
	Hosts []string `json:"Hosts"`
	Topic string   `json:"Topic"`

	MaxMessageBytes int                     `json:"MaxMessageBytes"`
	RequiredAcks    sarama.RequiredAcks     `json:"RequiredAcks"`
	Timeout         time.Duration           `json:"Timeout"`
	Compression     sarama.CompressionCodec `json:"Compression"`
	Partitioner     string                  `json:"Partitioner"`
	ReturnErrors    bool                    `json:"ReturnErrors"`
	Flush           KafkaFlush              `json:"Flush"`
	Retry           KafkaRetry              `json:"Retry"`

	Interval          int64              `json:"Interval"`
	AggregatorConfigs []AggregatorConfig `json:"AggregatorConfigs"`
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
	config := sarama.NewConfig()

	config.Producer.MaxMessageBytes = p.config.MaxMessageBytes
	config.Producer.RequiredAcks = p.config.RequiredAcks
	config.Producer.Timeout = p.config.Timeout
	config.Producer.Compression = p.config.Compression
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = p.config.ReturnErrors
	config.Producer.Flush.Bytes = p.config.Flush.Bytes
	config.Producer.Flush.Frequency = p.config.Flush.Frequency
	config.Producer.Flush.MaxMessages = p.config.Flush.MaxMessages
	config.Producer.Flush.Messages = p.config.Flush.Messages
	config.Producer.Retry.Backoff = p.config.Retry.Backoff
	config.Producer.Retry.Max = p.config.Retry.Max
	switch p.config.Partitioner {
	case "RandomPartitioner":
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	case "HashPartitioner":
		config.Producer.Partitioner = sarama.NewHashPartitioner
	case "ManualPartitioner":
		config.Producer.Partitioner = sarama.NewManualPartitioner
	case "RoundRobinPartitioner":
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	default:
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		log.Debug("[sender]Partitioner：%v is Invalid", p.config.Partitioner)
	}

	producer, err := sarama.NewSyncProducer(p.config.Hosts, config)
	if err != nil {
		log.Infof("[Send] producer err:%v", err)
		return
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic:     p.config.Topic,
		Partition: int32(-1),
		Key:       sarama.StringEncoder("key"),
	}
	value, err := json.Marshal(fields)
	if err != nil {
		log.Infof("[Send] fields Marshal err:%v", err)
		return
	}
	msg.Value = sarama.ByteEncoder(value)
	paritition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Infof("Send Message Fail")
	}

	log.Infof("[Send]Partion = %d, offset = %d, value = %v \n", paritition, offset, fields)
	//p.measurments.MeasurmentRecall(fields)
}
