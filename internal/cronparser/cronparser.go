package cronparser

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type CronParser struct {
	parser cron.Parser
}

func NewParser() *CronParser {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	return &CronParser{parser: parser}
}

func (cp *CronParser) NextTime(crontab string) (time.Time, error) {
	schedule, err := cp.parser.Parse(crontab)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing cron expression: %w", err)
	}
	nextTime := schedule.Next(time.Now())
	return nextTime, nil
}
