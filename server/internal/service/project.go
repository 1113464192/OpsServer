package service

import (
	"errors"
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
func (s *ProjectService) UpdateProject(params *api.UpdateProjectReq) (projectInfo any, err error) {
	var project model.Project
	var count int64
	// 判断项目名是否已被使用
	if model.DB.Model(&project).Where("name = ? AND id != ?", params.Name, params.ID).Count(&count); count > 0 {
		return project, errors.New("项目名已被使用")
	}
	// ID查询
	if params.ID != 0 {
		if !util2.CheckIdExists(&project, params.ID) {
			return project, errors.New("项目不存在")
		}

		if err := model.DB.Where("id = ?", params.ID).First(&project).Error; err != nil {
			return project, errors.New("项目数据库查询失败")
		}

		project.Name = params.Name
		project.Status = params.Status
		project.UserId = params.UserId
		project.GroupId = params.GroupId
		err = model.DB.Save(&project).Error
		if err != nil {
			return project, errors.New("数据保存失败")
		}
		var result *[]api.ProjectRes
		if result, err = s.GetResults(&project); err != nil {
			return nil, err
		}
		return result, err
	} else {
		project = model.Project{
			Name:    params.Name,
			Status:  params.Status,
			UserId:  params.UserId,
			GroupId: params.GroupId,
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
	var project []model.Project
	tx := model.DB.Begin()
	if err = tx.Find(&project, ids).Error; err != nil {
		return errors.New("查询项目信息失败")
	}
	if err = tx.Model(&project).Association("Hosts").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 项目与服务器关联 失败")
	}
	if err = tx.Where("id in (?)", ids).Delete(&model.Project{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除项目失败")
	}
	tx.Commit()
	return err
}

// 项目关联服务器
func (s *ProjectService) UpdateHostAss(params *api.UpdateProjectAssHostReq) (err error) {
	var project model.Project
	var host []model.Host
	// 判断所有项目是否都存在
	if err = util2.CheckIdsExists(model.Host{}, params.Hids); err != nil {
		return err
	}

	if !util2.CheckIdExists(&project, params.Pid) {
		return errors.New("项目ID不存在")
	}

	if err = model.DB.Find(&host, params.Hids).Error; err != nil {
		return errors.New("服务器数据库查询操作失败")
	}
	if err = model.DB.First(&project, params.Pid).Error; err != nil {
		return errors.New("项目数据库查询操作失败")
	}
	if err = model.DB.Model(&project).Association("Hosts").Replace(&host); err != nil {
		return errors.New("项目与服务器数据库关联操作失败")
	}
	if err != nil {
		return err
	}
	return err
}

// 获取项目
func (s *ProjectService) GetProject(params *api.GetProjectReq) (projectObj any, total int64, err error) {
	var project []model.Project
	db := model.DB.Model(&project)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &project,
		PageInfo:  params.PageInfo,
	}
	if params.Name != "" {
		name := "%" + strings.ToUpper(params.Name) + "%"
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
	var result *[]api.ProjectRes
	if result, err = s.GetResults(&project); err != nil {
		return nil, 0, err
	}
	return result, total, err
}

// 获取用户对应项目
func (s *ProjectService) GetSelfProjectList(groupList *[]model.UserGroup, page *api.PageInfo) (projectObj any, total int64, err error) {
	var idList []uint
	var projectList []model.Project
	for _, group := range *groupList {
		idList = append(idList, group.ID)
	}

	db := model.DB.Model(&projectList).Where("group_id IN ?", idList)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &projectList,
		PageInfo:  *page,
	}
	if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
		return nil, 0, err
	}
	var result *[]api.ProjectRes
	if result, err = s.GetResults(&projectList); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 获取项目关联的服务器
func (s *ProjectService) GetHostAss(params *api.GetHostAssReq) (hostInfo any, total int64, err error) {
	var project model.Project
	if !util2.CheckIdExists(&project, params.ProjectId) {
		return nil, 0, errors.New("项目ID不存在")
	}
	if err = model.DB.Preload("Hosts").Where("id = ?", params.ProjectId).First(&project).Error; err != nil {
		return nil, 0, errors.New("项目查询失败")
	}
	assQueryReq := &api.AssQueryReq{
		Condition: model.DB.Model(&model.Host{}),
		Table:     &project.Hosts,
		PageInfo:  params.PageInfo,
	}
	if total, err = dbOper.DbOper().AssDbFind(assQueryReq); err != nil {
		return nil, 0, err
	}
	var result *[]api.HostRes
	if result, err = Host().GetResults(&project.Hosts); err != nil {
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
			Status:  project.Status,
			UserId:  project.UserId,
			GroupId: project.GroupId,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换项目结果失败")
}
