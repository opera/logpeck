package logpeck

import (
	"fmt"
	"testing"
)

func TestGetIndexName(*testing.T) {
	{
		proto := "logpeck"
		if proto != GetIndexName(proto) {
			panic(proto)
		}
	}

	{
		proto := "logpeck-%{+2006.01.02}"
		indexName := GetIndexName(proto)
		fmt.Printf("proto: %s, indexName: %s\n", proto, indexName)
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
