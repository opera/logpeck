package logpeck

import (
	"fmt"
	"testing"
)

func TestGetIndexName(*testing.T) {
	{
		ESConfig := ElasticSearchConfig{
			Hosts: []string{"127.0.0.1:9200"},
			Index: "logpeck-%{+2006.01.02}",
			Type:  "hello",
		}
		config := SenderConfig{
			Name:   "ElasticsearchConfig",
			Config: ESConfig,
		}

		sender := NewElasticSearchSender(&config, nil)
		proto := "logpeck"
		if proto != sender.GetIndexName() {
			panic(proto)
		}
	}

	{
		ESConfig := ElasticSearchConfig{
			Hosts: []string{"127.0.0.1:9200"},
			Index: "logpeck-%{+2006.01.02}",
			Type:  "hello",
		}
		config := SenderConfig{
			Name:   "ElasticsearchConfig",
			Config: ESConfig,
		}
		sender := NewElasticSearchSender(&config, nil)
		indexName := sender.GetIndexName()
		fmt.Printf("proto: %s, indexName: %s\n", config.Config.(ElasticSearchConfig).Index, indexName)
		if len(indexName) != 18 {
			panic(indexName)
		}
	}
}
