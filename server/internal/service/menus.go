package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util2"
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
func (s *MenuService) UpdateMenu(param *api.UpdateMenuReq) (menuInfo any, err error) {
	var menu model.Menus
	var count int64
	// 判断菜单名/字段是否被占用
	if model.DB.Model(&menu).Where("name = ? AND id != ?", param.Name, param.ID).Or("title = ? AND id != ?", param.Title, param.ID).Count(&count); count > 0 {
		return menu, errors.New("菜单Name字段或菜单Title字段已被使用")
	}
	if param.ID != 0 {
		// 判断菜单是否存在
		if !util2.CheckIdExists(&menu, param.ID) {
			return menu, errors.New("该菜单不存在")
		}

		if err := model.DB.Where("id = ?", param.ID).Find(&menu).Error; err != nil {
			return menu, errors.New("用户组数据库查询失败")
		}

		menu.ID = param.ID
		menu.Name = param.Name
		menu.ParentId = param.ParentId
		menu.Mark = param.Mark
		menu.Type = param.Type
		menu.Title = param.Title
		menu.Url = param.Url
		menu.Sort = param.Sort
		menu.Icon = param.Icon
		menu.Author = param.Author
		menu.Component = param.Component

		if err = model.DB.Save(&menu).Error; err != nil {
			return menu, fmt.Errorf("数据保存失败: %v", err)
		}
		return menu, err
	} else {
		menu = model.Menus{
			Name:      param.Name,
			ParentId:  param.ParentId,
			Mark:      param.Mark,
			Type:      param.Type,
			Title:     param.Title,
			Url:       param.Url,
			Sort:      param.Sort,
			Icon:      param.Icon,
			Author:    param.Author,
			Component: param.Component,
		}
		if err = model.DB.Create(&menu).Error; err != nil {
			logger.Log().Error("Menu", "创建菜单失败", err)
			return menu, errors.New("创建菜单失败")
		}
		return menu, err
	}
}

// 关联用户组
func (s *MenuService) UpdateMenuAss(param *api.UpdateMenuAssReq) (menuObj any, err error) {
	var menu model.Menus
	var groups []model.UserGroup
	// 默认添加管理组
	// if !util.IsSliceContain(param.GroupIDs, 1) {
	// 	param.GroupIDs = append(param.GroupIDs, 1)
	// }
	if err = util2.CheckIdsExists(model.UserGroup{}, param.GroupIDs); err != nil {
		return nil, err
	}

	tx := model.DB.Begin()
	if err := tx.First(&menu, param.MenuID).Error; err != nil {
		tx.Rollback()
		return menu, errors.New("菜单不存在")
	}

	if err := tx.Find(&groups, param.GroupIDs).Error; err != nil {
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

// 获取用户组对应菜单
func (s *MenuService) GetMenuList(gid uint, isAdmin uint8) (menu *[]model.Menus, err error) {
	var group []model.UserGroup
	if gid != 0 {
		if err = model.DB.First(&group, gid).Error; err != nil {
			return menu, fmt.Errorf("查找gid的用户组失败: %v", err)
		}
		if err = model.DB.Model(&group).Association("Menus").Find(&menu); err != nil {
			return menu, fmt.Errorf("查找用户组关联的菜单失败: %v", err)
		}
		return menu, err
	}
	if isAdmin == 1 {
		if err := model.DB.Find(&menu).Error; err != nil {
			return menu, err
		}
		return menu, err
	}
	return nil, errors.New("如果不是管理员, 请输入菜单ID")
}

// 删除菜单
func (s *MenuService) DeleteMenu(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.UserGroup{}, ids); err != nil {
		return err
	}
	var menu []model.Menus
	tx := model.DB.Begin()
	if err = tx.Find(&menu, ids).Error; err != nil {
		return errors.New("查询菜单信息失败")
	}
	if err = tx.Model(&menu).Association("UserGroup").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 菜单与用户组关联 失败")
	}

	if err = tx.Where("id in (?)", ids).Delete(&model.Menus{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除菜单失败")
	}
	tx.Commit()
	return err
}
