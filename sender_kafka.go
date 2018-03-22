package logpeck

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	sjson "github.com/bitly/go-simplejson"
	"sync"
	"time"
)

type KafkaConfig struct {
	Brokers []string `json:"Brokers"`
	Topic   string   `json:"Topic"`

	MaxMessageBytes int                     `json:"MaxMessageBytes"`
	RequiredAcks    sarama.RequiredAcks     `json:"RequiredAcks"`
	Timeout         time.Duration           `json:"Timeout"`
	Compression     sarama.CompressionCodec `json:"Compression"`
	Partitioner     string                  `json:"Partitioner"`
	ReturnErrors    bool                    `json:"ReturnErrors"`
	Flush           KafkaFlush              `json:"Flush"`
	Retry           KafkaRetry              `json:"Retry"`
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
	mu            sync.Mutex
	lastIndexName string
	producer      sarama.SyncProducer
}

func NewKafkaSenderConfig(jbyte []byte) (KafkaConfig, error) {
	KafkaConfig := KafkaConfig{}
	err := json.Unmarshal(jbyte, &KafkaConfig)
	if err != nil {
		return KafkaConfig, err
	}
	log.Infof("[NewKafkaSenderConfig]ElasticSearchConfig: %v", KafkaConfig)
	return KafkaConfig, nil
}

func NewKafkaSender(senderConfig *SenderConfig) (*KafkaSender, error) {
	sender := KafkaSender{}
	config, ok := senderConfig.Config.(KafkaConfig)
	if !ok {
		return &sender, errors.New("New NewKafkaSender error ")
	}
	sender = KafkaSender{
		config: config,
	}
	return &sender, nil
}
func GetKafkaConfig(cJson *sjson.Json) (kafkaConfig KafkaConfig, e error) {
	kafkaConfig.Brokers, e = GetStringArray(cJson, "Brokers")
	if e != nil {
		return kafkaConfig, e
	}

	kafkaConfig.Topic, e = GetString(cJson, "Topic", true)
	if e != nil {
		return kafkaConfig, e
	}

	kafkaConfig.MaxMessageBytes, e = cJson.Get("MaxMessageBytes").Int()
	if e != nil {
		kafkaConfig.MaxMessageBytes = 1000000
	}

	kafkaJson := cJson.Get("RequiredAcks")
	if kafkaJson.Interface() == nil {
		kafkaConfig.RequiredAcks = 1
	} else {
		kafkaByte, e := kafkaJson.MarshalJSON()
		if e != nil {
			return kafkaConfig, e
		}
		e = json.Unmarshal(kafkaByte, &kafkaConfig.RequiredAcks)
		if e != nil {
			return kafkaConfig, e
		}
	}

	kafkaJson = cJson.Get("Timeout")
	if kafkaJson.Interface() == nil {
		kafkaConfig.Timeout = 10 * time.Second
	} else {
		kafkaByte, e := kafkaJson.MarshalJSON()
		if e != nil {
			return kafkaConfig, e
		}
		e = json.Unmarshal(kafkaByte, &kafkaConfig.Timeout)
		if e != nil {
			return kafkaConfig, e
		}
	}

	kafkaJson = cJson.Get("Compression")
	if kafkaJson.Interface() == nil {
		kafkaConfig.Compression = 0
	} else {
		kafkaByte, e := kafkaJson.MarshalJSON()
		if e != nil {
			return kafkaConfig, e
		}
		e = json.Unmarshal(kafkaByte, &kafkaConfig.Compression)
		if e != nil {
			return kafkaConfig, e
		}
	}

	kafkaConfig.Partitioner, e = GetString(cJson, "Partitioner", true)
	if e != nil {
		kafkaConfig.Partitioner = "RandomPartitioner"
	}

	kafkaJson = cJson.Get("ReturnErrors")
	if kafkaJson.Interface() == nil {
		kafkaConfig.ReturnErrors = true
	} else {
		kafkaByte, e := kafkaJson.MarshalJSON()
		if e != nil {
			return kafkaConfig, e
		}
		e = json.Unmarshal(kafkaByte, &kafkaConfig.ReturnErrors)
		if e != nil {
			return kafkaConfig, e
		}
	}

	kafkaJson = cJson.Get("Flush")
	if kafkaJson.Interface() == nil {

	} else {
		kafkaByte, e := kafkaJson.MarshalJSON()
		if e != nil {
			return kafkaConfig, e
		}
		e = json.Unmarshal(kafkaByte, &kafkaConfig.Flush)
		if e != nil {
			return kafkaConfig, e
		}
	}

	kafkaJson = cJson.Get("Retry")
	if kafkaJson.Interface() == nil {
		kafkaConfig.Retry.Max = 3
		kafkaConfig.Retry.Backoff = 100 * time.Millisecond
	} else {
		kafkaByte, e := kafkaJson.MarshalJSON()
		if e != nil {
			return kafkaConfig, e
		}
		e = json.Unmarshal(kafkaByte, &kafkaConfig.Retry)
		if e != nil {
			return kafkaConfig, e
		}
	}

	return kafkaConfig, nil
}

func (p *KafkaSender) Start() error {
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
		log.Debug("[Start]Partitioner：%v is Invalid", p.config.Partitioner)
	}

	producer, err := sarama.NewSyncProducer(p.config.Brokers, config)
	if err != nil {
		log.Error("[Start] producer err:%v", err)
		return err
	}
	p.producer = producer
	return nil
}

func (p *KafkaSender) Stop() error {
	if p.producer == nil {
		return nil
	} else if err := p.producer.Close(); err != nil {
		return err
	}
	return nil
}

func (p *KafkaSender) Send(fields map[string]interface{}) {
	msg := &sarama.ProducerMessage{
		Topic:     p.config.Topic,
		Partition: int32(-1),
		Key:       sarama.StringEncoder("key"),
	}
	value, err := json.Marshal(fields)
	if err != nil {
		log.Error("[Send] fields Marshal err:%v", err)
		return
	}
	msg.Value = sarama.ByteEncoder(value)
	defer func(){
		if err:=recover();err!=nil{
			log.Info("[KafkaSender]error:%v",err)
		}
	}()
	paritition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Error("Send Message Fail")
	}

	log.Debug("[Send]Partion = %d, offset = %d, value = %v \n", paritition, offset, fields)
	//p.measurments.MeasurmentRecall(fields)
}
