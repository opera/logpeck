package logpeck

import ()

type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat
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
	p.Stat.Stop = false
}

func (p *PeckTask) Stop() {
	p.Stat.Stop = true
}

func (p *PeckTask) IsStop() bool {
	return p.Stat.Stop
}

func (p *PeckTask) Process(content string) {
	if p.Stat.Stop {
		return
	}
	SendToElasticSearch(p.Config.ESConfig.URL, p.Config.ESConfig.Index, p.Config.ESConfig.Type, content)
}
