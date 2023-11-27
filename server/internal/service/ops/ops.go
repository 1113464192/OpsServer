package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/api/ops"
	"fqhWeb/pkg/util"
	"fqhWeb/pkg/util2"
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

// 提交工单
func (s *OpsService) SubmitTask(param ops.SubmitTaskReq) (result *[]ops.TaskRecordRes, err error) {
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

	sshReq := &[]api.SSHClientConfigReq{}
	sftpReq := &[]api.SFTPClientConfigReq{}

	// 入口
	var typeParam string
	// 获取操作类型
	switch {
	case strings.Contains(task.TypeName, consts.OperationInstallServerType):
		typeParam = consts.OperationInstallServerType
	}

	// 执行指定操作
	switch typeParam {
	// 单服装服操作
	case consts.OperationInstallServerType:
		if sshReq, sftpReq, err = s.opsInstallServer(pathCount, &task, &hosts, &user, &args, sshReq); err != nil {
			return nil, fmt.Errorf("提交%s工单失败: %v", consts.OperationInstallServerType, err)
		}
	// 服务端更新操作
	// case consts.OperationUpdateServerType:

	// 未知类型
	default:
		return nil, errors.New("没有对应模板类型, 请检查模板类型是否定义正确, 如需添加请联系运维")
	}
	// Auditor
	taskRecord, err = s.writingTaskRecord(sshReq, sftpReq, &user, &task, &args, param.Auditor)
	if err != nil {
		return nil, fmt.Errorf("写入TaskRecord失败: %v", err)
	}
	result, err = s.GetResults(taskRecord)
	if err != nil {
		return nil, fmt.Errorf("转换结果输出失败: %v", err)
	}
	return result, err
}

// 获取工单
func (s *OpsService) GetTask(param *ops.GetTaskReq) (result *[]ops.TaskRecordRes, total int64, err error) {
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

// 获取执行参数
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
	return &sshReq, &sftpReq, err
}

// 对应用户审核工单
func (s *OpsService) ApproveTask(param ops.ApproveTaskReq, uid uint) (res string, err error) {
	var task model.TaskRecord
	var nonApproverSlice []uint
	if nonApproverSlice, err = s.getNonApprover(param.Id); err != nil {
		return "", err
	}
	if !util.IsSliceContain(nonApproverSlice, uid) {
		return "", fmt.Errorf("未审批人:%v中 不包含 %v", nonApproverSlice, uid)
	}
	// 判断审批通过还是拒绝
	if param.Status == 1 {
		res = "审批通过"
		nonApprover := util.DeleteUintSlice(nonApproverSlice, uid)
		if len(nonApprover) != 0 {
			// 更改taskRecord表的nonApprover
			nonApproverMap := make(map[string][]uint)
			nonApproverMap["ids"] = nonApprover
			var data []byte
			data, err = json.Marshal(nonApproverMap)
			if err != nil {
				return "", fmt.Errorf("map转换json失败: %v", err)
			}
			if err = model.DB.Model(&task).Where("id = ?", param.Id).Update("non_approver", string(data)).Error; err != nil {
				return "", fmt.Errorf("更改TaskRecord表的NonApprover失败: %v", err)
			}
			// 向下一个审批者发送审批信息
			fmt.Println("==========向下一个审批者发送审批信息===========")
		} else {
			if err = model.DB.Model(&task).Where("id = ?", param.Id).Updates(model.TaskRecord{Status: 1, NonApprover: `{"ids":[]}`}).Error; err != nil {
				return "", fmt.Errorf("更改工单状态为可执行状态失败: %v", err)
			}
			// 向操作者发送信息
			fmt.Println("==========向操作者发送信息===========")
		}
	} else if param.Status == 4 {
		res = "审批拒绝"
		if err = model.DB.Model(&task).Where("id = ?", param.Id).Update("status", 4).Error; err != nil {
			return "", fmt.Errorf("更改工单状态为可执行状态失败: %v", err)
		}
	} else {
		res = "Status传参错误, 参数不在给定条件中"
		return "", errors.New("status传参错误, 参数不在给定条件中")
	}

	return res, err
}

