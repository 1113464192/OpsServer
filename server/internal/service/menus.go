package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils"
	"fqhWeb/pkg/utils2"
)

type MenuService struct {
}

var (
	insMenu = &MenuService{}
)

func Menu() *MenuService {
	return insMenu
}

// UpdateMenu
func (s *MenuService) UpdateMenu(params *api.UpdateMenuReq) (menuInfo any, err error) {
	var menu model.Menus
	var count int64
	if model.DB.Model(&menu).Where("name = ? AND id != ?", params.Name, params.ID).Or("title = ? AND id != ?", params.Title, params.ID).Count(&count); count > 0 {
		return menu, errors.New("菜单Name字段或菜单Title字段已被使用")
	}
	if params.ID != 0 {
		// 修改
		if !utils2.CheckIdExists(&menu, &params.ID) {
			return menu, errors.New("该菜单不存在")
		}

		if err := model.DB.Where("id = ?", params.ID).Find(&menu).Error; err != nil {
			return menu, errors.New("用户组数据库查询失败")
		}

		menu.ID = params.ID
		menu.Name = params.Name
		menu.ParentId = params.ParentId
		menu.Mark = params.Mark
		menu.Type = params.Type
		menu.Title = params.Title
		menu.Url = params.Url
		menu.Sort = params.Sort
		menu.Icon = params.Icon
		menu.Author = params.Author
		menu.Component = params.Component

		err = model.DB.Save(&menu).Error
		if err != nil {
			return menu, errors.New("数据保存失败")
		}
		return menu, err
	} else {
		menu = model.Menus{
			Name:      params.Name,
			ParentId:  params.ParentId,
			Mark:      params.Mark,
			Type:      params.Type,
			Title:     params.Title,
			Url:       params.Url,
			Sort:      params.Sort,
			Icon:      params.Icon,
			Author:    params.Author,
			Component: params.Component,
		}
		if err = model.DB.Create(&menu).Error; err != nil {
			logger.Log().Error("Menu", "创建菜单失败", err)
			return menu, errors.New("创建菜单失败")
		}
		return menu, err
	}
}

// 关联用户组
func (s *MenuService) UpdateMenuAss(params *api.UpdateMenuAssReq) (menuObj any, err error) {
	var menu model.Menus
	var groups []model.UserGroup
	// 默认添加管理组
	if !utils.IsContain(params.GroupIDs, 1) {
		params.GroupIDs = append(params.GroupIDs, 1)
	}
	var noExistId []uint

	// 判断用户组是否都存在
	for _, gid := range params.GroupIDs {
		uBool := utils2.CheckIdExists(&groups, &gid)
		if !uBool {
			noExistId = append(noExistId, gid)
		}
	}
	if len(noExistId) != 0 {
		return groups, fmt.Errorf("%v %s", noExistId, "用户组不存在")
	}

	tx := model.DB.Begin()
	if err := tx.First(&menu, params.MenuID).Error; err != nil {
		tx.Rollback()
		return menu, errors.New("菜单不存在")
	}

	if err := tx.Find(&groups, params.GroupIDs).Error; err != nil {
		tx.Rollback()
		return groups, errors.New("用户不存在")
	}

	if err := tx.Model(&menu).Association("UserGroups").Replace(&groups); err != nil {
		tx.Rollback()
		return menu, err
	}
	tx.Commit()

	return menu, err
}

// 获取菜单对应用户组
// 后续优化公共关联方法
func (s *MenuService) GetMenuList(groupIdStr *string) (menu *[]model.Menus, err error) {
	var group []model.UserGroup
	if *groupIdStr != "" {
		var gid uint
		gid, err = utils.StringToUint(groupIdStr)
		if err != nil {
			return nil, err
		}
		if err := model.DB.Model(&group).Association("Menus").Find(&menu, gid); err != nil {
			return menu, err
		}
		return menu, err
	} else {
		if err := model.DB.Find(&menu).Error; err != nil {
			return menu, err
		}
		return menu, err
	}
}
