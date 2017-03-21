package logpeck

import (
	"fmt"
	"os"
)

type Filter struct {
}

type PeckTask struct {
	Name     string
	LogPath  string
	Filter   Filter
	ESConfig ElasticSearchConfig

	file *os.File
}

func NewPeckTask(c *PeckTaskConfig) (*PeckTask, error) {
	f, f_err := os.Open(c.LogPath)
	if f_err != nil {
		return nil, f_err
	}
	fmt.Println(f)
	return &PeckTask{}, nil
}

func (p *PeckTask) Run() {

}

func (p *PeckTask) Pause() error {
	return nil
}

func (p *PeckTask) Cancel() error {
	return nil
}
