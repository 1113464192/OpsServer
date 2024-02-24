package crontab

import (
	"github.com/robfig/cron"
)

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 5 * * *", CronMysqlLogRename)
	c.AddFunc("0 */30 * * * *", CronWrittenHostInfo)
	c.AddFunc("0 * * * * *", CronExistCheckIdKey)
	c.AddFunc("0 0 5 * * *", CronCheckCloudProject)
	c.AddFunc("0 0 5 * * *", CronRemoveExpiredLogFile)
	c.AddFunc("0 0 4 * * *", GetRenewPrice)
	c.Start()
}
