package logpeck

type PeckTaskConfig struct {
	Name     string
	LogPath  string
	ESConfig ElasticSearchConfig

	LogFormat  string
	FilterExpr string
	Fields     []PeckField
	Delimiters string
}

type PeckField struct {
	Name   string
	Value  string
	Type   string
	ESType string
}

type ElasticSearchConfig struct {
	Hosts []string
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
