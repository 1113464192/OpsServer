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

	sshReq := &[]api.SSHClientConfigReq{}
	sftpReq := &[]api.SFTPClientConfigReq{}

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
		if task.CmdTem == "" || task.ConfigTem == "" {
			return nil, errors.New("任务的命令和传输文件内容都为空")
		}
		var hostList *[]model.Host
		// 如果有设置条件 则筛选符合条件的主机
		var memSize float32
		if task.Condition != "" {
			if err = s.filterConditionHost(&hosts, &user, &task, sshReq, &memSize); err != nil {
				return nil, fmt.Errorf("筛选符合条件的主机失败: %v", err)
			}
		}
		// 筛选端口规则
		if task.PortRule != "" {
			hostList, err = s.filterPortRuleHost(&hosts, &user, &task, sshReq, &args, memSize)
			if err != nil {
				return nil, fmt.Errorf("端口筛选报错: %v", err)
			}
			hosts = *hostList
		}
		sshReq, sftpReq, err = s.getInstallServer(&hosts, &task, &user, pathCount, &args)
		if err != nil {
			return nil, fmt.Errorf("获取装服参数报错: %v", err)
		}

	// case "更新":

	// 关联机器全操作
	default:
		return nil, errors.New("没有对应模板类型, 请检查模板类型是否定义正确, 如需添加请联系运维")
	}
	// Auditor
	taskRecord, err = s.writingTaskRecord(sshReq, sftpReq, &user, &task, param.Auditor)
	if err != nil {
		return nil, fmt.Errorf("写入TaskRecord失败: %v", err)
	}
	result, err = s.GetResults(taskRecord)
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
		if err = db.Where("id = ?", param.Tid).First(&task).Error; err != nil {
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
			searchReq.Condition = db.Where("UPPER(task_name) LIKE ?", name).Order("id desc")
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
		// 非ID展示则去掉执行Json，否则传输数据太大
		for i := 0; i < len(task); i++ {
			task[i].SSHJson = ""
			task[i].SFTPJson = ""
		}
	}
	result, err = s.GetResults(&task)
	if err != nil {
		return nil, total, fmt.Errorf("转换结果输出失败: %v", err)
	}
	return result, total, err

}

func (s *OpsService) GetExecParam(tid uint) (*[]api.SSHClientConfigReq, *[]api.SFTPClientConfigReq, error) {
	var sshReq []api.SSHClientConfigReq
	var sftpReq []api.SFTPClientConfigReq
	var err error
	var task model.TaskRecord
	if err = model.DB.First(&task, tid).Error; err != nil {
		return nil, nil, fmt.Errorf("查询TaskRecord数据失败: %v", err)
	}

	var user model.User
	if err = model.DB.First(&user, task.OperatorId).Error; err != nil {
		return nil, nil, fmt.Errorf("查询用户数据失败: %v", err)
	}

	if err = json.Unmarshal([]byte(task.SSHJson), &sshReq); err != nil {
		return nil, nil, errors.New("sshJson进行json解析失败")
	}
	if err = json.Unmarshal([]byte(task.SFTPJson), &sftpReq); err != nil {
		return nil, nil, errors.New("sshJson进行json解析失败")
	}
	// hostIp := sshJson["ipv4"]
	// username := sshJson["username"]
	// sshPort := sshJson["sshPort"]
	// cmd := sshJson["cmd"]
	// config := sshJson["config"]

	// sshReq = &api.RunSSHCmdAsyncReq{
	// 	HostIp:     hostIp,
	// 	Username:   username,
	// 	SSHPort:    sshPort,
	// 	Cmd:        cmd,
	// 	Key:        user.PriKey,
	// 	Passphrase: user.Passphrase,
	// }
	// if config != nil {
	// 	path := sshJson["configPath"]
	// 	sftpReq = &api.RunSFTPAsyncReq{
	// 		HostIp:      hostIp,
	// 		Username:    username,
	// 		SSHPort:     sshPort,
	// 		FileContent: config,
	// 		Path:        path,
	// 		Key:         user.PriKey,
	// 		Passphrase:  user.Passphrase,
	// 	}
	// 	return sshReq, sftpReq, err
	// }
	return &sshReq, &sftpReq, err
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
		// 更改taskRecord表的nonApprover
		nonApproverMap := make(map[string][]uint)
		nonApproverMap["ids"] = nonApprover
		var data []byte
		data, err = json.Marshal(nonApproverMap)
		if err != nil {
			return fmt.Errorf("map转换json失败: %v", err)
		}
		if err = model.DB.Model(&task).Where("id = ?", task.ID).Update("non_approver", string(data)).Error; err != nil {
			return fmt.Errorf("更改TaskRecord表的NonApprover失败: %v", err)
		}
		// 向下一个审批者发送审批信息
		fmt.Println("==========向下一个审批者发送审批信息===========")
	} else {
		if err = model.DB.Model(&task).Where("id = ?", task.ID).Updates(model.TaskRecord{Status: 1, NonApprover: `{"ids":[]}`}).Error; err != nil {
			return fmt.Errorf("更改工单状态为可执行状态失败")
		}
		// 向操作者发送信息
		fmt.Println("==========向操作者发送信息===========")
	}
	return err
}

