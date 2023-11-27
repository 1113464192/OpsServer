package dbOper

import (
	"errors"
	"fmt"
	"fqhWeb/pkg/api"
)

type DbOperService struct {
}

var (
	insDbOper = &DbOperService{}
)

func DbOper() *DbOperService {
	return insDbOper
}

func (s *DbOperService) DbFind(param *api.SearchReq) (int64, error) {
	var count int64
	var err error
	if err = param.Condition.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("记录数查询失败: %v", err)
	}
	if count < 1 {
		return 0, fmt.Errorf("查询到的记录为0: %v", err)
	}
	if param.PageInfo.PageSize != 0 && param.PageInfo.Page != 0 {
		limit := param.PageInfo.PageSize
		offset := (param.PageInfo.Page - 1) * param.PageInfo.PageSize
		if err = param.Condition.Limit(limit).Offset(offset).Find(param.Table).Error; err != nil {
			return 0, fmt.Errorf("数据库查询失败: %v", err)
		}
	} else {
		if err = param.Condition.Find(param.Table).Error; err != nil {
			return 0, fmt.Errorf("数据库查询失败: %v", err)
		}
	}
	if err != nil {
		return 0, errors.New("数据库操作失败")
	}
	return count, err
}

func (s *DbOperService) AssDbFind(param *api.AssQueryReq) (int64, error) {
	var count int64
	var err error
	if err = param.Condition.Count(&count).Error; err != nil {
		return 0, errors.New("记录数查询失败")
	}
	if count < 1 {
		return 0, errors.New("查询到的记录为0")
	}
	if param.PageInfo.PageSize != 0 && param.PageInfo.Page != 0 {
		limit := param.PageInfo.PageSize
		offset := (param.PageInfo.Page - 1) * param.PageInfo.PageSize
		if err = param.Condition.Offset(offset).Limit(limit).Find(param.Table).Error; err != nil {
			return 0, errors.New("数据库关联查询失败")
		}
	} else {
		if err = param.Condition.Find(param.Table).Error; err != nil {
			return 0, errors.New("数据库关联查询失败")
		}
	}
	return count, err

}
