package dbOper

import (
	"errors"
	"fmt"
	"fqhWeb/pkg/api"
	"reflect"
	"sort"
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
		return 0, nil
	}
	if param.PageInfo.PageSize != 0 && param.PageInfo.Page != 0 {
		limit := param.PageInfo.PageSize
		offset := (param.PageInfo.Page - 1) * param.PageInfo.PageSize
		if err = param.Condition.Limit(limit).Offset(offset).Find(param.Table).Error; err != nil {
			return 0, fmt.Errorf("数据库查询失败: %v", err)
		}
	} else {
		if err = param.Condition.Find(param.Table).Error; err != nil {
			return 0, fmt.Errorf("数据库操作失败: %v", err)
		}
	}
	if err != nil {
		return 0, fmt.Errorf("数据库操作失败: %v", err)
	}
	return count, err
}

// 排序并分页
func (s *DbOperService) PaginateAndSortModels(menusPtr interface{}, pageInfo api.PageInfo, lessFunc func(i, j int) bool) error {
	v := reflect.ValueOf(menusPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("传入的值不是一个切片")
	}

	if pageInfo.PageSize != 0 && pageInfo.Page != 0 {
		// 获取到menusPtr的元素并变成interface类型做排序
		sort.Slice(v.Elem().Interface(), lessFunc)

		// 反射到原本指针指向的内容
		sliceValue := v.Elem()
		limit := pageInfo.PageSize
		offset := (pageInfo.Page - 1) * pageInfo.PageSize
		start := offset
		end := offset + limit
		if end > sliceValue.Len() {
			end = sliceValue.Len()
		}
		paginatedSlice := sliceValue.Slice(start, end)
		v.Elem().Set(paginatedSlice)
	}

	return nil
}

// 分页
func (s *DbOperService) PaginateModels(menusPtr interface{}, pageInfo api.PageInfo) error {
	v := reflect.ValueOf(menusPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("传入的值不是一个切片")
	}

	if pageInfo.PageSize != 0 && pageInfo.Page != 0 {
		// 反射到原本指针指向的内容
		sliceValue := v.Elem()
		limit := pageInfo.PageSize
		offset := (pageInfo.Page - 1) * pageInfo.PageSize
		start := offset
		end := offset + limit
		if end > sliceValue.Len() {
			end = sliceValue.Len()
		}
		paginatedSlice := sliceValue.Slice(start, end)
		v.Elem().Set(paginatedSlice)
	}

	return nil
}
