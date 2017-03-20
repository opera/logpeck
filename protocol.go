package logpeck

type ElasticSearchConfig struct {
	URL   string
	Index string
	Type  string
}

type PeckTaskConfig struct {
	Name       string
	LogPath    string
	FilterExpr string
	ESConfig   ElasticSearchConfig
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
	PeckTaskStats   []Stat
}

type PeckerStat struct {
	Name     string
	Stat     Stat
	LogStats []LogStat
}
