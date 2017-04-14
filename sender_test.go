package logpeck

import (
	"testing"
)

func TestSendToElasticSearch(*testing.T) {
	config := &ElasticSearchConfig{
		URL:   "http://127.0.0.1:9200",
		Index: "testes",
		Type:  "log",
	}
	InitElasticSearchMapping(config)
	//	panic(time.Now().Format("2006-01-02 15:04:05"))
	//url := "127.0.0.1:9200"
	//index := "test"
	//ty := "helloes"
	//content := "well done"
	//	SendToElasticSearch(url, index, ty, content)
}