// 删除工单
func (s *OpsService) DeleteTask(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.TaskRecord{}, ids); err != nil {
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
	if err = tx.Where("id IN ?", ids).Delete(&model.TaskRecord{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除工单失败")
	}
	tx.Commit()
	return err
}

// 获取未审批用户
func (s *OpsService) getNonApprover(id uint) (nonApproverSlice []uint, err error) {
	var task model.TaskRecord
	if err = model.DB.First(&task, id).Error; err != nil {
		return nil, fmt.Errorf("查询TaskRecord数据失败: %v", err)
	}
	nonApproverJson := make(map[string][]uint)
	if err = json.Unmarshal([]byte(task.NonApprover), &nonApproverJson); err != nil {
		return nil, errors.New("sshJson进行json解析失败")
	}
	nonApproverSlice = nonApproverJson["ids"]
	return nonApproverSlice, err
}

// 写入工单操作后的结果入库
func (s *OpsService) execTaskResultWriteDB(status uint8, tid uint, data *map[string][]api.SSHResultRes) (err error) {
	var task model.TaskRecord
	var res []byte
	if res, err = json.Marshal(*data); err != nil {
		return fmt.Errorf("json编码报错: %v", err)
	}
	if err = model.DB.Model(&task).Where("id = ?", tid).Updates(model.TaskRecord{Status: status, Response: string(res)}).Error; err != nil {
		return fmt.Errorf("写入数据到工单结果表中报错: %v", err)
	}
	return err
}

func (s *OpsService) recordServerList(data *string, pid uint) (err error) {
	var args map[string][]string
	if err = json.Unmarshal([]byte(*data), &args); err != nil {
		return fmt.Errorf("json解析失败: %v", err)
	}
	pathCount := len(args["path"])
	if pathCount == len(args["serverName"]) {
		var serverList []model.ServerList
		var server model.ServerList
		var hostIds []uint
		strHostIds := args["hostId"]
		if hostIds, err = util.StringSliceToUintSlice(&strHostIds); err != nil {
			return err
		}
		for i := 0; i < pathCount; i++ {
			server = model.ServerList{
				Flag:       args["flag"][i],
				Path:       args["path"][i],
				ServerName: args["serverName"][i],
				HostId:     hostIds[i],
				ProjectId:  pid,
			}
			serverList = append(serverList, server)
		}
		if err = model.DB.Create(serverList).Error; err != nil {
			return errors.New("单服列表写入失败")
		}
		return err
	} else {
		return fmt.Errorf("%s: path数量与serverName数量不对等", consts.OperationInstallServerType)
	}
}

// 执行工单操作
func (s *OpsService) OpsExecTask(id uint) (map[string][]api.SSHResultRes, error) {
	var result *[]api.SSHResultRes
	var err error
	var task model.TaskRecord
	if err = model.DB.Preload("Template").First(&task, id).Error; err != nil {
		return nil, fmt.Errorf("取TaskRecord值报错: %v", err)
	}
	if task.Status != 1 {
		return nil, errors.New("当前工单状态不可执行")
	}

	data := make(map[string][]api.SSHResultRes)
	sshReq, sftpReq, err := s.GetExecParam(id)
	if err != nil {
		return nil, fmt.Errorf("获取执行参数失败: %v", err)
	}
	if len(*sftpReq) != 0 {
		result, err = service.SSH().RunSFTPAsync(sftpReq)
		data["ssh"] = *result
		if err != nil {
			err2 := s.execTaskResultWriteDB(3, id, &data)
			return nil, fmt.Errorf("SFTP执行失败: %v。%v", err, err2)
		}
	}
	if len(*sshReq) != 0 {
		result, err = service.SSH().RunSSHCmdAsync(sshReq)
		data["sftp"] = *result
		if err != nil {
			err2 := s.execTaskResultWriteDB(3, id, &data)
			return nil, fmt.Errorf("SSH执行失败: %v。%v", err, err2)
		}
	}
	if len(data) == 0 {
		return nil, errors.New("没有获取到任务结果")
	}
	if err = s.execTaskResultWriteDB(2, id, &data); err != nil {
		return nil, err
	}

	// 判断是否属于装服类型
	if strings.Contains(task.Template.TypeName, consts.OperationInstallServerType) {
		if err = s.recordServerList(&task.Args, task.Template.Pid); err != nil {
			return data, fmt.Errorf("操作执行完成, 但写入serverlist报错: %v", err)
		}
	}
	return data, err
}

// 返回工单结果
func (s *OpsService) GetResults(taskInfo any) (*[]ops.TaskRecordRes, error) {
	result := &[]ops.TaskRecordRes{}
	var err error
	if task, ok := taskInfo.(*model.TaskRecord); ok {
		var res ops.TaskRecordRes
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

		res = ops.TaskRecordRes{
			ID:          task.ID,
			TaskName:    task.TaskName,
			TemplateId:  task.TemplateId,
			OperatorId:  task.OperatorId,
			Status:      task.Status,
			Response:    task.Response,
			Args:        task.Args,
			SSHReqs:     res.SSHReqs,
			NonApprover: nonApprover,
			Auditor:     auditorIds,
		}
		*result = append(*result, res)
		return result, err
	} else if tasks, ok := taskInfo.(*[]model.TaskRecord); ok {
		for _, task := range *tasks {
			var res ops.TaskRecordRes
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

			res = ops.TaskRecordRes{
				ID:          task.ID,
				TaskName:    task.TaskName,
				TemplateId:  task.TemplateId,
				OperatorId:  task.OperatorId,
				Status:      task.Status,
				Response:    task.Response,
				Args:        task.Args,
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
