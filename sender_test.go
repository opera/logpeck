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
			SenderName: "ElasticsearchConfig",
			Config:     ESConfig,
		}
		sender, err := NewSender(&config, nil)
		if err != nil {
			fmt.Printf("New sender error")
		}
		proto := "logpeck"
		Esender := sender.(*ElasticSearchSender)
		if proto != Esender.GetIndexName() {
			//panic(proto)
		}
	}

	{
		ESConfig := ElasticSearchConfig{
			Hosts: []string{"127.0.0.1:9200"},
			Index: "logpeck-%{+2006.01.02}",
			Type:  "hello",
		}
		config := SenderConfig{
			SenderName: "ElasticsearchConfig",
			Config:     ESConfig,
		}
		sender, err := NewSender(&config, nil)
		if err != nil {
			fmt.Printf("New sender error")
		}
		Esender := sender.(*ElasticSearchSender)
		indexName := Esender.GetIndexName()
		fmt.Printf("proto: %s, indexName: %s\n", config.Config.(ElasticSearchConfig).Index, indexName)
		if len(indexName) != 18 {
			panic(indexName)
		}
	}
}
