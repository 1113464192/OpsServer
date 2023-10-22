package crontab

import (
	"github.com/robfig/cron"
)

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 5 * * *", CronMysqlLogRename)
	c.Start()
}
