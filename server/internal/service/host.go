package service

import (
	"database/sql"
	"errors"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/util"
	"fqhWeb/pkg/util2"
	"strconv"
	"strings"
)

type HostService struct {
}

var (
	insHost = &HostService{}
)

func Host() *HostService {
	return insHost
}

// 新增/修改服务器
func (s *HostService) UpdateHost(param *api.UpdateHostReq) (hostInfo any, err error) {
	var host *model.Host
	var count int64
	// NULL不会参与分配
	if model.DB.Model(host).Where("ipv4 = ?", param.Ipv4).Or("ipv6 = ?", param.Ipv6).Count(&count); count > 0 {
		return nil, errors.New("IP已被使用")
	}
	if param.ID != 0 {
		// 修改
		if !util2.CheckIdExists(host, param.ID) {
			return nil, errors.New("服务器ID不存在")
		}

		if err := model.DB.Where("id = ?", param.ID).First(host).Error; err != nil {
			return nil, errors.New("服务器在数据库中查询失败")
		}

		host.Ipv4 = sql.NullString{String: param.Ipv4, Valid: true}
		if param.Ipv6 != "" {
			host.Ipv6 = sql.NullString{String: param.Ipv6, Valid: true}
		}
		host.User = param.User
		host.Password, err = util.EncryptAESCBC(param.Password, []byte(consts.AesKey), []byte(consts.AesIv))
		if err != nil {
			return nil, fmt.Errorf("主机密码加密失败: %v", err)
		}
		host.Port = param.Port
		host.Zone = param.Zone
		host.ZoneTime = param.ZoneTime
		host.BillingType = param.BillingType
		host.Cost = param.Cost
		host.Cloud = param.Cloud
		host.System = param.System
		host.Type = param.Type
		host.Cores = param.Cores
		host.SystemDisk = param.SystemDisk
		host.DataDisk = param.DataDisk
		host.Iops = param.Iops
		host.Mbps = param.Mbps
		host.Mem = uint64(param.Mem) * uint64(1024)
		// 当前数据则只支持从代码中获取
		// host.CurrDisk = param.CurrDisk
		// host.CurrMem = param.CurrMem
		// host.CurrIowait = param.CurrIowait
		// host.CurrIdle = param.CurrIdle
		// host.CurrLoad = param.CurrLoad
		// 入库
		if err = model.DB.Save(host).Error; err != nil {
			return host, fmt.Errorf("数据保存失败: %v", err)
		}
		var result *[]api.HostRes
		if result, err = s.GetResults(host); err != nil {
			return nil, err
		}
		return result, err
	} else {
		var aesPassword []byte
		aesPassword, err = util.EncryptAESCBC(param.Password, []byte(consts.AesKey), []byte(consts.AesIv))
		if err != nil {
			return nil, fmt.Errorf("主机密码加密失败: %v", err)
		}
		host = &model.Host{
			Ipv4:        sql.NullString{String: param.Ipv4, Valid: true},
			User:        param.User,
			Password:    aesPassword,
			Port:        param.Port,
			Zone:        param.Zone,
			ZoneTime:    param.ZoneTime,
			BillingType: param.BillingType,
			Cost:        param.Cost,
			Cloud:       param.Cloud,
			System:      param.System,
			Type:        param.Type,
			Cores:       param.Cores,
			SystemDisk:  param.SystemDisk,
			DataDisk:    param.DataDisk,
			Iops:        param.Iops,
			Mbps:        param.Mbps,
			Mem:         uint64(param.Mem) * 1024,
			// CurrDisk:    param.CurrDisk,
			// CurrMem:     param.CurrMem,
			// CurrIowait:  param.CurrIowait,
			// CurrIdle:    param.CurrIdle,
			// CurrLoad:    param.CurrLoad,
		}
		if param.Ipv6 != "" {
			host.Ipv6 = sql.NullString{String: param.Ipv6, Valid: true}
		}
		if err = model.DB.Create(host).Error; err != nil {
			return host, errors.New("创建服务器失败")
		}
		var result *[]api.HostRes
		if result, err = s.GetResults(host); err != nil {
			return nil, err
		}
		return result, err
	}
}

// 获取服务器密码
func (s *HostService) GetHostPasswd(id uint) (string, error) {
	var err error
	var host model.Host
	if err = model.DB.First(&host, id).Error; err != nil {
		return "", fmt.Errorf("查找服务器失败: %v", err)
	}

	if host.Password != nil {
		var passwd []byte
		passwd, err = util.DecryptAESCBC(host.Password, []byte(consts.AesKey), []byte(consts.AesIv))
		if err != nil {
			return "", fmt.Errorf("主机密码解密失败: %v", err)

		}
		return string(passwd), err
	} else {
		return "", errors.New("密码为空")
	}
}

