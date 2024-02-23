package crontab

import (
	"fqhWeb/configs"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util"
	"os"
	"path/filepath"
	"time"
)

func CronRemoveExpiredLogFile() {
	// 设置日志文件的过期时间，例如7天
	expiration := time.Now().AddDate(0, 0, configs.Conf.Logger.ExpiredDay)

	// 遍历logs目录下的所有文件和子目录
	if err := filepath.Walk(util.GetRootPath()+"/logs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果文件的修改时间早于过期时间，则删除该文件
		if !info.IsDir() && info.ModTime().Before(expiration) {
			if err = os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		logger.Log().Error("Log", "删除过期日志文件失败", err)
		return
	}
}
