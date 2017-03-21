package logpeck

import (
//	"fmt"
)

type Filter struct {
}

type PeckTask struct {
	Name     string
	LogPath  string
	Filter   Filter
	ESConfig ElasticSearchConfig
}

func (p *PeckTask) Init(c *PeckTaskConfig) error {
	return nil
}

func (p *PeckTask) Run() {

}

func (p *PeckTask) Pause() error {
	return nil
}

func (p *PeckTask) Cancel() error {
	return nil
}
