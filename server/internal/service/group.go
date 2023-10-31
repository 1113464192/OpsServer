package service

import (
	"database/sql"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils2"
	"strings"
)

type GroupService struct {
}

var (
	insGroup = &GroupService{}
)

func Group() *GroupService {
	return insGroup
}

// 修改或新增用户组
func (s *GroupService) UpdateGroup(params *api.UpdateGroupReq) (groupInfo any, err error) {
	var group model.UserGroup
	var count int64
	if model.DB.Model(&group).Where("name = ? AND id != ?", params.Name, params.ID).Count(&count); count > 0 {
		return group, errors.New("组名已被使用")
	}
	if params.ID != 0 {
		// 修改
		if !utils2.CheckIdExists(&group, &params.ID) {
			return group, errors.New("用户组不存在")
		}

		if err := model.DB.Where("id = ?", params.ID).First(&group).Error; err != nil {
			return group, errors.New("用户组数据库查询失败")
		}

		group.Name = params.Name
		group.ParentId = params.ParentId
		if params.Mark == "" {
			group.Mark = sql.NullString{String: "", Valid: false}
		} else {
			group.Mark = sql.NullString{String: params.Mark, Valid: true}
		}

		err = model.DB.Save(&group).Error
		if err != nil {
			return group, errors.New("数据保存失败")
		}
		var result []api.GroupRes
		if result, err = s.GetResults(&group); err != nil {
			return nil, err
		}
		return result, err
	} else {
		group = model.UserGroup{
			Name:     params.Name,
			ParentId: params.ParentId,
		}
		if params.Mark == "" {
			group.Mark = sql.NullString{String: "", Valid: false}
		} else {
			group.Mark = sql.NullString{String: params.Mark, Valid: true}
		}
		if err = model.DB.Create(&group).Error; err != nil {
			logger.Log().Error("Group", "创建用户组失败", err)
			return group, errors.New("创建用户组失败")
		}
		var result []api.GroupRes
		if result, err = s.GetResults(&group); err != nil {
			return nil, err
		}
		return result, err
	}
}

// 关联用户
func (s *GroupService) UpdateUserAss(params *api.UpdateUserAssReq) (err error) {
	var group model.UserGroup
	var user []model.User
	// 判断用户是否都存在
	var noExistId []uint
	for _, uid := range params.UserIDs {
		uBool := utils2.CheckIdExists(&user, &uid)
		if !uBool {
			noExistId = append(noExistId, uid)
		}
	}
	if len(noExistId) != 0 {
		return fmt.Errorf("%v %s", noExistId, "用户不存在")
	}

	tx := model.DB.Begin()
	if err := tx.First(&group, params.GroupID).Error; err != nil {
		tx.Rollback()
		return errors.New("用户组不存在")
	}
	if err := tx.Find(&user, params.UserIDs).Error; err != nil {
		tx.Rollback()
		return errors.New("用户不存在")
	}
	if err := tx.Model(&group).Association("Users").Replace(&user); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

// 删除用户组
func (s *GroupService) DeleteUserGroup(ids []uint) (err error) {
	for _, i := range ids {
		if !utils2.CheckIdExists(&model.UserGroup{}, &i) {
			return errors.New("用户组不存在")
		}
	}
	var group []model.UserGroup
	tx := model.DB.Begin()
	if err = tx.Find(&group, ids).Error; err != nil {
		return errors.New("查询用户组信息失败")
	}
	if err = tx.Model(&group).Association("Users").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 用户组与用户关联 失败")
	}
	if err = tx.Model(&group).Association("Menus").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 用户组与菜单关联 失败")
	}
	if err = tx.Where("id in (?)", ids).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除用户失败")
	}
	tx.Commit()
	return err
}

// 获取用户组
func (s *GroupService) GetGroupList(params *api.GetGroupReq) (groupObj any, total int64, err error) {
	var group []model.UserGroup
	db := model.DB.Model(&group)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &group,
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
	var result []api.GroupRes
	if result, err = s.GetResults(&group); err != nil {
		return nil, 0, err
	}
	return result, total, err
}

// 获取用户组对应用户
func (s *GroupService) GetAssUser(params *api.GetGroupAssIdReq) (userObj any, total int64, err error) {
	var group model.UserGroup
	if !utils2.CheckIdExists(&group, &params.Id) {
		return nil, 0, errors.New("组ID不存在")
	}
	if err = model.DB.Preload("Users").Where("id = ?", params.Id).First(&group).Error; err != nil {
		return nil, 0, errors.New("组查询失败")
	}
	assQueryReq := &api.AssQueryReq{
		Condition: model.DB.Model(&model.UserGroup{}),
		Table:     &group.Users,
		PageInfo:  params.PageInfo,
	}
	if total, err = dbOper.DbOper().AssDbFind(assQueryReq); err != nil {
		return nil, 0, err
	}
	var result []api.UserRes
	if result, err = User().GetResults(&group.Users); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 获取用户组对应项目
func (s *GroupService) GetAssProject(params *api.GetGroupAssIdReq) (projectObj any, total int64, err error) {
	var group model.UserGroup
	if !utils2.CheckIdExists(&group, &params.Id) {
		return nil, 0, errors.New("组ID不存在")
	}
	if err = model.DB.Preload("Users").Where("id = ?", params.Id).First(&group).Error; err != nil {
		return nil, 0, errors.New("组查询失败")
	}
	assQueryReq := &api.AssQueryReq{
		Condition: model.DB.Model(&model.UserGroup{}),
		Table:     &group.Project,
		PageInfo:  params.PageInfo,
	}
	if total, err = dbOper.DbOper().AssDbFind(assQueryReq); err != nil {
		return nil, 0, err
	}
	var result []api.ProjectRes
	if result, err = Project().GetResults(&group.Project); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 返回用户组结果
func (s *GroupService) GetResults(groupInfo any) (result []api.GroupRes, err error) {
	var res api.GroupRes
	if groups, ok := groupInfo.(*[]model.UserGroup); ok {
		for _, group := range *groups {
			res = api.GroupRes{
				ID:       group.ID,
				Name:     group.Name,
				ParentId: group.ParentId,
				Mark:     group.Mark.String,
			}
			result = append(result, res)
		}
		return result, err
	}
	if group, ok := groupInfo.(*model.UserGroup); ok {
		res = api.GroupRes{
			ID:       group.ID,
			Name:     group.Name,
			ParentId: group.ParentId,
			Mark:     group.Mark.String,
		}
		result = append(result, res)
		return result, err
	}
	return result, errors.New("转换组结果失败")
}
