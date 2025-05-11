package crontotime

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type CronToTime struct {
	parser cron.Parser
}

func NewParser() *CronToTime {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	return &CronToTime{parser: parser}
}

func (ctd *CronToTime) NextTime(crontab string) (time.Time, error) {
	schedule, err := ctd.parser.Parse(crontab)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing cron expression: %w", err)
	}
	nextTime := schedule.Next(time.Now())
	return nextTime, nil
}
