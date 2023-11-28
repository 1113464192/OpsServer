package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/util"
	"fqhWeb/pkg/util/jwt"
	"fqhWeb/pkg/util2"
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
func (s *UserService) UpdateUser(param *api.UpdateUserReq) (any, string, error) {
	var err error
	// 判断电话和邮箱是否正确
	if param.Mobile != "" && !util.CheckMobile(param.Mobile) {
		return param.Mobile, "", errors.New("电话格式错误")
	}

	if param.Email != "" && !util.CheckEmail(param.Email) {
		return param.Email, "", errors.New("邮箱格式错误")
	}

	var user model.User
	var count int64
	if param.ID != 0 {
		// 修改
		if !util2.CheckIdExists(&user, param.ID) {
			return nil, "", errors.New("用户不存在")
		}

		if err := model.DB.Where("id = ?", param.ID).Find(&user).Error; err != nil {
			return nil, "", errors.New("用户数据库查询失败")
		}
		user.Username = param.Username
		user.Name = param.Name
		user.Expiration = param.Expiration
		user.Mobile = param.Mobile
		user.Email = param.Email

		// 判断username是否和现有用户重复
		if model.DB.Model(&user).Where("username = ? AND id != ?", param.Username, param.ID).Count(&count); count > 0 {
			return nil, "", errors.New("用户名已被使用")
		}

		if err = model.DB.Save(&user).Error; err != nil {
			return nil, "", fmt.Errorf("数据保存失败: %v", err)
		}

		// 返回过滤后的JSON
		var result *[]api.UserRes
		result, err = s.GetResults(&user)
		if err != nil {
			return nil, "", err
		}
		return result, "", err
	} else {
		// 判断username是否和现有用户重复
		if model.DB.Model(&user).Where("username = ?", param.Username).Count(&count); count > 0 {
			return user, "", errors.New("账号已经注册")
		}
		user = model.User{
			Username:   param.Username,
			Name:       param.Name,
			Expiration: param.Expiration,
			Mobile:     param.Mobile,
			Email:      param.Email,
			IsAdmin:    param.IsAdmin,
		}
		// 生成初始化密码
		password := util.RandStringRunes(12)
		user.Password, err = util.GenerateFromPassword(password)
		if err != nil {
			return user, "", errors.New("用户密码bcrypt加密失败")
		}
		if err = model.DB.Create(&user).Error; err != nil {
			return user, "", errors.New("创建用户失败")
		}
		var result *[]api.UserRes
		result, err = s.GetResults(&user)
		if err != nil {
			return nil, "", err
		}
		// 创建用户时返回初始密码
		return result, password, err
	}
}

// 获取用户切片
func (s *UserService) GetUserList(param api.SearchStringReq) (list any, total int64, err error) {
	var user []model.User
	db := model.DB.Model(&user)
	// 有ID优先ID
	if param.Id != 0 {
		if err = db.Where("id = ?", param.Id).Count(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("查询ids总数错误: %v", err)
		}
		if err = db.Where("id = ?", param.Id).Find(&user).Error; err != nil {
			return nil, 0, fmt.Errorf("查询ids错误: %v", err)
		}
	} else {
		// 分页查询
		searchReq := &api.SearchReq{
			Condition: db,
			Table:     &user,
			PageInfo:  param.PageInfo,
		}
		// 用户名模糊查询
		if param.String != "" {
			name := "%" + strings.ToUpper(param.String) + "%"
			db = model.DB.Model(&user).Where("UPPER(name) LIKE ?", name)
			searchReq.Condition = db
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
			// 全部返回
		} else {
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
		}
	}
	// 过滤结果
	var result *[]api.UserRes
	result, err = s.GetResults(&user)
	if err != nil {
		return nil, 0, err
	}
	return result, total, err
}

// 删除用户
func (s *UserService) DeleteUser(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.User{}, ids); err != nil {
		return err
	}
	var user []model.User
	// 开启事务
	tx := model.DB.Begin()
	// 返回用户对象
	if err = tx.Find(&user, ids).Error; err != nil {
		return errors.New("查询用户信息失败")
	}
	// 清除和组的关联
	if err = tx.Model(&user).Association("UserGroups").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 用户与用户组关联 失败")
	}
	// 伪删除用户
	if err = tx.Where("id in (?)", ids).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除用户失败")
	}
	tx.Commit()
	return err
}

// 修改用户状态
func (s *UserService) UpdateStatus(param *api.StatusReq) (err error) {
	if !util2.CheckIdExists(&model.User{}, param.ID) {
		return errors.New("用户不存在")
	}
	err = model.DB.Model(&model.User{}).Where("id = ?", param.ID).Update("status", param.Status).Error
	return err
}

