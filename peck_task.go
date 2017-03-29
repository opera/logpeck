package logpeck

import ()

type PeckTask struct {
	Name     string
	Filter   string
	ESConfig ElasticSearchConfig

	pause bool
}

func NewPeckTask(c *PeckTaskConfig) *PeckTask {
	task := &PeckTask{
		Name:     c.Name,
		ESConfig: c.ESConfig,
		pause:    true,
	}
	return task
}

func (p *PeckTask) Start() {
	p.pause = false
}

func (p *PeckTask) Pause() {
	p.pause = true
}

func (p *PeckTask) Process(content string) {
	if p.pause {
		return
	}
	SendToElasticSearch(p.ESConfig.URL, p.ESConfig.Index, p.ESConfig.Type, content)
}
