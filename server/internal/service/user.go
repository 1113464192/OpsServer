package service

import (
	"errors"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/utils"
	"fqhWeb/pkg/utils/jwt"
	"fqhWeb/pkg/utils2"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

type UserService struct{}

var insUser UserService

func User() *UserService {
	return &insUser
}

// 修改/添加用户
func (s *UserService) UpdateUser(params *api.UpdateUserReq) (userInfo any, err error) {
	if params.Mobile != "" && !utils.CheckMobile(params.Mobile) {
		return params.Mobile, errors.New("电话格式错误")
	}

	if params.Email != "" && !utils.CheckEmail(params.Email) {
		return params.Email, errors.New("邮箱格式错误")
	}
	var user model.User
	var count int64
	if params.ID != 0 {
		// 修改
		if !utils2.CheckIdExists(&user, &params.ID) {
			return nil, errors.New("用户不存在")
		}

		if err := model.DB.Where("id = ?", params.ID).Find(&user).Error; err != nil {
			return nil, errors.New("用户数据库查询失败")
		}
		user.Username = params.Username
		user.Name = params.Name
		user.Expiration = params.Expiration
		user.Mobile = params.Mobile
		user.Email = params.Email

		if model.DB.Model(&user).Where("username = ? AND id != ?", params.Username, params.ID).Count(&count); count > 0 {
			return nil, errors.New("用户名已被使用")
		}

		err = model.DB.Save(&user).Error
		if err != nil {
			return nil, errors.New("数据保存失败")
		}
		var result []api.UserRes
		result, err = s.GetResults(&user)
		if err != nil {
			return nil, err
		}
		return result, err
	} else {
		model.DB.Model(&user).Where("username = ?", params.Username).Count(&count)
		if count > 0 {
			return user, errors.New("账号已经注册")
		}
		user = model.User{
			Username:   params.Username,
			Name:       params.Name,
			Expiration: params.Expiration,
			Mobile:     params.Mobile,
			Email:      params.Email,
			IsAdmin:    params.IsAdmin,
		}
		password := utils.RandStringRunes(12)
		user.Password, err = utils.GenerateFromPassword(password)
		if err != nil {
			return user, errors.New("用户密码bcrypt加密失败")
		}
		if err = model.DB.Create(&user).Error; err != nil {
			return user, errors.New("创建用户失败")
		}
		var result []api.UserRes
		result, err = s.GetResults(&user)
		if err != nil {
			return nil, err
		}
		return result, err
	}
}

// 获取用户切片
func (s *UserService) GetUserList(params api.PageInfo, username string) (list any, total int64, err error) {
	var user []model.User
	db := model.DB.Model(&user)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &user,
		PageInfo:  params,
	}
	if username != "" {
		name := "%" + strings.ToUpper(username) + "%"
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
	var result []api.UserRes
	result, err = s.GetResults(&user)
	if err != nil {
		return nil, 0, err
	}
	return result, total, err
}

// 删除用户
func (s *UserService) DeleteUser(ids []uint) (err error) {
	for _, i := range ids {
		if !utils2.CheckIdExists(&model.User{}, &i) {
			return errors.New("用户不存在")
		}
	}
	var user []model.User
	tx := model.DB.Begin()
	if err = tx.Find(&user, ids).Error; err != nil {
		return errors.New("查询用户信息失败")
	}
	if err = tx.Model(&user).Association("UserGroups").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 用户与用户组关联 失败")
	}
	if err = tx.Where("id in (?)", ids).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除用户失败")
	}
	tx.Commit()
	return err
}

// 修改用户状态
func (s *UserService) UpdateStatus(params *api.StatusReq) (err error) {
	if !utils2.CheckIdExists(&model.User{}, &params.ID) {
		return errors.New("用户不存在")
	}
	err = model.DB.Model(&model.User{}).Where("id = ?", params.ID).Update("status", params.Status).Error
	return err
}

// 修改密码
func (s *UserService) UpdatePasswd(passwd *api.PasswordReq) (err error) {
	if !utils2.CheckIdExists(&model.User{}, &passwd.ID) {
		return errors.New("用户不存在")
	}
	err = model.DB.Model(&model.User{}).Where("id = ?", passwd.ID).Update("password", passwd.Password).Error
	return err
}

