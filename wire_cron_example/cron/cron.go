package cron

import (
	"github.com/robfig/cron/v3"

	"asong.cloud/Golang_Dream/wire_cron_example/cron/task"
)

type Cron struct {
	Scanner *task.Scanner
	Schedule *cron.Cron
}

func NewCron(scanner *task.Scanner) *Cron {
	return &Cron{
		Scanner: scanner,
		Schedule: cron.New(cron.WithSeconds()),
	}
}

func (s *Cron)Start()  error{
	_,err := s.Schedule.AddJob("*/1 * * * *",s.Scanner)
	if err != nil{
		return err
	}
	s.Schedule.Start()
	return nil
}