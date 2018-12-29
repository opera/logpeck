package logpeck

// PeckTaskConfig http request protocol
type PeckTaskConfig struct {
	Name       string
	LogPath    string
	Extractor  ExtractorConfig
	Sender     SenderConfig
	Aggregator AggregatorConfig

	Keywords string
	Test     TestModule
}

// ExtractorConfig .
type ExtractorConfig struct {
	Name   string
	Config interface{}
}

// SenderConfig .
type SenderConfig struct {
	Name   string
	Config interface{}
}

// PeckField config
type PeckField struct {
	Name  string
	Value string
}

// PeckTaskStat .
type PeckTaskStat struct {
	Name        string
	LinesPerSec int64
	BytesPerSec int64
	LinesTotal  int64
	BytesTotal  int64
	Stop        bool
}

// Stat .
type Stat struct {
	Name        string
	LinesPerSec int64
	BytesPerSec int64
	LinesTotal  int64
	BytesTotal  int64
}

// LogStat .
type LogStat struct {
	LogPath         string
	PeckTaskConfigs []PeckTaskConfig
	PeckTaskStats   []PeckTaskStat
}

// PeckerStat .
type PeckerStat struct {
	Name     string
	Stat     Stat
	LogStats []LogStat
}

// TestModule .
type TestModule struct {
	TestNum int
	Timeout int
}
