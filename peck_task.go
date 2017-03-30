package logpeck

import ()

type PeckTask struct {
	Name     string
	Filter   string
	ESConfig ElasticSearchConfig

	Stop bool
}

func NewPeckTask(c *PeckTaskConfig) *PeckTask {
	task := &PeckTask{
		Name:     c.Name,
		ESConfig: c.ESConfig,
		Stop:     true,
	}
	return task
}

func (p *PeckTask) Start() {
	p.Stop = false
}

func (p *PeckTask) Pause() {
	p.Stop = true
}

func (p *PeckTask) Process(content string) {
	if p.Stop {
		return
	}
	SendToElasticSearch(p.ESConfig.URL, p.ESConfig.Index, p.ESConfig.Type, content)
}
