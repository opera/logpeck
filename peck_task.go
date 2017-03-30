package logpeck

import (
	"log"
)

type PeckTask struct {
	Config PeckTaskConfig
	Stat   PeckTaskStat
}

func NewPeckTask(c *PeckTaskConfig, s *PeckTaskStat) *PeckTask {
	if s == nil {
		return &PeckTask{
			Config: *c,
			Stat: PeckTaskStat{
				Name:    c.Name,
				LogPath: c.LogPath,
				Stop:    true,
			},
		}
	} else {
		if c.LogPath != s.LogPath || c.Name != s.Name {
			log.Fatalf("Config[%s], Stat[%s]", c, s)
		}
		return &PeckTask{
			Config: *c,
			Stat:   *s,
		}
	}
}

func (p *PeckTask) Start() {
	log.Printf("[PeckTask] Start")
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
		log.Println("PeckTask stopped" + content)
		return
	}
	SendToElasticSearch(p.Config.ESConfig.URL, p.Config.ESConfig.Index, p.Config.ESConfig.Type, content)
}
