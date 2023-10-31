package service

import (
	"database/sql"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/utils"
	"fqhWeb/pkg/utils2"
)

type HostService struct {
}

var (
	insHost = &HostService{}
)

func Host() *HostService {
	return insHost
}

func (s *HostService) UpdateHost(params *api.UpdateHostReq) (hostInfo any, err error) {
	var host *model.Host
	var count int64
	// NULL不会参与分配
	if model.DB.Model(host).Where("ipv4 = ?", params.Ipv4).Or("ipv6 = ?", params.Ipv6).Count(&count); count > 0 {
		return nil, errors.New("IP已被使用")
	}
	if params.ID != 0 {
		// 修改
		if !utils2.CheckIdExists(host, &params.ID) {
			return nil, errors.New("服务器ID不存在")
		}

		if err := model.DB.Where("id = ?", params.ID).First(host).Error; err != nil {
			return nil, errors.New("服务器在数据库中查询失败")
		}

		host.Ipv4 = sql.NullString{String: params.Ipv4, Valid: true}
		if params.Ipv6 != "" {
			host.Ipv6 = sql.NullString{String: params.Ipv6, Valid: true}
		}
		host.Password, err = utils.GenerateFromPassword(params.Password)
		if err != nil {
			return host, errors.New("服务器密码bcrypt加密失败")
		}
		host.User = params.User
		host.Password = params.Password
		host.Port = params.Port
		host.Zone = params.Zone
		host.ZoneTime = params.ZoneTime
		host.BillingType = params.BillingType
		host.Cloud = params.Cloud
		host.System = params.System
		host.Type = params.Type
		host.Cores = params.Cores
		host.SystemDisk = params.SystemDisk
		host.DataDisk = params.DataDisk
		host.Iops = params.Iops
		host.Mbps = params.Mbps
		host.Mem = params.Mem
		// 只支持从代码中获取
		// host.CurrDisk = params.CurrDisk
		// host.CurrMem = params.CurrMem
		// host.CurrIowait = params.CurrIowait
		// host.CurrIdle = params.CurrIdle
		// host.CurrLoad = params.CurrLoad
		err = model.DB.Save(host).Error
		if err != nil {
			return host, errors.New("数据保存失败")
		}
		var result []api.HostRes
		if result, err = s.GetResults(host); err != nil {
			return nil, err
		}
		return result, err
	} else {
		host = &model.Host{
			Ipv4:        sql.NullString{String: params.Ipv4, Valid: true},
			User:        params.User,
			Password:    params.Password,
			Port:        params.Port,
			Zone:        params.Zone,
			ZoneTime:    params.ZoneTime,
			BillingType: params.BillingType,
			Cloud:       params.Cloud,
			System:      params.System,
			Type:        params.Type,
			Cores:       params.Cores,
			SystemDisk:  params.SystemDisk,
			DataDisk:    params.DataDisk,
			Iops:        params.Iops,
			Mbps:        params.Mbps,
			Mem:         params.Mem,
			// CurrDisk:    params.CurrDisk,
			// CurrMem:     params.CurrMem,
			// CurrIowait:  params.CurrIowait,
			// CurrIdle:    params.CurrIdle,
			// CurrLoad:    params.CurrLoad,
		}
		if params.Ipv6 != "" {
			host.Ipv6 = sql.NullString{String: params.Ipv6, Valid: true}
		}
		if err = model.DB.Create(host).Error; err != nil {
			return host, errors.New("创建服务器失败")
		}
		var result []api.HostRes
		if result, err = s.GetResults(host); err != nil {
			return nil, err
		}
		return result, err
	}
}

func (s *HostService) DeleteHost(ids []uint) (err error) {
	for _, i := range ids {
		if !utils2.CheckIdExists(&model.Host{}, &i) {
			return errors.New("服务器不存在")
		}
	}
	var host []model.Host
	tx := model.DB.Begin()
	if err = tx.Find(&host, ids).Error; err != nil {
		return errors.New("查询服务器信息失败")
	}
	if err = tx.Model(&host).Association("Domains").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 服务器与域名关联 失败")
	}
	if err = tx.Model(&host).Association("Projects").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 服务器与项目关联 失败")
	}
	if err = tx.Model(&host).Association("TaskTemplate").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 服务器与任务模板关联 失败")
	}
	if err = tx.Where("id in (?)", ids).Delete(&model.Host{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除服务器失败")
	}
	tx.Commit()
	return err
}

