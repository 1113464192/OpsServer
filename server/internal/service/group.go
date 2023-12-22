package service

import (
	"database/sql"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util2"
	"strconv"
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

func (s *GroupService) permissionInheritance(parentId uint, group *model.UserGroup) (err error) {
	// 顶层组无权限继承
	if parentId == 0 {
		return nil
	}
	var parentGroup model.UserGroup
	if err = model.DB.First(&parentGroup, parentId).Error; err != nil {
		return fmt.Errorf("查询用户组报错: %v", err)
	}
	// 继承工作室菜单
	var parentMenus []model.Menus
	if err = model.DB.Model(&parentGroup).Association("Menus").Find(&parentMenus); err != nil {
		return fmt.Errorf("查询工作室菜单报错: %v", err)
	}
	if err = model.DB.Model(&group).Association("Menus").Replace(&parentMenus); err != nil {
		return fmt.Errorf("继承工作室菜单报错: %v", err)
	}

	// 继承工作室API
	var anySlice []any
	var parentApiIds []uint
	if anySlice, err = CasbinServiceApp().GetPolicyPathByGroupIds([]uint{parentId}); err != nil {
		return fmt.Errorf("获取工作室的APIIDs失败: %v", err)
	}
	for _, anyUint := range anySlice {
		parentApiId, ok := anyUint.(uint)
		if !ok {
			return fmt.Errorf("转换工作室的APIIDs失败: %v", err)
		}
		parentApiIds = append(parentApiIds, parentApiId)
	}
	updateCasbinReq := api.UpdateCasbinReq{
		GroupId: strconv.Itoa(int(group.ID)),
		Ids:     parentApiIds,
	}
	if err = CasbinServiceApp().UpdateCasbin(updateCasbinReq); err != nil {
		return fmt.Errorf("继承工作室API报错: %v", err)
	}
	return err
}

// 修改或新增用户组
func (s *GroupService) UpdateGroup(param *api.UpdateGroupReq) (groupInfo any, err error) {
	var group model.UserGroup
	var count int64
	// 判断组是否被占用
	if model.DB.Model(&group).Where("name = ? AND id != ?", param.Name, param.ID).Count(&count); count > 0 {
		return group, errors.New("组名已被使用")
	}
	// 根据ID查询组
	if param.ID != 0 {
		// 判断组是否存在
		if !util2.CheckIdExists(&group, param.ID) {
			return group, errors.New("用户组不存在")
		}
		// 获取组对象
		if err := model.DB.Where("id = ?", param.ID).First(&group).Error; err != nil {
			return group, errors.New("用户组数据库查询失败")
		}

		group.Name = param.Name
		group.ParentId = param.ParentId

		// 有标识则写入，没有默认为Null
		if param.Mark != "" {
			group.Mark = sql.NullString{String: param.Mark, Valid: true}
		}
		// 入库
		if err = model.DB.Save(&group).Error; err != nil {
			return group, fmt.Errorf("数据保存失败: %v", err)
		}
		// 过滤结果
		var result *[]api.GroupRes
		if result, err = s.GetResults(&group); err != nil {
			return nil, err
		}
		return result, err
	} else {
		group = model.UserGroup{
			Name:     param.Name,
			ParentId: param.ParentId,
		}
		if param.Mark != "" {
			group.Mark = sql.NullString{String: param.Mark, Valid: true}
		}
		if err = model.DB.Create(&group).Error; err != nil {
			logger.Log().Error("Group", "创建用户组失败", err)
			return group, fmt.Errorf("创建用户组失败: %v", err)
		}
		if err = s.permissionInheritance(uint(group.ParentId), &group); err != nil {
			return group, fmt.Errorf("创建组成功，但是继承工作室权限失败: %v", err)
		}
		var result *[]api.GroupRes
		if result, err = s.GetResults(&group); err != nil {
			return nil, err
		}
		return result, err
	}
}

// 关联用户
func (s *GroupService) UpdateGroupAssUser(param *api.UpdateGroupAssUserReq) (err error) {
	var group model.UserGroup
	var user []model.User
	// 判断用户是否都存在
	if err = util2.CheckIdsExists(user, param.UserIDs); err != nil {
		return err
	}

	tx := model.DB.Begin()
	if err := tx.First(&group, param.GroupID).Error; err != nil {
		tx.Rollback()
		return errors.New("用户组不存在")
	}
	if err := tx.Find(&user, param.UserIDs).Error; err != nil {
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
	if err = util2.CheckIdsExists(model.UserGroup{}, ids); err != nil {
		return err
	}
	var group []model.UserGroup
	tx := model.DB.Begin()
	if err = tx.Find(&group, ids).Error; err != nil {
		return errors.New("查询用户组信息失败")
	}
	// 如果删除工作室则判断是否存在未删除的项目组
	var count int64
	for _, g := range group {
		if g.ParentId == 0 {
			if err = tx.Model(&model.UserGroup{}).Where("parent_id = ?", g.ID).Count(&count).Error; err != nil || count > 0 {
				tx.Rollback()
				if err != nil {
					return fmt.Errorf("查询用户组下是否有子组失败: %v", err)
				}
				return fmt.Errorf("%d 用户组下有子组存在", g.ID)
			}
		}
	}
	if err = tx.Model(&group).Association("Users").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 用户组与用户关联 失败")
	}
	if err = tx.Model(&group).Association("Menus").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 用户组与菜单关联 失败")
	}
	if err = tx.Where("id IN (?)", ids).Delete(&model.UserGroup{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除用户失败")
	}
	tx.Commit()
	return err
}

// 获取用户组
func (s *GroupService) GetGroupList(param *api.SearchIdStringReq) (groupObj any, total int64, err error) {
	var group []model.UserGroup
	db := model.DB.Model(&group)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &group,
		PageInfo:  param.PageInfo,
	}
	if param.Id != 0 {
		if err = db.Where("id = ?", param.Id).Count(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id总数错误: %v", err)
		}
		if err = db.Where("id = ?", param.Id).Find(&group).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id错误: %v", err)
		}
	} else {
		if param.String != "" {
			name := "%" + strings.ToUpper(param.String) + "%"
			db = model.DB.Model(&group).Where("UPPER(name) LIKE ?", name)
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

	var result *[]api.GroupRes
	if result, err = s.GetResults(&group); err != nil {
		return nil, 0, err
	}
	return result, total, err
}

// 获取用户组对应用户
func (s *GroupService) GetAssUser(param *api.GetPagingMustByIdsReq) (userObj any, total int64, err error) {
	var groups []model.UserGroup
	if err = util2.CheckIdsExists(&groups, param.Ids); err != nil {
		return nil, 0, err
	}

	// 获取组数据
	if err = model.DB.Find(&groups, param.Ids).Error; err != nil {
		return nil, 0, errors.New("查询用户组报错")
	}

	// 统计被关联个数
	if total = model.DB.Model(&groups).Association("Users").Count(); total == 0 {
		return "没有关联数据", 0, nil
	}

	// 取出关联数据
	var users []model.User
	if err = model.DB.Model(&groups).Order("id asc").Association("Users").Find(&users); err != nil {
		return &users, total, fmt.Errorf("获取关联的数据失败: %v", err)
	}

	// 取出所有关联的数据并去重
	var dedupliusers []model.User
	userMap := make(map[uint]struct{})
	for _, user := range users {
		if _, ok := userMap[user.ID]; !ok {
			dedupliusers = append(dedupliusers, user)
			userMap[user.ID] = struct{}{}
		}
	}

	// 分页
	if err = dbOper.DbOper().PaginateModels(&dedupliusers, param.PageInfo); err != nil {
		return nil, 0, err
	}

	// 过滤结果
	var result *[]api.UserRes
	if result, err = User().GetResults(&dedupliusers); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 获取用户组对应项目
func (s *GroupService) GetAssProject(param *api.GetPagingMustByIdsReq) (result *[]api.ProjectRes, total int64, err error) {
	var group model.UserGroup
	// 判断是否有不存在的ID
	if err = util2.CheckIdsExists(group, param.Ids); err != nil {
		return nil, 0, err
	}
	var projects []model.Project
	db := model.DB.Model(&projects).Where("group_id IN (?)", param.Ids)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &projects,
		PageInfo:  param.PageInfo,
	}
	if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
		return nil, 0, err
	}
	if result, err = Project().GetResults(&projects); err != nil {
		return nil, total, err
	}
	return result, total, err
}

