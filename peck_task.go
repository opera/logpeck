package logpeck

import ()

type Filter struct {
}

type PeckTask struct {
	Name       string
	FilterExpr string
	ESConfig   ElasticSearchConfig
	Running    bool
}

func NewPeckTask(c *PeckTaskConfig) (*PeckTask, error) {
	task := &PeckTask{
		Name:     c.Name,
		ESConfig: c.ESConfig,
		Running:  false,
	}
	return task, nil
}

func (p *PeckTask) Run() {

}

func (p *PeckTask) Pause() error {
	return nil
}

func (p *PeckTask) Cancel() error {
	return nil
}
