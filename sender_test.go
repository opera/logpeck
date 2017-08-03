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

func TestSendToElasticSearch(*testing.T) {
	//config := &ElasticSearchConfig{
	//	URL:   "http://127.0.0.1:9200",
	//	Index: "testes",
	//	Type:  "log",
	//}
	//InitElasticSearchMapping(config)
	//	panic(time.Now().Format("2006-01-02 15:04:05"))
	//url := "127.0.0.1:9200"
	//index := "test"
	//ty := "helloes"
	//content := "well done"
	//	SendToElasticSearch(url, index, ty, content)
}