// 删除服务器
func (s *HostService) DeleteHost(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.Host{}, ids); err != nil {
		return err
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

// 删除域名
func (s *HostService) DeleteDomain(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.Domain{}, ids); err != nil {
		return err
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

// 获取主机
func (s *HostService) GetHost(param *api.GetHostReq) (hostInfo any, count int64, err error) {
	var host []model.Host
	ipstr := "%" + param.Ip + "%"
	if err := model.DB.Model(&host).Where("UPPER(name) LIKE ?", ipstr).Count(&count).Error; err != nil || count < 1 {
		return nil, 0, errors.New("记录总数查询失败或不存在该搜索内容")
	}
	db := model.DB.Model(&host)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &host,
		PageInfo:  param.PageInfo,
	}
	name := "%" + param.Ip + "%"
	if param.Ip != "" {
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
	var result *[]api.HostRes
	if result, err = s.GetResults(&host); err != nil {
		return nil, 0, err
	}
	return result, count, err
}

// 获取域名关联的主机
func (s *HostService) GetDomainAssHost(param *api.GetPagingByIdReq) (hostInfo any, total int64, err error) {
	var domain model.Domain
	if !util2.CheckIdExists(&domain, param.Id) {
		return nil, 0, errors.New("域名ID不存在")
	}
	if err = model.DB.Preload("Hosts").Where("id = ?", param.Id).First(&domain).Error; err != nil {
		return nil, 0, errors.New("域名查询失败")
	}
	assQueryReq := &api.AssQueryReq{
		Condition: model.DB.Model(&model.Host{}),
		Table:     &domain.Hosts,
		PageInfo:  param.PageInfo,
	}
	if total, err = dbOper.DbOper().AssDbFind(assQueryReq); err != nil {
		return nil, 0, err
	}
	var result *[]api.HostRes
	if result, err = Host().GetResults(&domain.Hosts); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 新增或修改域名
func (s *HostService) UpdateDomain(param *api.UpdateDomainReq) (domain *model.Domain, err error) {
	var count int64
	// NULL不会参与分配
	if model.DB.Model(domain).Where("Value = ?", param.Value).Count(&count); count > 0 {
		return nil, errors.New("域名已被使用")
	}
	if param.Id != 0 {
		// 修改
		if !util2.CheckIdExists(domain, param.Id) {
			return nil, errors.New("域名ID不存在")
		}

		if err := model.DB.Where("id = ?", param.Id).First(domain).Error; err != nil {
			return nil, errors.New("服域名在数据库中查询失败")
		}
		domain.Value = param.Value
		if err = model.DB.Save(domain).Error; err != nil {
			return domain, fmt.Errorf("数据保存失败: %v", err)
		}
		return domain, err
	} else {
		domain = &model.Domain{
			Value: param.Value,
		}
		if err = model.DB.Create(domain).Error; err != nil {
			return domain, errors.New("创建域名失败")
		}
		return domain, err
	}
}

// 更新域名关联的主机
func (s *HostService) UpdateDomainAss(param *api.UpdateDomainAssHostReq) (err error) {
	var host []model.Host
	var domain model.Domain
	// 判断所有项目是否都存在
	if err = util2.CheckIdsExists(model.Host{}, param.Hids); err != nil {
		return err
	}

	if !util2.CheckIdExists(&host, param.Did) {
		return errors.New("域名ID不存在")
	}

	if err = model.DB.Find(&host, param.Hids).Error; err != nil {
		return errors.New("服务器数据库查询操作失败")
	}
	if err = model.DB.First(&domain, param.Did).Error; err != nil {
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

// 获取主机当前状态
func (s *HostService) GetHostCurrData(param *[]api.SSHClientConfigReq) (*api.HostInfoRes, error) {
	// systemDiskShell := `df -Th | awk '{if ($NF=="/")print$(NF-2)}' | grep -Eo "[0-9]+"`
	// dataDiskShell := `df -Th | awk '{if ($NF=="/data")print$(NF-2)}' | grep -Eo "[0-9]+"`
	// memShell := `free -m | awk '/Mem/{print $NF}'`
	// iowaitShell := `iostat | awk '/avg-cpu:/ {getline; print $(NF-2)}'`
	// idleShell := `iostat | awk '/avg-cpu:/ {getline; print $(NF)}'`
	// loadShell := `uptime | awk -F"[, ]+" '{print $(NF-1)}'`
	cmdShell := `systemDisk=$(df -Th | awk '{if ($NF=="/")print$(NF-2)}' | grep -Eo "[0-9]+")
				 dataDisk=$(df -Th | awk '{if ($NF=="/data")print$(NF-2)}' | grep -Eo "[0-9]+")
				 if [[ -z ${dataDisk} ]];then
				 	dataDisk=-1
				 fi
				 mem=$(free -m | awk '/Mem/{print $NF}')
				 iowait=$(iostat | awk '/avg-cpu:/ {getline; print $(NF-2)}')
				 idle=$(iostat | awk '/avg-cpu:/ {getline; print $(NF)}')
				 load=$(uptime | awk -F"[, ]+" '{print $(NF-1)}')
				 echo "$systemDisk $dataDisk $mem $iowait $idle $load" | awk '{print $1,$2,$3,$4,$5,$6}'`

	var hostInfo api.HostInfoRes
	var err error
	// 对各个sshReq写入CMD
	for i := 0; i < len(*param); i++ {
		(*param)[i].Cmd = cmdShell
	}
	hostDataRes, err := SSH().RunSSHCmdAsync(param)
	// 返回*[]SSHResultRes
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(*hostDataRes); i++ {
		splitRes := strings.Fields((*hostDataRes)[i].Response)
		// 这里请为我实现将所有value依次赋给对应的key
		if len(splitRes) == 6 {
			hostInfo.CurrSystemDisk = append(hostInfo.CurrSystemDisk, api.SSHResultRes{
				HostIp:   (*hostDataRes)[i].HostIp,
				Status:   (*hostDataRes)[i].Status,
				Response: splitRes[0],
			})
			hostInfo.CurrDataDisk = append(hostInfo.CurrDataDisk, api.SSHResultRes{
				HostIp:   (*hostDataRes)[i].HostIp,
				Status:   (*hostDataRes)[i].Status,
				Response: splitRes[1],
			})
			hostInfo.CurrMem = append(hostInfo.CurrMem, api.SSHResultRes{
				HostIp:   (*hostDataRes)[i].HostIp,
				Status:   (*hostDataRes)[i].Status,
				Response: splitRes[2],
			})
			hostInfo.CurrIowait = append(hostInfo.CurrIowait, api.SSHResultRes{
				HostIp:   (*hostDataRes)[i].HostIp,
				Status:   (*hostDataRes)[i].Status,
				Response: splitRes[3],
			})
			hostInfo.CurrIdle = append(hostInfo.CurrIdle, api.SSHResultRes{
				HostIp:   (*hostDataRes)[i].HostIp,
				Status:   (*hostDataRes)[i].Status,
				Response: splitRes[4],
			})
			hostInfo.CurrLoad = append(hostInfo.CurrLoad, api.SSHResultRes{
				HostIp:   (*hostDataRes)[i].HostIp,
				Status:   (*hostDataRes)[i].Status,
				Response: splitRes[5],
			})
		}
	}
	return &hostInfo, err
}

// hostInfo.CurrSystemDisk = systemDiskRes

// param.Cmd = []string{systemDiskShell}
// systemDiskRes, err := SSH().RunSSHCmdAsync(param)
// if err != nil {
// 	return nil, err
// }
// for i := range *systemDiskRes {
// 	(*systemDiskRes)[i].Response = strings.TrimSpace((*systemDiskRes)[i].Response)
// }
// hostInfo.CurrSystemDisk = systemDiskRes

// param.Cmd = []string{dataDiskShell}
// dataDiskRes, err := SSH().RunSSHCmdAsync(param)
// if err != nil {
// 	return nil, err
// }
// for i := range *dataDiskRes {
// 	(*dataDiskRes)[i].Response = strings.TrimSpace((*dataDiskRes)[i].Response)
// }
// hostInfo.CurrDataDisk = dataDiskRes

// param.Cmd = []string{memShell}
// memRes, err := SSH().RunSSHCmdAsync(param)
// if err != nil {
// 	return nil, err
// }
// for i := range *memRes {
// 	memDataStr := strings.TrimSpace((*memRes)[i].Response)
// 	memData, err := strconv.Atoi(memDataStr)
// 	if err != nil {
// 		return nil, fmt.Errorf(" 字符串转换整数失败: %v", err)
// 	}
// 	// 将内存除以 1024，并转换为以 "G" 为单位的大小
// 	memSize := float64(memData) / float64(1024)
// 	(*memRes)[i].Response = strconv.FormatFloat(memSize, 'f', 2, 32)
// }
// hostInfo.CurrMem = memRes

// param.Cmd = []string{iowaitShell}
// iowaitRes, err := SSH().RunSSHCmdAsync(param)
// if err != nil {
// 	return nil, err
// }
// for i := range *iowaitRes {
// 	(*iowaitRes)[i].Response = strings.TrimSpace((*iowaitRes)[i].Response)
// }
// hostInfo.CurrIowait = iowaitRes

// param.Cmd = []string{idleShell}
// idleRes, err := SSH().RunSSHCmdAsync(param)
// if err != nil {
// 	return nil, err
// }
// for i := range *idleRes {
// 	(*idleRes)[i].Response = strings.TrimSpace((*idleRes)[i].Response)
// }
// hostInfo.CurrIdle = idleRes

// param.Cmd = []string{loadShell}
// loadRes, err := SSH().RunSSHCmdAsync(param)
// if err != nil {
// 	return nil, err
// }
// for i := range *loadRes {
// 	(*loadRes)[i].Response = strings.TrimSpace((*loadRes)[i].Response)
// }
// hostInfo.CurrLoad = loadRes

// 写入主机信息到数据库
func (s *HostService) WritieToDatabase(data *api.HostInfoRes) error {
	var host model.Host
	tx := model.DB.Begin()
	// 如果status非0则全部-1
	// for _, hostRes := range *data.CurrSystemDisk {
	for i := 0; i < len(data.CurrSystemDisk); i++ {
		if (data.CurrSystemDisk)[i].Status != 0 {
			(data.CurrSystemDisk)[i].Response = "-1"
		}
		currSystemDisk, err := strconv.ParseFloat((data.CurrSystemDisk)[i].Response, 32)
		if err != nil {
			return fmt.Errorf("字符串转换浮点数错误: %v", err)
		}
		if (data.CurrSystemDisk)[i].Status != 0 {
			(data.CurrDataDisk)[i].Response = "-1"
		}
		currDataDisk, err := strconv.ParseFloat((data.CurrDataDisk)[i].Response, 32)
		if err != nil {
			return fmt.Errorf("字符串转换浮点数错误: %v", err)
		}
		if (data.CurrSystemDisk)[i].Status != 0 {
			(data.CurrMem)[i].Response = "-1"
		}
		currMem, err := strconv.ParseFloat((data.CurrMem)[i].Response, 32)
		if err != nil {
			return fmt.Errorf("字符串转换浮点数错误: %v", err)
		}
		if (data.CurrSystemDisk)[i].Status != 0 {
			(data.CurrIdle)[i].Response = "-1"
		}
		currIdle, err := strconv.ParseFloat((data.CurrIdle)[i].Response, 32)
		if err != nil {
			return fmt.Errorf("字符串转换浮点数错误: %v", err)
		}
		if (data.CurrSystemDisk)[i].Status != 0 {
			(data.CurrIowait)[i].Response = "-1"
		}
		currIowait, err := strconv.ParseFloat((data.CurrIowait)[i].Response, 32)
		if err != nil {
			return fmt.Errorf("字符串转换浮点数错误: %v", err)
		}
		if (data.CurrSystemDisk)[i].Status != 0 {
			(data.CurrLoad)[i].Response = "-1"
		}
		currLoad, err := strconv.ParseFloat((data.CurrLoad)[i].Response, 32)
		if err != nil {
			return fmt.Errorf("字符串转换浮点数错误: %v", err)
		}
		if err = tx.Model(&host).Where("ipv4 = ?", (data.CurrSystemDisk)[i].HostIp).Updates(model.Host{CurrSystemDisk: float32(currSystemDisk), CurrDataDisk: float32(currDataDisk), CurrMem: float32(currMem), CurrIowait: float32(currIowait), CurrIdle: float32(currIdle), CurrLoad: float32(currLoad)}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新当前服务器状态到数据库失败: %v", err)
		}
	}
	tx.Commit()
	return nil
}

// 返回结果
func (s *HostService) GetResults(hostInfo any) (*[]api.HostRes, error) {
	var res api.HostRes
	var result []api.HostRes
	var err error
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
				Cost:           host.Cost,
				Cloud:          host.Cloud,
				System:         host.System,
				Type:           host.Type,
				Cores:          host.Cores,
				SystemDisk:     host.SystemDisk,
				DataDisk:       host.DataDisk,
				Iops:           host.Iops,
				Mbps:           host.Mbps,
				Mem:            uint32(host.Mem) / uint32(1024),
				CurrSystemDisk: host.CurrSystemDisk,
				CurrDataDisk:   host.CurrDataDisk,
				CurrMem:        host.CurrMem,
				CurrIowait:     host.CurrIowait,
				CurrIdle:       host.CurrIdle,
				CurrLoad:       host.CurrLoad,
			}
			result = append(result, res)
		}
		return &result, err
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
			Cost:           host.Cost,
			Cloud:          host.Cloud,
			System:         host.System,
			Type:           host.Type,
			Cores:          host.Cores,
			SystemDisk:     host.SystemDisk,
			DataDisk:       host.DataDisk,
			Iops:           host.Iops,
			Mbps:           host.Mbps,
			Mem:            uint32(host.Mem) / uint32(1024),
			CurrSystemDisk: host.CurrSystemDisk,
			CurrDataDisk:   host.CurrDataDisk,
			CurrMem:        host.CurrMem,
			CurrIowait:     host.CurrIowait,
			CurrIdle:       host.CurrIdle,
			CurrLoad:       host.CurrLoad,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换服务器结果失败")
}