// 修改密码
func (s *UserService) UpdatePasswd(passwd *api.PasswordReq) (err error) {
	if !util2.CheckIdExists(&model.User{}, passwd.ID) {
		return errors.New("用户不存在")
	}
	err = model.DB.Model(&model.User{}).Where("id = ?", passwd.ID).Update("password", passwd.Password).Error
	return err
}

// 获取用户个人信息
func (s *UserService) GetSelfInfo(id uint) (userInfo any, err error) {
	var user model.User
	if !util2.CheckIdExists(&model.User{}, id) {
		return nil, errors.New("用户不存在")
	}
	if err = model.DB.Where("id in (?)", id).First(&user).Error; err != nil {
		return user, errors.New("查询用户个人信息失败")
	}
	var result *[]api.UserRes
	result, err = s.GetResults(&user)
	if err != nil {
		return nil, err
	}
	return result, err
}

// 获取用户关联组信息
func (s *UserService) GetAssGroup(id uint) (groupInfo any, err error) {
	var user model.User
	var group []model.UserGroup
	if err != nil {
		return nil, err
	}
	if !util2.CheckIdExists(&model.User{}, id) {
		return nil, errors.New("用户不存在")
	}
	if err = model.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	if err = model.DB.Model(&user).Association("UserGroups").Find(&group); err != nil {
		return user, err
	}
	return group, err
}

// 登录
func (s *UserService) Login(u *model.User) (userInfo *api.AuthLoginRes, err error) {
	var user model.User
	if err = model.DB.Where("username = ?", u.Username).Preload("UserGroups").First(&user).Error; err != nil {
		return userInfo, errors.New("获取用户对象失败")
	}
	if !util.CheckPassword(user.Password, u.Password) {
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
	// 遍历所有组id
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
func (s *UserService) UpdateKeyFileContext(file *multipart.FileHeader, passphrase string, id uint) error {
	fileP, err := file.Open()
	if err != nil {
		return err
	}
	defer fileP.Close()

	fileBytes, err := io.ReadAll(fileP)
	if err != nil {
		return err
	}

	// AES加密并写入prikey
	var data []byte
	data, err = util.EncryptAESCBC(fileBytes, []byte(consts.AesKey), []byte(consts.AesIv))
	if err != nil {
		return fmt.Errorf("用户私钥加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("pri_key", data).Error; err != nil {
		return errors.New("私钥写入数据库失败")
	}

	// AES加密并写入passphrase
	data, err = util.EncryptAESCBC([]byte(passphrase), []byte(consts.AesKey), []byte(consts.AesIv))
	if err != nil {
		return fmt.Errorf("用户passphrase加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("passphrase", data).Error; err != nil {
		return errors.New("通行证密码写入数据库失败")
	}
	return nil
}

// 通过字符串更新私钥内容
func (s *UserService) UpdateKeyContext(key string, passphrase string, id uint) (err error) {
	// AES加密并写入prikey
	var data []byte
	data, err = util.EncryptAESCBC([]byte(key), []byte(consts.AesKey), []byte(consts.AesIv))
	if err != nil {
		return fmt.Errorf("用户私钥加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("pri_key", data).Error; err != nil {
		return errors.New("私钥字符串写入数据库失败")
	}

	// AES加密并写入passphrase
	data, err = util.EncryptAESCBC([]byte(passphrase), []byte(consts.AesKey), []byte(consts.AesIv))
	if err != nil {
		return fmt.Errorf("用户passphrase加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("passphrase", data).Error; err != nil {
		return errors.New("通行证密码写入数据库失败")
	}
	return nil
}

// 返回用户结果
func (s *UserService) GetResults(userInfo any) (*[]api.UserRes, error) {
	var res api.UserRes
	var result []api.UserRes
	var err error
	if users, ok := userInfo.(*[]model.User); ok {
		for _, user := range *users {
			res = api.UserRes{
				ID:         user.ID,
				Name:       user.Name,
				Username:   user.Username,
				Status:     user.Status,
				Email:      user.Email,
				Mobile:     user.Mobile,
				LoginTime:  uint64(user.LoginTime),
				Expiration: user.Expiration,
				IsAdmin:    user.IsAdmin,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if user, ok := userInfo.(*model.User); ok {
		res = api.UserRes{
			ID:         user.ID,
			Name:       user.Name,
			Username:   user.Username,
			Status:     user.Status,
			Email:      user.Email,
			Mobile:     user.Mobile,
			LoginTime:  uint64(user.LoginTime),
			Expiration: user.Expiration,
			IsAdmin:    user.IsAdmin,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换用户结果失败")
}
