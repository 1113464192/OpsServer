package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/utils"
	"fqhWeb/pkg/utils2"
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

func (s *OpsService) SubmitTask(param api.SubmitTaskReq) (result *[]api.TaskRecordRes, err error) {
	var taskRecord *model.TaskRecord
	var task model.TaskTemplate
	var user model.User
	if err = model.DB.First(&task, param.Tid).Error; err != nil {
		return nil, errors.New("根据id查询任务失败")
	}
	if err = model.DB.First(&user, param.Uid).Error; err != nil {
		return nil, errors.New("根据id查询用户失败")
	}

	if task.CmdTem == "" && task.ConfigTem == "" {
		return nil, errors.New("任务的命令和传输文件内容都为空")
	}

	var resParam *api.RunSSHCmdAsyncReq
	var resConfig *api.RunSFTPAsyncReq
	resParam = new(api.RunSSHCmdAsyncReq)
	if task.CmdTem != "" {
		resParam.Key = user.PriKey
		resParam.Passphrase = user.KeyPasswd
	}
	resConfig = new(api.RunSFTPAsyncReq)
	if task.ConfigTem != "" {
		resConfig.Key = user.PriKey
		resConfig.Passphrase = user.KeyPasswd
	}

	// 获取模板/项目对应主机
	var hosts []model.Host
	count := model.DB.Model(&task).Where("id = ?", task.ID).Association("Hosts").Count()
	if count == 0 {
		var project model.Project
		if err = model.DB.Preload("Hosts").Where("id = ?", task.Pid).First(&project).Error; err != nil {
			return nil, errors.New("查询项目关联主机失败")
		}
		hosts = project.Hosts
	} else {
		if err = model.DB.Preload("Hosts").Where("id = ?", task.ID).Find(&task).Error; err != nil {
			return nil, errors.New("查询task关联主机失败")
		}
		hosts = task.Hosts
	}

	// json参数编出到map
	var args map[string][]string
	if task.Args != "" {
		if err = json.Unmarshal([]byte(task.Args), &args); err != nil {
			return nil, fmt.Errorf("参数字段进行json解析失败: %v", err)
		}
	}

	pathCount := len(args["path"])

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
			return nil, errors.New("path参数数量为0")
		}
		var hostList *[]model.Host
		// 如果有设置条件 则筛选符合条件的主机
		var memSize float32
		if task.Condition != "" {
			if err = s.filterConditionHost(&hosts, &task, resParam, &memSize); err != nil {
				return nil, fmt.Errorf("筛选符合条件的主机失败: %v", err)
			}
		}
		if task.PortRule != "" {
			hostList, err = s.filterPortRuleHost(&hosts, &task, resParam, &args, memSize)
			if err != nil {
				return nil, fmt.Errorf("端口筛选报错: %v", err)
			}
			hosts = *hostList
		}
		resParam, resConfig, err = s.getInstallServer(&hosts, &task, &user, pathCount, &args, resConfig)
		if err != nil {
			return nil, fmt.Errorf("获取装服参数报错: %v", err)
		}

	// case "更新":

	// 关联机器全操作
	default:
		// 看有无需求, 要做成多cmd轮流执行, 没有的话就单cmd, 装服再多CMD
		resParam, resConfig, err = s.getGeneral(&hosts, &task, &args, resParam, resConfig)
		if err != nil {
			return nil, fmt.Errorf("获取参数报错: %v", err)
		}
	}
	// Auditor
	taskRecord, err = s.writingTaskRecord(resParam, resConfig, &user, &task, param.Auditor)
	if err != nil {
		return nil, fmt.Errorf("写入TaskRecord失败: %v", err)
	}
	result, err = s.GetResults(&taskRecord)
	if err != nil {
		return nil, fmt.Errorf("转换结果输出失败: %v", err)
	}
	return result, err
}