func (s *HostService) DeleteDomain(ids []uint) (err error) {
	for _, i := range ids {
		if !utils2.CheckIdExists(&model.Host{}, &i) {
			return errors.New("域名不存在")
		}
	}
	var domain []model.Domain
	tx := model.DB.Begin()
	if err = tx.Find(&domain, ids).Error; err != nil {
		return errors.New("查询域名信息失败")
	}
	if err = tx.Model(&domain).Association("Hosts").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 域名与服务器关联 失败")
	}
	if err = tx.Where("id in (?)", ids).Delete(&model.Domain{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除域名失败")
	}
	tx.Commit()
	return err
}

func (s *HostService) UpdateProjectAss(params *api.UpdateHostAssProjectReq) (err error) {
	var host model.Host
	var noExistId []uint
	var project []model.Project
	// 判断所有项目是否都存在
	for _, pid := range params.Pids {
		uBool := utils2.CheckIdExists(&project, &pid)
		if !uBool {
			noExistId = append(noExistId, pid)
		}
	}
	if len(noExistId) != 0 {
		return fmt.Errorf("%v %s", noExistId, "项目不存在")
	}

	if !utils2.CheckIdExists(&host, &params.Hid) {
		return errors.New("服务器ID不存在")
	}

	if err = model.DB.Find(&project, params.Pids).Error; err != nil {
		return errors.New("项目数据库查询操作失败")
	}
	if err = model.DB.First(&host, params.Hid).Error; err != nil {
		return errors.New("服务器数据库查询操作失败")
	}
	if err = model.DB.Model(&host).Association("Projects").Replace(&project); err != nil {
		return errors.New("项目与服务器数据库关联操作失败")
	}
	if err != nil {
		return err
	}
	return err
}

func (s *HostService) GetHost(params *api.GetHostReq) (hostInfo any, count int64, err error) {
	var host []model.Host
	ipstr := "%" + params.Ip + "%"
	if err := model.DB.Model(&host).Where("UPPER(name) LIKE ?", ipstr).Count(&count).Error; err != nil || count < 1 {
		return nil, 0, errors.New("记录总数查询失败或不存在该搜索内容")
	}
	db := model.DB.Model(&host)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &host,
		PageInfo:  params.PageInfo,
	}
	name := "%" + params.Ip + "%"
	if params.Ip != "" {
		db = model.DB.Where("name LIKE ?", name)
		searchReq.Condition = db
		if count, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
	} else {
		if count, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
	}
	var result []api.HostRes
	if result, err = s.GetResults(&host); err != nil {
		return nil, 0, err
	}
	return result, count, err
}

