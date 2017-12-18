package logpeck

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"testing"
)

func TestStartSend(*testing.T) {
	log.Infof("[aggregator_test] TestStartSend")

	interval := int64(30)
	name := "getTest"
	aggregators := map[string]AggregatorConfig{}
	test := AggregatorConfig{
		Tags:         []string{"cost"},
		Aggregations: []string{"cnt"},
		Target:       "cost",
		Time:         "time",
	}
	aggregators["test"] = test
	aggregator := NewAggregator(interval, name, &aggregators)

	start, _ := aggregator.StartSend(int64(29))
	if start == true {
		panic(aggregator)
	}
	start, _ = aggregator.StartSend(int64(31))
	if start == false {
		panic(aggregator)
	}
}

func TestRecord(*testing.T) {
	interval := int64(30)
	name := "test"
	aggregators := map[string]AggregatorConfig{}
	test := AggregatorConfig{
		Tags:         []string{"cost"},
		Aggregations: []string{"cnt"},
		Target:       "cost",
		Time:         "time",
	}
	aggregators["getTest"] = test
	aggregator := NewAggregator(interval, name, &aggregators)

	fields := make(map[string]interface{})
	fields["test"] = "getTest"
	fields["cost"] = "2"
	fields["time"] = "15"
	if aggregator.Record(fields) != int64(15) {
		panic(fields)
	}
	fmt.Print("[aggregator_test] TestRecord: buckets[getTest]= %v", aggregator.buckets["getTest"])
}

func TestDump(*testing.T) {
	interval := int64(30)
	name := "test"
	aggregators := map[string]AggregatorConfig{}
	test := AggregatorConfig{
		Tags:         []string{"cost"},
		Aggregations: []string{"cnt"},
		Target:       "cost",
		Time:         "time",
	}
	aggregators["getTest"] = test
	aggregator := NewAggregator(interval, name, &aggregators)

	fields := make(map[string]interface{})
	fields["test"] = "getTest"
	fields["cost"] = "2"
	fields["time"] = "15"
	if aggregator.Record(fields) != int64(15) {
		panic(fields)
	}
	dump := aggregator.Dump(int64(1))
	if dump["getTest"] != "getTest,cost=2 cnt=1 1" {
		panic(dump)
	}
}
