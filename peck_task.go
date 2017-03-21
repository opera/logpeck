package logpeck

import (
	"fmt"
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

}

func (p *PeckTask) Run() {
}

func (p *PeckTask) Pause() error {
}

func (p *PeckTask) Cancel() error {
}