func (s *OpsService) GetTask(param *api.GetTaskReq) (result *[]api.TaskRecordRes, total int64, err error) {
	var task []model.TaskRecord
	db := model.DB.Model(&task)
	// id存在返回id对应model
	if param.Tid != 0 {
		if err = db.Where("id = ?", param.Tid).Count(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id数错误: %v", err)
		}
		if err = db.Where("id = ?", param.Tid).Find(&task).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id错误: %v", err)
		}
	} else {
		searchReq := &api.SearchReq{
			Condition: db,
			Table:     &task,
			PageInfo:  param.PageInfo,
		}
		// 返回name的模糊匹配
		if param.TaskName != "" {
			name := "%" + strings.ToUpper(param.TaskName) + "%"
			db = model.DB.Where("UPPER(task_name) LIKE ?", name).Order("id desc")
			searchReq.Condition = db
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
			// 返回所有
		} else {
			searchReq.Condition = db.Order("id desc")
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
		}
		// 全展示即首页展示，去掉sshJson，否则传输数据太大
		for i := 0; i < len(task); i++ {
			task[i].SSHJson = ""
		}
	}
	result, err = s.GetResults(&task)
	if err != nil {
		return nil, total, fmt.Errorf("转换结果输出失败: %v", err)
	}
	return result, total, err

}

func (s *OpsService) GetExecParam(tid uint) (resParam *api.RunSSHCmdAsyncReq, resConfig *api.RunSFTPAsyncReq, err error) {
	var task model.TaskRecord
	if err = model.DB.First(&task, tid).Error; err != nil {
		return nil, nil, fmt.Errorf("查询TaskRecord数据失败: %v", err)
	}

	var user model.User
	if err = model.DB.First(&user, task.OperatorId).Error; err != nil {
		return nil, nil, fmt.Errorf("查询用户数据失败: %v", err)
	}

	sshJson := make(map[string][]string)
	if err = json.Unmarshal([]byte(task.SSHJson), &sshJson); err != nil {
		return nil, nil, errors.New("sshJson进行json解析失败")
	}
	hostIp := sshJson["ipv4"]
	username := sshJson["username"]
	sshPort := sshJson["sshPort"]
	cmd := sshJson["cmd"]
	config := sshJson["config"]

	resParam = &api.RunSSHCmdAsyncReq{
		HostIp:     hostIp,
		Username:   username,
		SSHPort:    sshPort,
		Cmd:        cmd,
		Key:        user.PriKey,
		Passphrase: user.KeyPasswd,
	}
	if config != nil {
		path := sshJson["configPath"]
		resConfig = &api.RunSFTPAsyncReq{
			HostIp:      hostIp,
			Username:    username,
			SSHPort:     sshPort,
			FileContent: config,
			Path:        path,
			Key:         user.PriKey,
			Passphrase:  user.KeyPasswd,
		}
		return resParam, resConfig, err
	}
	return resParam, nil, err
}

func (s *OpsService) ApproveTask(tid uint, uid uint) error {
	var err error
	var task model.TaskRecord
	if err = model.DB.First(&task, tid).Error; err != nil {
		return fmt.Errorf("查询TaskRecord数据失败: %v", err)
	}
	nonApproverJson := make(map[string][]uint)
	if err = json.Unmarshal([]byte(task.NonApprover), &nonApproverJson); err != nil {
		return errors.New("sshJson进行json解析失败")
	}
	nonApproverSlice := nonApproverJson["ids"]
	if !utils.IsSliceContain(nonApproverSlice, uid) {
		return fmt.Errorf("未审批人:%v中 不包含 %v", nonApproverSlice, uid)
	}
	nonApprover := utils.DeleteUintSlice(nonApproverSlice, uid)
	if len(nonApprover) != 0 {
		// 向下一个审批者发送审批信息
		fmt.Println("==========向下一个审批者发送审批信息===========")
	} else {
		if err = model.DB.Model(&task).Where("id = ?", task.ID).Update("status", 1).Error; err != nil {
			return fmt.Errorf("更改工单状态为可执行失败")
		}
		// 向操作者发送信息
		fmt.Println("==========向操作者发送信息===========")
	}
	return err
}

func (s *OpsService) DeleteTask(id uint) (err error) {
	if !utils2.CheckIdExists(&model.TaskRecord{}, &id) {
		return errors.New("工单不存在")
	}
	var task model.TaskRecord
	tx := model.DB.Begin()
	if err = tx.First(&task, id).Error; err != nil {
		return errors.New("查询工单信息失败")
	}
	if err = tx.Model(&task).Association("Auditor").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 工单与用户关联 失败")
	}
	if err = tx.Where("id = ?", id).Delete(&model.TaskRecord{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除工单失败")
	}
	tx.Commit()
	return err
}

