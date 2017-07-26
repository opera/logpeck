package logpeck

import (
	JsonUtil "github.com/bitly/go-simplejson"
)

type PeckTaskConfig struct {
	Name      string
	LogPath   string
	ESConfig  ElasticSearchConfig
	Extractor Extractor

	// Deprecated
	FilterExpr string
	Fields     []PeckField
	Delimiters string
}

type Extractor struct {
	Name string
}

type PlainExtractor struct {
	Extractor Extractor

	FilterExpr string
	Fields     []PeckField
	Delimiters string
}

type JsonExtractor struct {
	Extractor Extractor

	FilterExpr string
	Fields     []PeckField
}

type PeckField struct {
	Name   string
	Value  string
	Type   string
	ESType string
}

type ElasticSearchConfig struct {
	URL   string
	Index string
	Type  string
}

type PeckTaskStat struct {
	Name        string
	LogPath     string
	LinesPerSec int64
	BytesPerSec int64
	LinesTotal  int64
	BytesTotal  int64
	Stop        bool
}

type Stat struct {
	Name        string
	LinesPerSec int64
	BytesPerSec int64
	LinesTotal  int64
	BytesTotal  int64
}

type LogStat struct {
	LogPath         string
	PeckTaskConfigs []PeckTaskConfig
	PeckTaskStats   []PeckTaskStat
}

type PeckerStat struct {
	Name     string
	Stat     Stat
	LogStats []LogStat
}
