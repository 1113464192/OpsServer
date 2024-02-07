package service

import (
	"errors"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
)

type ApiService struct {
}

var (
	insApi = &ApiService{}
)

func Api() *ApiService {
	return insApi
}

// GetApiList
// @description:  获取API列表
// @param: param api.PageInfo
// @return:  list any, total int64, err error
func (s *ApiService) GetApiList(param api.PageInfo) (list any, total int64, err error) {
	var modelApi []model.Api
	db := model.DB.Model(&modelApi)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &modelApi,
		PageInfo:  param,
	}
	if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
		return nil, 0, err
	}
	return &modelApi, total, err
}

// UpdateApi
// @description:  新增/修改API
// @param: param *api.UpdateApiReq
// @return:  apiInter model.Api, err error
func (s *ApiService) UpdateApi(param *api.UpdateApiReq) (apiInter model.Api, err error) {
	var apiRes model.Api
	//判断有无ID过来，有就是修改，没有就是新增
	if param.ID != 0 {
		// 修改
		if err = model.DB.Where("id = ?", param.ID).Find(&apiRes).Error; err != nil {
			return apiInter, err
		}
		if err = CasbinServiceApp().UpdateCasbinApi(apiRes.Path, param.Path, apiRes.Method, param.Method); err != nil {
			return apiInter, err
		}
		apiRes.Path = param.Path
		apiRes.Description = param.Description
		apiRes.ApiGroup = param.ApiGroup
		apiRes.Method = param.Method
		if err = model.DB.Save(&apiRes).Error; err != nil {
			return apiInter, err
		}
	} else {
		// 新增
		count := int64(0)
		model.DB.Model(&model.Api{}).Where("path = ?", param.Path).Count(&count)
		if count > 0 {
			return apiInter, errors.New("api已存在")
		}
		apiRes = model.Api{
			Path:        param.Path,
			Description: param.Description,
			ApiGroup:    param.ApiGroup,
			Method:      param.Method,
		}
		if err = model.DB.Save(&apiRes).Error; err != nil {
			return apiInter, err
		}
	}
	return apiRes, err
}

// DeleteApi
// @description:  删除API
// @param ids []uint
// @return:  err error
func (s *ApiService) DeleteApi(ids []uint) (err error) {
	err = model.DB.Where("id IN (?)", ids).Delete(&model.Api{}).Error
	return err
}

// FreshCasbin
// @description:  刷新casbin缓存
// @return:  err error
func (s *ApiService) FreshCasbin() (err error) {
	e := CasbinServiceApp().Casbin()
	err = e.LoadPolicy()
	return err
}
