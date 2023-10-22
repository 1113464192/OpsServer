package service

import (
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
)

type ActRecord struct {
}

var (
	insLog = ActRecord{}
)

func Record() *ActRecord {
	return &insLog
}

// RecordCreate 插入日志
func (s *ActRecord) RecordCreate(log *model.ActRecord) (err error) {
	err = model.DB.Create(&log).Error
	return err
}

// GetRecordList
// @description:  获取日志列表
// @param: param api.GetRecordReq
// @return: list any, total int64, err error
func (s *ActRecord) GetRecordList(param *api.GetRecordReq) (list any, total int64, err error) {
	var logs []model.ActRecord
	db := model.DB.Model(&logs).Order("id desc")
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &logs,
		PageInfo:  param.PageInfo,
	}
	if param.Username != "" {
		db = db.Where("username = ?", param.Username)
		searchReq.Condition = db
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
	} else {
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
	}
	return &logs, total, err
}
