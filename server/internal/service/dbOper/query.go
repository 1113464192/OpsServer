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

func (s *DbOperService) DbFind(params *api.SearchReq) (int64, error) {
	var count int64
	var err error
	if err = params.Condition.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("记录数查询失败: %v", err)
	}
	if count < 1 {
		return 0, fmt.Errorf("查询到的记录为0: %v", err)
	}
	if params.PageInfo.PageSize != 0 && params.PageInfo.Page != 0 {
		limit := params.PageInfo.PageSize
		offset := (params.PageInfo.Page - 1) * params.PageInfo.PageSize
		if err = params.Condition.Limit(limit).Offset(offset).Find(params.Table).Error; err != nil {
			return 0, fmt.Errorf("数据库查询失败: %v", err)
		}
	} else {
		if err = params.Condition.Find(params.Table).Error; err != nil {
			return 0, fmt.Errorf("数据库查询失败: %v", err)
		}
	}
	if err != nil {
		return 0, errors.New("数据库操作失败")
	}
	return count, err
}

func (s *DbOperService) AssDbFind(params *api.AssQueryReq) (int64, error) {
	var count int64
	var err error
	if err = params.Condition.Count(&count).Error; err != nil {
		return 0, errors.New("记录数查询失败")
	}
	if count < 1 {
		return 0, errors.New("查询到的记录为0")
	}
	if params.PageInfo.PageSize != 0 && params.PageInfo.Page != 0 {
		limit := params.PageInfo.PageSize
		offset := (params.PageInfo.Page - 1) * params.PageInfo.PageSize
		if err = params.Condition.Offset(offset).Limit(limit).Find(params.Table).Error; err != nil {
			return 0, errors.New("数据库关联查询失败")
		}
	} else {
		if err = params.Condition.Find(params.Table).Error; err != nil {
			return 0, errors.New("数据库关联查询失败")
		}
	}
	return count, err

}
