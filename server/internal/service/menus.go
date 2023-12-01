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

// 获取菜单
func (s *MenuService) GetMenuList(param api.SearchIdStringReq) (*[]model.Menus, int64, error) {
	var err error
	var total int64
	var menus []model.Menus
	db := model.DB.Model(&menus)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &menus,
		PageInfo:  param.PageInfo,
	}
	if param.Id != 0 {
		if err = db.Where("id IN (?)", param.Id).Count(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("查询ids总数错误: %v", err)
		}
		if err = db.Where("id IN (?)", param.Id).Find(&menus).Error; err != nil {
			return nil, 0, fmt.Errorf("查询ids错误: %v", err)
		}
	} else {
		if param.String != "" {
			title := "%" + strings.ToUpper(param.String) + "%"
			db = model.DB.Where("UPPER(title) LIKE ?", title)
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

	// var result *[]api.GroupRes
	// if result, err = s.get(&group); err != nil {
	// 	return nil, 0, err
	// }
	return &menus, total, err
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

	if err = tx.Where("id IN (?)", ids).Delete(&model.Menus{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除菜单失败")
	}
	tx.Commit()
	return err
}
