package utils2

// 此文件是为躲避互相import导致报错

import (
	"fqhWeb/internal/model"
	"fqhWeb/pkg/logger"
)

// 查询ID记录是否存在
func CheckIdExists(table any, id *uint) bool {
	var count int64
	if err := model.DB.Model(&table).Where("id = ?", *id).Count(&count).Error; err != nil {
		logger.Log().Error("Service", "Mysql查询ID", err)
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}
