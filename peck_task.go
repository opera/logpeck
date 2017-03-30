package logpeck

import ()

type PeckTask struct {
	config PeckTaskConfig
	stat   PeckTaskStat
}

func NewPeckTask(c *PeckTaskConfig) *PeckTask {
	task := &PeckTask{
		PeckTaskConfig{
			Name:     c.Name,
			LogPath:  c.LogPath,
			ESConfig: c.ESConfig,
		},
		PeckTaskStat{
			Name:    c.Name,
			LogPath: c.LogPath,
			Stop:    true,
		},
	}
	return task
}

func (p *PeckTask) Start() {
	p.stat.Stop = false
}

func (p *PeckTask) Stop() {
	p.stat.Stop = true
}

func (p *PeckTask) IsStop() bool {
	return p.stat.Stop
}

func (p *PeckTask) Process(content string) {
	if p.stat.Stop {
		return
	}
	SendToElasticSearch(p.config.ESConfig.URL, p.config.ESConfig.Index, p.config.ESConfig.Type, content)
}