// 关联菜单
func (s *MenuService) UpdateGroupAssMenus(param *api.UpdateGroupAssMenusReq) (any, error) {
	var err error
	var menus []model.Menus
	var group model.UserGroup
	// 默认添加管理组
	// if !util.IsSliceContain(param.GroupIDs, 1) {
	// 	param.GroupIDs = append(param.GroupIDs, 1)
	// }
	if err = util2.CheckIdsExists(model.Menus{}, param.MenuIDs); err != nil {
		return nil, err
	}

	tx := model.DB.Begin()
	if err := tx.Find(&group, param.GroupID).Error; err != nil {
		tx.Rollback()
		return group, errors.New("用户组不存在")
	}

	if err := tx.Find(&menus, param.MenuIDs).Error; err != nil {
		tx.Rollback()
		return menus, errors.New("菜单不存在")
	}

	if err := tx.Model(&group).Association("Menus").Replace(&menus); err != nil {
		tx.Rollback()
		return group, err
	}
	tx.Commit()

	return menus, err
}

// 获取用户组对应菜单
func (s *GroupService) GetGroupAssMenus(param *api.GetPagingMustByIdsReq) (MenuObj any, total int64, err error) {
	var groups []model.UserGroup
	if err = util2.CheckIdsExists(&groups, param.Ids); err != nil {
		return nil, 0, err
	}
	// 获取组数据
	if err = model.DB.Find(&groups, param.Ids).Error; err != nil {
		return nil, 0, errors.New("查询用户组报错")
	}
	// 统计被关联个数
	if total = model.DB.Model(&groups).Association("Menus").Count(); total == 0 {
		return "没有关联数据", 0, nil
	}
	// 取出关联数据
	var menus []model.Menus
	if err = model.DB.Model(&groups).Order("id asc").Association("Menus").Find(&menus); err != nil {
		return &menus, total, fmt.Errorf("获取关联的数据失败: %v", err)
	}

	// 取出所有预加载的表并去重
	var deduplimenus []model.Menus
	menuMap := make(map[uint]struct{})
	for _, menu := range menus {
		if _, ok := menuMap[menu.ID]; !ok {
			deduplimenus = append(deduplimenus, menu)
			menuMap[menu.ID] = struct{}{}
		}
	}

	// 分页
	if err = dbOper.DbOper().PaginateModels(&deduplimenus, param.PageInfo); err != nil {
		return nil, 0, err
	}
	return &deduplimenus, total, err
}

// 返回用户组结果
func (s *GroupService) GetResults(groupInfo any) (*[]api.GroupRes, error) {
	var res api.GroupRes
	var result []api.GroupRes
	var err error
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
		return &result, err
	}
	if group, ok := groupInfo.(*model.UserGroup); ok {
		res = api.GroupRes{
			ID:       group.ID,
			Name:     group.Name,
			ParentId: group.ParentId,
			Mark:     group.Mark.String,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换组结果失败")
}
