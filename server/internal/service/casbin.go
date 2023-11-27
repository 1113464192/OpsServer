package service

import (
	"errors"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/util"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	casbinUtil "github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type CasbinService struct {
}

var (
	InsCasbin = CasbinService{}
)

func CasbinServiceApp() *CasbinService {
	return &InsCasbin
}

// @function: UpdateCasbin
// @description: 更新casbin权限
// @param: authorityId string, casbinInfos []request.CasbinInfo
// @return: error
func (s *CasbinService) UpdateCasbin(groupId string, casbinIds []uint) error {
	s.ClearCasbin(0, groupId)
	var apiList []model.Api
	if err := model.DB.Where("id in (?)", casbinIds).Find(&apiList).Error; err != nil {
		return err
	}
	var rules [][]string
	for _, v := range apiList {
		rules = append(rules, []string{groupId, v.Path, v.Method})
	}
	e := s.Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

// @function: UpdateCasbinApi
// @description: API更新随动
// @param: oldPath string, newPath string, oldMethod string, newMethod string
// @return: error
func (s *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := model.DB.Model(&model.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]any{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}

// @function: GetPolicyPathByGroupId
// @description: 获取权限列表
// @param: groupId string
// @return: res []any
func (s *CasbinService) GetPolicyPathByGroupId(groupId string) (res []any, err error) {
	e := s.Casbin()
	list := e.GetFilteredPolicy(0, groupId)
	var pathList []string
	for _, v := range list {
		pathList = append(pathList, v[1])
	}
	var apiList []model.Api
	if err := model.DB.Where("path in (?)", pathList).Find(&apiList).Error; err != nil {
		return nil, err
	}
	for _, v := range apiList {
		res = append(res, v.ID)
	}
	return res, nil
}

// @function: ClearCasbin
// @description: 清除匹配的权限
// @param: v int, p ...string
// @return: bool
func (s *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := s.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success

}

// @function: Casbin
// @description: 持久化到数据库  引入自定义规则
// @return: *casbin.Enforcer
var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func (s *CasbinService) Casbin() *casbin.SyncedEnforcer {
	once.Do(func() {
		a, _ := gormadapter.NewAdapterByDB(model.DB)
		syncedEnforcer, _ = casbin.NewSyncedEnforcer(util.GetRootPath()+"/configs/casbin.conf", a)
		syncedEnforcer.AddFunction("ParamsMatch", s.ParamsMatchFunc)
	})
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}

// @function: ParamsMatch
// @description: 自定义规则函数
// @param: fullNameKey1 string, key2 string
// @return: bool
func (s *CasbinService) ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return casbinUtil.KeyMatch2(key1, key2)
}

// @function: ParamsMatchFunc
// @description: 自定义规则函数
// @param: args ...any
// @return: any, error
func (s *CasbinService) ParamsMatchFunc(args ...any) (any, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return s.ParamsMatch(name1, name2), nil
}
