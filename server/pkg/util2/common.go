package util2

// 此文件是为躲避互相import导致报错

import (
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/logger"
)

// 查询ID记录是否存在
func CheckIdExists(table any, id uint) bool {
	var count int64
	if err := model.DB.Model(&table).Where("id = ?", id).Count(&count).Error; err != nil {
		logger.Log().Error("Service", "Mysql查询ID", err)
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

// 查询IDs记录是否都存在，不存在返回所有不存在的组
func CheckIdsExists(table any, ids []uint) (err error) {
	var nonExistId []uint
	var uBool bool
	for _, id := range ids {
		if uBool = CheckIdExists(table, id); !uBool {
			// 如果不存在则添加到noexistid切片
			nonExistId = append(nonExistId, id)
		}
	}
	if len(nonExistId) != 0 {
		return fmt.Errorf("%v %s", nonExistId, "id不存在")
	}
	return err
}
