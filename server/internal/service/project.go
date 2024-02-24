package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util2"
	"strings"
)

type ProjectService struct {
}

var (
	insProject = &ProjectService{}
)

func Project() *ProjectService {
	return insProject
}

// 修改或新增项目
func (s *ProjectService) UpdateProject(param *api.UpdateProjectReq) (projectInfo any, err error) {
	var project model.Project
	var count int64
	// 判断项目名是否已被使用
	if model.DB.Model(&project).Where("name = ? AND id != ?", param.Name, param.ID).Count(&count); count > 0 {
		return &project, errors.New("项目名已被使用")
	}
	// 判断是否为支持的云商
	switch param.Cloud {
	case "腾讯云":
	default:
		return nil, errors.New("云平台不支持")
	}
	// ID查询
	if param.ID != 0 {
		if !util2.CheckIdExists(&project, param.ID) {
			return project, errors.New("项目不存在")
		}

		if err := model.DB.Where("id = ?", param.ID).First(&project).Error; err != nil {
			return &project, errors.New("项目数据库查询失败")
		}
		if project.Name != param.Name {
			return nil, errors.New("项目名不允许修改,建议删除重建新项目。 （因为牵扯服务太多容易出问题,如必要请通知运维逐个服务添加递归修改)")
		}
		if project.Cloud != param.Cloud {
			return nil, errors.New("项目云平台不允许修改")
		}
		// 更改云平台项目属性
		if project.Status != param.Status {
			cloudPid, err := Cloud().GetCloudProjectId(project.Cloud, project.Name)
			if err != nil {
				return nil, err
			}
			if err = Cloud().UpdateCloudProject(project.Cloud, cloudPid, param.Name, int64(param.Status)); err != nil {
				return nil, err
			}
		}

		project.Name = param.Name
		project.Status = param.Status
		project.UserId = param.UserId
		project.GroupId = param.GroupId
		if err = model.DB.Save(&project).Error; err != nil {
			return project, fmt.Errorf("数据保存失败: %v", err)
		}
		var result *[]api.ProjectRes
		if result, err = s.GetResults(&project); err != nil {
			return nil, err
		}
		return result, err
	} else {
		project = model.Project{
			Name:    param.Name,
			Cloud:   param.Cloud,
			Status:  param.Status,
			UserId:  param.UserId,
			GroupId: param.GroupId,
		}

		// 创建云平台项目
		if err = Cloud().CreateCloudProject(param.Cloud, param.Name); err != nil {
			return nil, err
		}

		if err = model.DB.Create(&project).Error; err != nil {
			logger.Log().Error("project", "创建项目失败", err)
			return project, errors.New("创建项目失败")
		}
		var result *[]api.ProjectRes
		if result, err = s.GetResults(&project); err != nil {
			return nil, err
		}
		return result, err
	}
}

// 删除项目
func (s *ProjectService) DeleteProject(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.Project{}, ids); err != nil {
		return err
	}
	var projects []model.Project
	tx := model.DB.Begin()
	if err = tx.Find(&projects, ids).Error; err != nil {
		return errors.New("查询项目信息失败")
	}
	//if err = tx.Model(&projects).Association("Hosts").Clear(); err != nil {
	//	tx.Rollback()
	//	return errors.New("清除表信息 项目与服务器关联 失败")
	//}
	if err = tx.Where("id IN (?)", ids).Delete(&model.Project{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除项目失败")
	}
	tx.Commit()

	// 删除云平台项目
	for _, project := range projects {
		cloudPid, err := Cloud().GetCloudProjectId(project.Cloud, project.Name)
		if err != nil {
			return err
		}
		if err = Cloud().UpdateCloudProject(project.Cloud, cloudPid, project.Name, 2); err != nil {
			return err
		}
	}

	return err
}

// 项目关联服务器
func (s *ProjectService) UpdateHostAss(param *api.UpdateProjectAssHostReq) (err error) {
	var project model.Project
	var hosts []model.Host
	// 判断所有服务器是否都存在
	if err = util2.CheckIdsExists(model.Host{}, param.Hids); err != nil {
		return err
	}

	if !util2.CheckIdExists(&project, param.Pid) {
		return errors.New("项目ID不存在")
	}

	if err = model.DB.Find(&hosts, param.Hids).Error; err != nil {
		return errors.New("服务器数据库查询操作失败")
	}
	if err = model.DB.First(&project, param.Pid).Error; err != nil {
		return errors.New("项目数据库查询操作失败")
	}
	// host切片所有表Pid改为project.Id
	for _, h := range hosts {
		h.Pid = project.ID
		if err = model.DB.Save(&h).Error; err != nil {
			return fmt.Errorf("保存服务器失败: %v", err)
		}

	}
	if err != nil {
		return err
	}
	return err
}

// 获取项目
func (s *ProjectService) GetProject(param api.SearchIdStringReq) (projectObj any, total int64, err error) {
	var project []model.Project
	db := model.DB.Model(&project)
	if param.Id != 0 {
		if err = db.Where("id = ?", param.Id).Find(&project).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id错误: %v", err)
		}
	} else {
		searchReq := &api.SearchReq{
			Condition: db,
			Table:     &project,
			PageInfo:  param.PageInfo,
		}
		if param.String != "" {
			name := "%" + strings.ToUpper(param.String) + "%"
			db = model.DB.Where("UPPER(name) LIKE ?", name)
			searchReq.Condition = db
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
		} else {
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
		}
	}
	var result *[]api.ProjectRes
	if result, err = s.GetResults(&project); err != nil {
		return nil, 0, err
	}
	return result, total, err
}

// 获取项目关联的服务器
func (s *ProjectService) GetHostAss(param *api.GetPagingMustByIdReq) (hostInfo any, total int64, err error) {
	var project model.Project

	if !util2.CheckIdExists(&project, param.Id) {
		return nil, 0, errors.New("项目ID不存在")
	}

	// 统计被关联个数
	if err = model.DB.Find(&project, param.Id).Error; err != nil {
		return nil, 0, errors.New("查询项目报错")
	}
	// host表中所有pid为project.Id的个数
	if err = model.DB.Model(&model.Host{}).Where("pid = ?", project.ID).Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询关联服务器个数报错")
	}

	// 取出关联数据
	var hosts []model.Host
	offset := (param.PageInfo.Page - 1) * param.PageInfo.PageSize
	if err = model.DB.Where("pid = ?", project.ID).Find(&hosts).Offset(offset).Limit(param.PageInfo.PageSize).Error; err != nil {
		return nil, 0, errors.New("查询关联服务器报错")

	}

	// 分页
	var result *[]api.HostRes
	if result, err = Host().GetResults(&hosts); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 返回项目JSON结果
func (s *ProjectService) GetResults(projectInfo any) (*[]api.ProjectRes, error) {
	var res api.ProjectRes
	var result []api.ProjectRes
	var err error
	if projects, ok := projectInfo.(*[]model.Project); ok {
		for _, project := range *projects {
			res = api.ProjectRes{
				ID:      project.ID,
				Name:    project.Name,
				Cloud:   project.Cloud,
				Status:  project.Status,
				UserId:  project.UserId,
				GroupId: project.GroupId,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if project, ok := projectInfo.(*model.Project); ok {
		res = api.ProjectRes{
			ID:      project.ID,
			Name:    project.Name,
			Cloud:   project.Cloud,
			Status:  project.Status,
			UserId:  project.UserId,
			GroupId: project.GroupId,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换项目结果失败")
}
