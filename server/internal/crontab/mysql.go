package crontab

import (
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util"
	"os"
	"time"
)

func CronMysqlLogRename() {
	now := time.Now().Local()
	previousDay := now.AddDate(0, 0, -1) // 获取前一天的日期
	logFileName := fmt.Sprintf(util.GetRootPath()+"/logs/mysql/%s.log", previousDay.Format("20060102"))

	// 关闭之前的日志文件描述符
	if err := model.LogFile.Close(); err != nil {
		logger.Log().Error("Mysql", "关闭mysql日志文件描述符失败", err)
		return
	}

	// 重命名日志文件
	if err := os.Rename(util.GetRootPath()+"/logs/mysql/mysql.log", logFileName); err != nil {
		logger.Log().Error("Mysql", "重命名mysql日志失败", err)
		return
	}

}