// 获取用户个人信息
func (s *UserService) GetSelfInfo(id *uint) (userInfo any, err error) {
	var user model.User
	if utils2.CheckIdExists(&model.User{}, id) {
		err = model.DB.Where("id in (?)", *id).First(&user).Error
		if err != nil {
			return user, errors.New("查询用户个人信息失败")
		}
		var result []api.UserRes
		result, err = s.GetResults(&user)
		if err != nil {
			return nil, err
		}
		return result, err
	} else {
		return nil, errors.New("用户不存在")
	}

}

// 获取用户关联组信息
func (s *UserService) GetAssGroup(id uint) (groupInfo any, err error) {
	var user model.User
	var group []model.UserGroup
	if err != nil {
		return nil, err
	}
	if utils2.CheckIdExists(&model.User{}, &id) {
		if err := model.DB.First(&user, id).Error; err != nil {
			return user, err
		}
		if err := model.DB.Model(&user).Association("UserGroups").Find(&group); err != nil {
			return user, err
		}

		return group, err
	} else {
		return nil, errors.New("用户不存在")
	}

}

// 登录
func (s *UserService) Login(u *model.User) (userInfo *api.AuthLoginRes, err error) {
	var user model.User
	if err = model.DB.Where("username = ?", u.Username).Preload("UserGroups").First(&user).Error; err != nil {
		return userInfo, errors.New("获取用户对象失败")
	}
	if !utils.CheckPassword(user.Password, u.Password) {
		return userInfo, errors.New("账号或密码错误")
	}
	if user.Status != 1 {
		return userInfo, errors.New("账号不存在或账号已被禁用")
	}
	if uint64(time.Now().Local().Sub(user.CreatedAt).Seconds()) > user.Expiration {
		return userInfo, errors.New("账号已过期，请联系管理员延长过期时间")
	}

	// 设置JWT-token
	token, err := jwt.GenToken(user)
	if err != nil {
		return userInfo, err
	}
	var groupIds []uint
	for _, userGroup := range user.UserGroups {
		groupIds = append(groupIds, userGroup.ID)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("login_time", time.Now().Local().Unix()).Error; err != nil {
		return nil, errors.New("写入登陆时间失败")
	}
	userInfo = &api.AuthLoginRes{
		ID:       user.ID,
		Username: user.Username,
		GroupIds: groupIds,
		Name:     user.Name,
		Token:    token,
	}
	return userInfo, err
}

// 通过文件更新私钥
func (s *UserService) UpdateKeyFileContext(file *multipart.FileHeader, keyPasswd string, id uint) error {
	fileP, err := file.Open()
	if err != nil {
		return err
	}
	defer fileP.Close()

	fileBytes, err := io.ReadAll(fileP)
	if err != nil {
		return err
	}

	fileContent := string(fileBytes)
	err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("pri_key", fileContent).Error
	if err != nil {
		return errors.New("私钥写入数据库失败")
	}
	err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("key_passwd", keyPasswd).Error
	if err != nil {
		return errors.New("通行证密码写入数据库失败")
	}
	return nil
}

// 通过字符串更新私钥内容
func (s *UserService) UpdateKeyContext(key string, keyPasswd string, id uint) (err error) {
	err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("pri_key", key).Error
	if err != nil {
		return errors.New("私钥字符串写入数据库失败")
	}
	err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("key_passwd", keyPasswd).Error
	if err != nil {
		return errors.New("通行证密码写入数据库失败")
	}
	return nil
}

// 返回用户结果
func (s *UserService) GetResults(userInfo any) (result []api.UserRes, err error) {
	var res api.UserRes
	if users, ok := userInfo.(*[]model.User); ok {
		for _, user := range *users {
			res = api.UserRes{
				ID:         user.ID,
				Name:       user.Name,
				Username:   user.Username,
				Password:   user.Password,
				Status:     user.Status,
				Email:      user.Email,
				Mobile:     user.Mobile,
				LoginTime:  uint64(user.LoginTime),
				Expiration: user.Expiration,
				IsAdmin:    user.IsAdmin,
			}
			result = append(result, res)
		}
		return result, err
	}
	if user, ok := userInfo.(*model.User); ok {
		res = api.UserRes{
			ID:         user.ID,
			Name:       user.Name,
			Username:   user.Username,
			Password:   user.Password,
			Status:     user.Status,
			Email:      user.Email,
			Mobile:     user.Mobile,
			LoginTime:  uint64(user.LoginTime),
			Expiration: user.Expiration,
			IsAdmin:    user.IsAdmin,
		}
		result = append(result, res)
		return result, err
	}
	return result, errors.New("转换用户结果失败")
}