func (s *OpsService) OpsExecTask(id uint) (map[string][]api.SSHResultRes, error) {
	var result *[]api.SSHResultRes
	var err error
	data := make(map[string][]api.SSHResultRes)
	resParam, resConfig, err := s.GetExecParam(id)
	if err != nil {
		return nil, fmt.Errorf("获取执行参数失败: %v", err)
	}
	if resConfig.FileContent != nil {
		result, err = service.SSH().RunSFTPAsync(resConfig)
		if err != nil {
			return nil, fmt.Errorf("测试执行失败: %v", err)
		}
		data["ssh"] = *result
	}
	if resParam.Cmd != nil {
		result, err = service.SSH().RunSFTPAsync(resConfig)
		if err != nil {
			return nil, fmt.Errorf("测试执行失败: %v", err)
		}
		data["sftp"] = *result
	}
	if len(data) == 0 {
		return nil, errors.New("没有获取到任务结果")
	}
	return data, err
}

// 返回工单结果
func (s *OpsService) GetResults(taskInfo any) (*[]api.TaskRecordRes, error) {
	var res api.TaskRecordRes
	var result []api.TaskRecordRes
	var err error
	if tasks, ok := taskInfo.(*[]model.TaskRecord); ok {
		for _, task := range *tasks {
			sshJson := make(map[string][]string)
			if err = json.Unmarshal([]byte(task.SSHJson), &sshJson); err != nil {
				return nil, errors.New("sshJson进行json解析失败")
			}
			hostIp := sshJson["ipv4"]
			username := sshJson["username"]
			sshPort := sshJson["sshPort"]
			cmd := sshJson["cmd"]
			config := sshJson["config"]
			var path []string
			if config != nil {
				path = sshJson["configPath"]
			}

			nonApproverJson := make(map[string][]uint)
			if err = json.Unmarshal([]byte(task.NonApprover), &nonApproverJson); err != nil {
				return nil, errors.New("sshJson进行json解析失败")
			}
			nonApprover := nonApproverJson["ids"]
			var auditor []model.User
			if err = model.DB.Model(&task).Association("Auditor").Find(&auditor); err != nil {
				return nil, fmt.Errorf("获取关联用户失败: %v", err)
			}
			var auditorIds []uint
			for _, a := range auditor {
				auditorIds = append(auditorIds, a.ID)
			}

			res = api.TaskRecordRes{
				ID:          task.ID,
				TaskName:    task.TaskName,
				TemplateId:  task.TemplateId,
				OperatorId:  task.OperatorId,
				Status:      task.Status,
				Response:    task.Response,
				HostIp:      hostIp,
				Username:    username,
				SSHPort:     sshPort,
				Cmd:         cmd,
				ConfigPath:  path,
				FileContent: config,
				NonApprover: nonApprover,
				Auditor:     auditorIds,
			}
			result = append(result, res)
		}
		return &result, err
	} else if task, ok := taskInfo.(*model.TaskRecord); ok {
		sshJson := make(map[string][]string)
		if err = json.Unmarshal([]byte(task.SSHJson), &sshJson); err != nil {
			return nil, errors.New("sshJson进行json解析失败")
		}
		hostIp := sshJson["ipv4"]
		username := sshJson["username"]
		sshPort := sshJson["sshPort"]
		cmd := sshJson["cmd"]
		config := sshJson["config"]
		var path []string
		if config != nil {
			path = sshJson["configPath"]
		}

		nonApproverJson := make(map[string][]uint)
		if err = json.Unmarshal([]byte(task.NonApprover), &nonApproverJson); err != nil {
			return nil, errors.New("sshJson进行json解析失败")
		}
		nonApprover := nonApproverJson["ids"]
		var auditor []model.User
		if err = model.DB.Model(&task).Association("Auditor").Find(&auditor); err != nil {
			return nil, fmt.Errorf("获取关联用户失败: %v", err)
		}
		var auditorIds []uint
		for _, a := range auditor {
			auditorIds = append(auditorIds, a.ID)
		}

		res = api.TaskRecordRes{
			ID:          task.ID,
			TaskName:    task.TaskName,
			TemplateId:  task.TemplateId,
			OperatorId:  task.OperatorId,
			Status:      task.Status,
			Response:    task.Response,
			HostIp:      hostIp,
			Username:    username,
			SSHPort:     sshPort,
			Cmd:         cmd,
			ConfigPath:  path,
			FileContent: config,
			NonApprover: nonApprover,
			Auditor:     auditorIds,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换taskRecord结果失败")
}
