package crontab

import (
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils"
	"os"
	"time"
)

func CronMysqlLogRename() {

	now := time.Now().Local()
	previousDay := now.AddDate(0, 0, -1) // 获取前一天的日期
	logFileName := fmt.Sprintf(utils.GetRootPath()+"/logs/mysql/%s.log", previousDay.Format("20060102"))

	model.LogFile.Close() // 关闭之前的日志文件句柄

	// 重命名日志文件
	if err := os.Rename(utils.GetRootPath()+"/logs/mysql/mysql.log", logFileName); err != nil {
		logger.Log().Error("mysql", "重命名mysql日志失败", err)
		return
	}

}