func (s *OpsService) DeleteTask(ids []uint) (err error) {
	if err = utils2.CheckIdsExists(model.TaskRecord{}, ids); err != nil {
		return err
	}
	var task []model.TaskRecord
	tx := model.DB.Begin()
	if err = tx.Find(&task, ids).Error; err != nil {
		return errors.New("查询工单信息失败")
	}
	if err = tx.Model(&task).Association("Auditor").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 工单与用户关联 失败")
	}
	if err = tx.Where("id IN ?", ids).Delete(&[]model.TaskRecord{}).Error; err != nil {
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
	sshReq, sftpReq, err := s.GetExecParam(id)
	if err != nil {
		return nil, fmt.Errorf("获取执行参数失败: %v", err)
	}
	if len(*sftpReq) != 0 {
		result, err = service.SSH().RunSFTPAsync(sftpReq)
		if err != nil {
			return nil, fmt.Errorf("SFTP执行失败: %v", err)
		}
		data["ssh"] = *result
	}
	if len(*sshReq) != 0 {
		result, err = service.SSH().RunSSHCmdAsync(sshReq)
		if err != nil {
			return nil, fmt.Errorf("SSH执行失败: %v", err)
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
	result := &[]api.TaskRecordRes{}
	var err error
	// var res api.TaskRecordRes
	// var result []api.TaskRecordRes
	// var err error
	// if tasks, ok := taskInfo.(*[]model.TaskRecord); ok {
	// 	for _, task := range *tasks {
	// 		sshJson := make(map[string][]string)
	// 		if err = json.Unmarshal([]byte(task.SSHJson), &sshJson); err != nil {
	// 			return nil, errors.New("sshJson进行json解析失败")
	// 		}
	// 		hostIp := sshJson["ipv4"]
	// 		username := sshJson["username"]
	// 		sshPort := sshJson["sshPort"]
	// 		cmd := sshJson["cmd"]
	// 		config := sshJson["config"]
	// 		var path []string
	// 		if config != nil {
	// 			path = sshJson["configPath"]
	// 		}

	// 		nonApproverJson := make(map[string][]uint)
	// 		if err = json.Unmarshal([]byte(task.NonApprover), &nonApproverJson); err != nil {
	// 			return nil, errors.New("sshJson进行json解析失败")
	// 		}
	// 		nonApprover := nonApproverJson["ids"]
	// 		var auditor []model.User
	// 		if err = model.DB.Model(&task).Association("Auditor").Find(&auditor); err != nil {
	// 			return nil, fmt.Errorf("获取关联用户失败: %v", err)
	// 		}
	// 		var auditorIds []uint
	// 		for _, a := range auditor {
	// 			auditorIds = append(auditorIds, a.ID)
	// 		}

	// 		res = api.TaskRecordRes{
	// 			ID:          task.ID,
	// 			TaskName:    task.TaskName,
	// 			TemplateId:  task.TemplateId,
	// 			OperatorId:  task.OperatorId,
	// 			Status:      task.Status,
	// 			Response:    task.Response,
	// 			HostIp:      hostIp,
	// 			Username:    username,
	// 			SSHPort:     sshPort,
	// 			Cmd:         cmd,
	// 			ConfigPath:  path,
	// 			FileContent: config,
	// 			NonApprover: nonApprover,
	// 			Auditor:     auditorIds,
	// 		}
	// 		result = append(result, res)
	// 	}
	// 	return &result, err
	// }
	if task, ok := taskInfo.(*model.TaskRecord); ok {
		var res api.TaskRecordRes
		sshCmdReq := []api.SSHClientConfigReq{}
		if task.SFTPJson != "" {
			// 除了CMD全部写入
			if err = json.Unmarshal([]byte(task.SFTPJson), &res.SSHReqs); err != nil {
				return nil, fmt.Errorf("sftpJson进行json解析失败: %v", err)
			}
			// 基于SFTP写入CMD
			if task.SSHJson != "" {
				if err = json.Unmarshal([]byte(task.SSHJson), &sshCmdReq); err != nil {
					return nil, fmt.Errorf("sshJson进行json解析失败: %v", err)
				}
			}
		} else if task.SSHJson != "" {
			// 只写入SSH
			if err = json.Unmarshal([]byte(task.SSHJson), &res.SSHReqs); err != nil {
				return nil, fmt.Errorf("sshJson进行json解析失败: %v", err)
			}
		}

		// 写入cmd给res.SSHReqs
		if task.SFTPJson != "" && task.SSHJson != "" {
			// 判断是否对等
			if len(sshCmdReq) == len(res.SSHReqs) {
				for i := 0; i < len(res.SSHReqs); i++ {
					res.SSHReqs[i].Cmd = sshCmdReq[i].Cmd
				}
			} else {
				return nil, errors.New("ssh切片和sftp切片数量不对等")
			}
		}

		nonApproverJson := make(map[string][]uint)
		if err = json.Unmarshal([]byte(task.NonApprover), &nonApproverJson); err != nil {
			return nil, fmt.Errorf("NonApprover进行json解析失败: %v", err)
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
			SSHReqs:     res.SSHReqs,
			NonApprover: nonApprover,
			Auditor:     auditorIds,
		}
		*result = append(*result, res)
		return result, err
	} else if tasks, ok := taskInfo.(*[]model.TaskRecord); ok {
		for _, task := range *tasks {
			var res api.TaskRecordRes
			sshCmdReq := []api.SSHClientConfigReq{}
			if task.SFTPJson != "" {
				// 除了CMD全部写入
				if err = json.Unmarshal([]byte(task.SFTPJson), &res.SSHReqs); err != nil {
					return nil, fmt.Errorf("sftpJson进行json解析失败: %v", err)
				}
				// 基于SFTP写入CMD
				if task.SSHJson != "" {
					if err = json.Unmarshal([]byte(task.SSHJson), &sshCmdReq); err != nil {
						return nil, fmt.Errorf("sshJson进行json解析失败: %v", err)
					}
				}
			} else if task.SSHJson != "" {
				// 只写入SSH
				if err = json.Unmarshal([]byte(task.SSHJson), &res.SSHReqs); err != nil {
					return nil, fmt.Errorf("sshJson进行json解析失败: %v", err)
				}
			}

			// 写入cmd给res.SSHReqs
			if task.SFTPJson != "" && task.SSHJson != "" {
				// 判断是否对等
				if len(sshCmdReq) == len(res.SSHReqs) {
					for i := 0; i < len(res.SSHReqs); i++ {
						res.SSHReqs[i].Cmd = sshCmdReq[i].Cmd
					}
				} else {
					return nil, errors.New("ssh切片和sftp切片数量不对等")
				}
			}

			nonApproverJson := make(map[string][]uint)
			if err = json.Unmarshal([]byte(task.NonApprover), &nonApproverJson); err != nil {
				return nil, fmt.Errorf("NonApprover进行json解析失败: %v", err)
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
				SSHReqs:     res.SSHReqs,
				NonApprover: nonApprover,
				Auditor:     auditorIds,
			}
			*result = append(*result, res)

		}
		return result, err
	} else {
		return result, errors.New("返回结果转换类型不配对, 转换taskRecord结果失败")
	}
}
