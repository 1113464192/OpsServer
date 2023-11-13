package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api"
	"strings"
)

type OpsService struct {
}

var (
	insOps = &OpsService{}
)

func Ops() *OpsService {
	return insOps
}

func (s *OpsService) GetExecParam(param api.GetExecParamReq) (resParam *api.RunSSHCmdAsyncReq, resConfig *api.SftpReq, err error) {
	var task model.TaskTemplate
	var user model.User
	if err = model.DB.First(&task, param.Tid).Error; err != nil {
		return nil, nil, errors.New("根据id查询任务失败")
	}
	if err = model.DB.First(&user, param.Uid).Error; err != nil {
		return nil, nil, errors.New("根据id查询用户失败")
	}
	resParam = new(api.RunSSHCmdAsyncReq)
	resParam.Key = user.PriKey
	resParam.Passphrase = user.KeyPasswd
	// 获取模板/项目对应主机
	var hosts []model.Host
	count := model.DB.Model(&task).Where("id = ?", task.ID).Association("Hosts").Count()
	if count == 0 {
		var project model.Project
		if err = model.DB.Preload("Hosts").Where("id = ?", task.Pid).First(&project).Error; err != nil {
			return nil, nil, errors.New("查询项目关联主机失败")
		}
		hosts = project.Hosts
	} else {
		if err = model.DB.Preload("Hosts").Where("id = ?", task.ID).Find(&task).Error; err != nil {
			return nil, nil, errors.New("查询task关联主机失败")
		}
		hosts = task.Hosts
	}

	// json参数编出到map
	var args map[string][]string
	if task.Args != "" {
		if err = json.Unmarshal([]byte(task.Args), &args); err != nil {
			return nil, nil, fmt.Errorf("参数字段进行json解析失败: %v", err)
		}
	}

	pathCount := len(args["path"])

	// 如果有设置条件 则筛选符合条件的主机
	var memSize float32
	if task.Condition != "" {
		if err = s.filterConditionHost(&hosts, &task, resParam, &memSize); err != nil {
			return nil, nil, fmt.Errorf("筛选符合条件的主机失败: %v", err)
		}
	}

	// 入口
	var typeParam string
	// 判断类型
	switch {
	case strings.Contains(task.TypeName, "装服"):
		typeParam = "装服"
		// case strings.Contains(task.TypeName, "更新"):
		// 	typeParam = "更新"
	}

	// 获取参数
	switch typeParam {
	case "装服":
		if pathCount == 0 {
			return nil, nil, errors.New("path参数数量为0")
		}
		var hostList *[]model.Host
		if task.PortRule != "" {
			hostList, err = s.filterPortRuleHost(&hosts, &task, resParam, &args, memSize)
			if err != nil {
				return nil, nil, fmt.Errorf("端口筛选报错: %v", err)
			}
			hosts = *hostList
		}
		resParam, resConfig, err = s.getInstallServer(&hosts, &task, &user, pathCount, &args)
		if err != nil {
			return nil, nil, fmt.Errorf("获取装服参数报错: %v", err)
		}
		return resParam, resConfig, err
	// case "更新":

	// 不指定path的全机器操作
	default:
		resParam, resConfig, err = s.getGeneral(&hosts, pathCount, &task, &args)
		if err != nil {
			return nil, nil, fmt.Errorf("获取参数报错: %v", err)
		}
		return resParam, resConfig, err
	}
}
