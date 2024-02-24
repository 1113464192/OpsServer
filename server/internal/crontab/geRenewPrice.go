package crontab

import (
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/logger"
	"strconv"
)

func GetRenewPrice() {
	//	获取所有主机的续费价格并写入host表
	var (
		err   error
		hosts []model.Host
	)
	if err = model.DB.Find(&hosts).Error; err != nil {
		logger.Log().Error("Host", "获取所有主机失败: ", err)
		return
	}
	for _, host := range hosts {
		cost, err := service.Cloud().GetCloudInsRenewPrice(host.Cloud, host.ID, host.Pid)
		if err != nil {
			logger.Log().Error("Host", "获取续费价格失败: ", err)
			return
		}
		// 字符串转float
		costFloar32, err := strconv.ParseFloat(cost, 32)
		if err != nil {
			logger.Log().Error("Host", "续费价格转换失败: ", err)
			return
		}
		host.Cost = float32(costFloar32)
		if err = model.DB.Save(&host).Error; err != nil {
			logger.Log().Error("Host", "续费价格写入失败: ", err)
			return
		}
	}
}
