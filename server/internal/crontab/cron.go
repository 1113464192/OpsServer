package crontab

import (
	"github.com/robfig/cron"
)

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 5 * * *", CronMysqlLogRename)
	c.AddFunc("0 */30 * * * *", CronWrittenHostInfo)
	c.AddFunc("0 * * * * *", CronCheckIdKey)
	c.Start()
}
