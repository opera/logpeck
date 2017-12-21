package logpeck

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
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
		Timestamp:    "time",
	}
	aggregators["test"] = test
	aggregator := NewAggregator(interval, name, &aggregators)

	deadline := aggregator.IsDeadline(int64(29))
	if deadline == true {
		panic(aggregator)
	}
	deadline = aggregator.IsDeadline(int64(31))
	if deadline == false {
		panic(aggregator)
	}
}

func TestRecord(*testing.T) {
	interval := int64(30)
	name := "test"
	aggregators := map[string]AggregatorConfig{}
	test := AggregatorConfig{
		Tags:         []string{"upstream"},
		Aggregations: []string{"cnt,avg"},
		Target:       "cost",
		Timestamp:    "time",
	}
	aggregators["getTest"] = test
	aggregator := NewAggregator(interval, name, &aggregators)

	fields := make(map[string]interface{})
	fields["test"] = "getTest"
	fields["upstream"] = "127.0.0.1"
	fields["cost"] = "2"
	fields["time"] = "15"
	if aggregator.Record(fields) != int64(15) {
		panic(fields)
	}
	if aggregator.buckets["getTest"][",upstream=127.0.0.1"][0] != 2 {
		panic(aggregator)
	}
	if aggregator.Record(fields) != int64(15) {
		panic(fields)
	}
	if aggregator.buckets["getTest"][",upstream=127.0.0.1"][0]+aggregator.buckets["getTest"][",upstream=127.0.0.1"][1] != 4 {
		panic(aggregator)
	}
}

func TestDump(*testing.T) {
	interval := int64(30)
	name := "test"
	aggregators := map[string]AggregatorConfig{}
	test := AggregatorConfig{
		Tags:         []string{"upstream"},
		Aggregations: []string{"cnt", "avg", "p99", "p50"},
		Target:       "cost",
		Timestamp:    "time",
	}
	aggregators["getTest"] = test
	aggregator := NewAggregator(interval, name, &aggregators)

	fields := make(map[string]interface{})
	fields["test"] = "getTest"
	fields["upstream"] = "127.0.0.1"
	fields["cost"] = "0"
	fields["time"] = "15"
	for i := 9; i >= 0; i-- {
		fields["cost"] = strconv.Itoa(i)
		aggregator.Record(fields)
	}
	dump := aggregator.Dump(int64(30))
	a := dump["getTest,upstream=127.0.0.1"].(map[string]int64)
	if a["cnt"] != 10 {
		panic(a)
	}
	if a["avg"] != 4 {
		panic(a)
	}
	if a["p99"] != 8 {
		panic(a)
	}
	if a["p50"] != 4 {
		panic(a)
	}
	if dump["timestamp"].(int64) != 30 {
		panic(dump)
	}
}
