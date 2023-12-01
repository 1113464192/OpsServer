package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/util2"
	"strings"
)

type ServerService struct {
}

var (
	insServer = &ServerService{}
)

func Server() *ServerService {
	return insServer
}

// 更改单服列表
func (s *ServerService) UpdateServerRecord(param api.UpdateServerRecordReq) (result *[]api.ServerRecordRes, err error) {
	var server model.ServerRecord
	var count int64

	if err = model.DB.First(&server, param.Id).Error; err != nil {
		return nil, fmt.Errorf("查询单服记录表错误: %v", err)
	}

	// 判断Flag是否已被使用
	if model.DB.Model(&model.ServerRecord{}).Where("flag = ? AND project_id = ? AND id != ?", param.Flag, server.ProjectId, param.Id).Count(&count); count > 0 {
		return nil, errors.New("Flag已被使用")
	}
	server.Flag = param.Flag
	server.Path = param.Path
	server.ServerName = param.ServerName
	server.HostId = param.HostId
	server.ProjectId = param.ProjectId
	if err = model.DB.Save(&server).Error; err != nil {
		return nil, fmt.Errorf("数据保存失败: %v", err)
	}
	result, err = s.GetResults(&server)
	return result, err
}

// 查询单服列表
func (s *ServerService) GetServerRecord(param api.GetServerRecordReq) (result *[]api.ServerRecordRes, total int64, err error) {
	var record []model.ServerRecord
	db := model.DB.Model(&record)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &record,
		PageInfo:  param.PageInfo,
	}
	// id存在返回id对应model
	if param.Id != 0 {
		if err = db.Where("id = ?", param.Id).Count(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id数错误: %v", err)
		}
		if err = db.Where("id = ?", param.Id).First(&record).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id错误: %v", err)
		}
	} else if param.Pid != 0 && param.Flag != "" || param.ServerName != "" {
		if param.Flag != "" && param.ServerName != "" {
			return nil, 0, errors.New("单服名和单服标识不能同时填写搜索")
		}
		// 返回flag的模糊匹配
		if param.Flag != "" {
			flag := "%" + strings.ToUpper(param.Flag) + "%"
			searchReq.Condition = db.Where("project_id = ? AND UPPER(flag) LIKE ?", param.Pid, flag).Order("id desc")
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
			// 返回单服名的模糊匹配
		} else if param.ServerName != "" {
			name := "%" + strings.ToUpper(param.ServerName) + "%"
			searchReq.Condition = db.Where("project_id = ? AND UPPER(server_name) LIKE ?", param.Pid, name).Order("id desc")
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
		}
		// 返回所有
	} else {
		if param.Pid == 0 {
			return nil, 0, errors.New("请选择项目ID")
		}
		searchReq.Condition = db.Where("project_id = ?", param.Pid).Order("id desc")
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
	}
	result, err = s.GetResults(&record)
	return result, total, err
}

// 删除单服记录
func (s *ServerService) DeleteServerRecord(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.ServerRecord{}, ids); err != nil {
		return err
	}
	if err = model.DB.Where("id IN (?)", ids).Delete(&model.ServerRecord{}).Error; err != nil {
		return errors.New("删除单服记录失败")
	}
	return err
}

// 返回用户结果
func (s *ServerService) GetResults(serverInfo any) (*[]api.ServerRecordRes, error) {
	var res api.ServerRecordRes
	var result []api.ServerRecordRes
	var err error
	if recrods, ok := serverInfo.(*[]model.ServerRecord); ok {
		for _, recrod := range *recrods {
			res = api.ServerRecordRes{
				Id:         recrod.ID,
				Flag:       recrod.Flag,
				Path:       recrod.Path,
				ServerName: recrod.ServerName,
				HostId:     recrod.HostId,
				ProjectId:  recrod.ProjectId,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if recrod, ok := serverInfo.(*model.ServerRecord); ok {
		res = api.ServerRecordRes{
			Id:         recrod.ID,
			Flag:       recrod.Flag,
			Path:       recrod.Path,
			ServerName: recrod.ServerName,
			HostId:     recrod.HostId,
			ProjectId:  recrod.ProjectId,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换server结果失败")
}
