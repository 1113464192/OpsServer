package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/util2"
)

type DomainService struct {
}

var (
	insDomain = &DomainService{}
)

func Domain() *DomainService {
	return insDomain
}

// 新增或修改域名
func (s *DomainService) UpdateDomain(param *api.UpdateDomainReq) (domain *model.Domain, err error) {
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
func (s *DomainService) UpdateDomainAss(param *api.UpdateDomainAssHostReq) (err error) {
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

// 删除域名
func (s *DomainService) DeleteDomain(ids []uint) (err error) {
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
	if err = tx.Where("id IN (?)", ids).Delete(&model.Domain{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除域名失败")
	}
	tx.Commit()
	return err
}

// 获取域名关联的主机
func (s *DomainService) GetDomainAssHost(param *api.GetPagingMustByIdReq) (hostInfo any, total int64, err error) {
	var domain model.Domain
	if !util2.CheckIdExists(&domain, param.Id) {
		return nil, 0, errors.New("域名ID不存在")
	}

	// 统计被关联个数
	if err = model.DB.Find(&domain, param.Id).Error; err != nil {
		return nil, 0, errors.New("查询域名报错")
	}
	if total = model.DB.Model(&domain).Association("Hosts").Count(); total == 0 {
		return "没有关联数据", 0, nil
	}

	// 取出关联数据
	var hosts []model.Host
	if err = model.DB.Model(&domain).Order("id asc").Association("Hosts").Find(&hosts); err != nil {
		return &hosts, total, fmt.Errorf("获取关联的数据失败: %v", err)
	}

	// 分页
	if err = dbOper.DbOper().PaginateModels(&hosts, param.PageInfo); err != nil {
		return nil, 0, err
	}
	var result *[]api.HostRes
	if result, err = Host().GetResults(&hosts); err != nil {
		return nil, total, err
	}
	return result, total, err
}
