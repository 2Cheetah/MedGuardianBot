package crontodate

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type CronToDate struct {
	parser cron.Parser
}

func NewParser() *CronToDate {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	return &CronToDate{parser: parser}
}

func (ctd *CronToDate) NextTime(crontab string) (time.Time, error) {
	schedule, err := ctd.parser.Parse(crontab)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing cron expression: %w", err)
	}
	nextTime := schedule.Next(time.Now())
	return nextTime, nil
}
