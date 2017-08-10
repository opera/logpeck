package logpeck

import (
	"fmt"
	"testing"
)

func TestGetIndexName(*testing.T) {
	{
		config := ElasticSearchConfig{
			Hosts: []string{"127.0.0.1:9200"},
			Index: "logpeck",
			Type:  "hello",
		}
		sender := NewElasticSearchSender(&config, nil)
		proto := "logpeck"
		if proto != sender.GetIndexName() {
			panic(proto)
		}
	}

	{
		config := ElasticSearchConfig{
			Hosts: []string{"127.0.0.1:9200"},
			Index: "logpeck-%{+2006.01.02}",
			Type:  "hello",
		}
		sender := NewElasticSearchSender(&config, nil)
		indexName := sender.GetIndexName()
		fmt.Printf("proto: %s, indexName: %s\n", config.Index, indexName)
		if len(indexName) != 18 {
			panic(indexName)
		}
	}
}
