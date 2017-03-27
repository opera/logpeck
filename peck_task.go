package logpeck

import ()

type Filter struct {
}

type PeckTask struct {
	Name       string
	FilterExpr string
	ESConfig   ElasticSearchConfig

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