// 获取对应关联项目
func (s *HostService) GetProject(params *api.GetHostAssProjectReq) (projectInfo any, total int64, err error) {
	var host model.Host
	if !utils2.CheckIdExists(&host, &params.Id) {
		return nil, 0, errors.New("主机ID不存在")
	}
	if err = model.DB.Preload("Projects").Where("id = ?", params.Id).First(&host).Error; err != nil {
		return nil, 0, errors.New("主机查询失败")
	}
	assQueryReq := &api.AssQueryReq{
		Condition: model.DB.Model(&model.Project{}),
		Table:     &host.Projects,
		PageInfo:  params.PageInfo,
	}
	if total, err = dbOper.DbOper().AssDbFind(assQueryReq); err != nil {
		return nil, 0, err
	}
	var result []api.ProjectRes
	if result, err = Project().GetResults(&host.Projects); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 获取对应关联项目
func (s *HostService) GetDomainAssHost(params *api.GetHostAssProjectReq) (hostInfo any, total int64, err error) {
	var domain model.Domain
	if !utils2.CheckIdExists(&domain, &params.Id) {
		return nil, 0, errors.New("域名ID不存在")
	}
	if err = model.DB.Preload("Hosts").Where("id = ?", params.Id).First(&domain).Error; err != nil {
		return nil, 0, errors.New("域名查询失败")
	}
	assQueryReq := &api.AssQueryReq{
		Condition: model.DB.Model(&model.Host{}),
		Table:     &domain.Hosts,
		PageInfo:  params.PageInfo,
	}
	if total, err = dbOper.DbOper().AssDbFind(assQueryReq); err != nil {
		return nil, 0, err
	}
	var result []api.HostRes
	if result, err = Host().GetResults(&domain.Hosts); err != nil {
		return nil, total, err
	}
	return result, total, err
}

func (s *HostService) UpdateDomain(params *api.UpdateDomainReq) (domain *model.Domain, err error) {
	var count int64
	// NULL不会参与分配
	if model.DB.Model(domain).Where("Value = ?", params.Value).Count(&count); count > 0 {
		return nil, errors.New("域名已被使用")
	}
	if params.Id != 0 {
		// 修改
		if !utils2.CheckIdExists(domain, &params.Id) {
			return nil, errors.New("域名ID不存在")
		}

		if err := model.DB.Where("id = ?", params.Id).First(domain).Error; err != nil {
			return nil, errors.New("服域名在数据库中查询失败")
		}
		domain.Value = params.Value
		err = model.DB.Save(domain).Error
		if err != nil {
			return domain, errors.New("数据保存失败")
		}
		return domain, err
	} else {
		domain = &model.Domain{
			Value: params.Value,
		}
		if err = model.DB.Create(domain).Error; err != nil {
			return domain, errors.New("创建域名失败")
		}
		return domain, err
	}
}

func (s *HostService) UpdateDomainAss(params *api.UpdateDomainAssHostReq) (err error) {
	var host []model.Host
	var noExistId []uint
	var domain model.Domain
	// 判断所有项目是否都存在
	for _, hid := range params.Hids {
		uBool := utils2.CheckIdExists(&host, &hid)
		if !uBool {
			noExistId = append(noExistId, hid)
		}
	}
	if len(noExistId) != 0 {
		return fmt.Errorf("%v %s", noExistId, "服务器不存在")
	}

	if !utils2.CheckIdExists(&host, &params.Did) {
		return errors.New("域名ID不存在")
	}

	if err = model.DB.Find(&host, params.Hids).Error; err != nil {
		return errors.New("服务器数据库查询操作失败")
	}
	if err = model.DB.First(&domain, params.Did).Error; err != nil {
		return errors.New("域名数据库查询操作失败")
	}
	if err = model.DB.Model(&domain).Association("Hosts").Replace(&host); err != nil {
		return errors.New("项目与服务器数据库关联操作失败")
	}
	if err != nil {
		return err
	}
	return err
}

// 返回结果
func (s *HostService) GetResults(hostInfo any) (result []api.HostRes, err error) {
	var res api.HostRes
	if hosts, ok := hostInfo.(*[]model.Host); ok {
		for _, host := range *hosts {
			res = api.HostRes{
				ID:             host.ID,
				Ipv4:           host.Ipv4.String,
				Ipv6:           host.Ipv6.String,
				Port:           host.Port,
				Zone:           host.Zone,
				ZoneTime:       host.ZoneTime,
				BillingType:    host.BillingType,
				Cloud:          host.Cloud,
				System:         host.System,
				Type:           host.Type,
				Cores:          host.Cores,
				SystemDisk:     host.SystemDisk,
				DataDisk:       host.DataDisk,
				Iops:           host.Iops,
				Mbps:           host.Mbps,
				Mem:            host.Mem,
				CurrSystemDisk: host.CurrSystemDisk,
				CurrDataDisk:   host.CurrDataDisk,
				CurrMem:        host.CurrMem,
				CurrIowait:     host.CurrIowait,
				CurrIdle:       host.CurrIdle,
				CurrLoad:       host.CurrLoad,
			}
			result = append(result, res)
		}
		return result, err
	}
	if host, ok := hostInfo.(*model.Host); ok {
		res = api.HostRes{
			ID:             host.ID,
			Ipv4:           host.Ipv4.String,
			Ipv6:           host.Ipv6.String,
			Port:           host.Port,
			Zone:           host.Zone,
			ZoneTime:       host.ZoneTime,
			BillingType:    host.BillingType,
			Cloud:          host.Cloud,
			System:         host.System,
			Type:           host.Type,
			Cores:          host.Cores,
			SystemDisk:     host.SystemDisk,
			DataDisk:       host.DataDisk,
			Iops:           host.Iops,
			Mbps:           host.Mbps,
			Mem:            host.Mem,
			CurrSystemDisk: host.CurrSystemDisk,
			CurrDataDisk:   host.CurrDataDisk,
			CurrMem:        host.CurrMem,
			CurrIowait:     host.CurrIowait,
			CurrIdle:       host.CurrIdle,
			CurrLoad:       host.CurrLoad,
		}
		result = append(result, res)
		return result, err
	}
	return result, errors.New("转换服务器结果失败")
}
