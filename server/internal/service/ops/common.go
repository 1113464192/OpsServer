package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

func (s *OpsService) getFlag(param string, args *map[string][]string) (flags []int, err error) {
	re := regexp.MustCompile(`\d+`)
	for _, path := range (*args)["path"] {
		matches := re.FindAllString(path, -1)
		if len(matches) != 1 {
			return nil, fmt.Errorf("从path取出的int不止一个: %v", matches)
		}
		match := matches[0]
		flag, err := strconv.Atoi(match)
		if err != nil {
			return nil, fmt.Errorf(path+" 字符串转换整数失败: %v", err)
		}
		flags = append(flags, flag)
	}
	return flags, err
}

func (s *OpsService) templateRender(task *model.TaskTemplate, args *map[string][]string, pathCount int) (cmd []string, config []string, err error) {
	var cmdTem *template.Template
	cmdTem, err = template.New("cmdTem").Parse(task.CmdTem)
	if err != nil {
		return cmd, config, fmt.Errorf("无法解析CMD模板: %v", err)
	}
	var configTem *template.Template
	configTem, err = template.New("configTem").Parse(task.ConfigTem)
	if err != nil {
		return cmd, config, fmt.Errorf("无法解析config模板: %v", err)
	}

	var buf strings.Builder
	var bufString string
	// 对map[string][]string进行拆解，以便模板渲染
	serverInfo := utils.SplitStringMap(*args)
	for i := 0; i < len(serverInfo); i++ {
		if err = cmdTem.Execute(&buf, serverInfo[i]); err != nil {
			return cmd, config, fmt.Errorf("无法渲染cmd模板: %v", err)
		}
		bufString = buf.String()
		if strings.Contains(bufString, "no value") {
			return cmd, config, fmt.Errorf("cmd模板有变量没有获取对应解析 %s", bufString)
		}
		cmd = append(cmd, bufString)
		buf.Reset()
		if err = configTem.Execute(&buf, serverInfo[i]); err != nil {
			return cmd, config, fmt.Errorf("无法渲染config模板: %v", err)
		}
		bufString = buf.String()
		if strings.Contains(bufString, "no value") {
			return cmd, config, fmt.Errorf("config模板有变量没有获取对应解析 %s", bufString)
		}
		config = append(config, bufString)
		buf.Reset()
	}
	// 兼容所有单服一套命令或配置
	cmdCount := len(cmd)
	if cmdCount != 1 && cmdCount != pathCount {
		return nil, nil, errors.New("CMD不等于1, 也不等于路径数量")
	}
	configCount := len(config)
	if configCount != 1 && configCount != pathCount {
		return nil, nil, errors.New("CONFIG不等于1, 也不等于路径数量")
	}
	return cmd, config, err
}

// 按条件筛选符合条件的服务器
func (s *OpsService) filterConditionHost(hosts *[]model.Host, user *model.User, task *model.TaskTemplate, sshReq *[]api.SSHClientConfigReq, memSize *float32) (err error) {
	// 赋值可用IP给ssh命令参数
	for i := 0; i < len(*sshReq); i++ {
		(*sshReq)[i].HostIp = (*hosts)[i].Ipv4.String
		(*sshReq)[i].Username = (*hosts)[i].User
		(*sshReq)[i].SSHPort = (*hosts)[i].Port
		(*sshReq)[i].Key = user.PriKey
		(*sshReq)[i].Passphrase = user.Passphrase
	}
	// 更新每个服务器的最新状态
	hostInfo, err := service.Host().GetHostCurrData(sshReq)
	if err != nil {
		logger.Log().Error("Host", "机器数据采集——数据结构有错误", err)
		return fmt.Errorf("机器数据采集——数据结构有错误: %v", err)
	}
	if err := service.Host().WritieToDatabase(hostInfo); err != nil {
		logger.Log().Error("Host", "机器数据采集——数据写入数据库失败", err)
		return fmt.Errorf("机器数据采集——数据写入数据库失败: %v", err)
	}

	// 获取允许执行命令的服务器资源条件
	var condition map[string][]string
	if err = json.Unmarshal([]byte(task.Condition), &condition); err != nil {
		return errors.New("筛选机器条件规则进行json解析失败")
	}
	var fields []string
	// 为了使用不定长参数的解包方法，所以要设置为[]any
	var values []any

	// 将条件组合成切片
	for key, value := range condition {
		switch key {
		case "mem":
			fields = append(fields, "curr_mem > ?")
			values = append(values, value[0])
			memFloat, err := strconv.ParseFloat(value[0], 64)
			if err != nil {
				return fmt.Errorf(value[0]+"转换为浮点数失败: %v", err)
			}
			*memSize = float32(memFloat) * float32(1024)

		case "data_disk":
			fields = append(fields, "curr_data_disk > ?")
			values = append(values, value[0])
		case "iowait":
			fields = append(fields, "curr_iowait < ?")
			values = append(values, value[0])
		case "idle":
			fields = append(fields, "curr_idle > ?")
			values = append(values, value[0])
		case "load":
			fields = append(fields, "curr_load < ?")
			values = append(values, value[0])
		default:
			return fmt.Errorf("%s 不属于ConditionSet中的任何一个", key)
		}
	}

	// 将可用IP做成切片，方便后续执行GORM命令
	var hostsIp []string
	for _, sshReq := range *sshReq {
		hostsIp = append(hostsIp, sshReq.HostIp)
	}
	if len(fields) > 0 {
		conditions := strings.Join(fields, " AND ")
		// 从关联的主机中获取符合条件的单服
		if err = model.DB.Where("ipv4 IN ?", hostsIp).Where(conditions, values...).Order("curr_mem").Find(hosts).Error; err != nil {
			return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
		}
	}
	return err
	// for _, c := range condition {
	// 	for key, value := range c {
	// 		switch key {
	// 		case "mem":
	// 			if err = model.DB.Where("curr_mem > ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "data_disk":
	// 			if err = model.DB.Where("curr_data_disk > ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "iowait":
	// 			if err = model.DB.Where("curr_iowait < ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "idle":
	// 			if err = model.DB.Where("curr_idle > ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "load":
	// 			if err = model.DB.Where("curr_load < ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		default:
	// 			return fmt.Errorf("%s 不属于ConditionSet中的其中一个", key)
	// 		}
	// 	}
	// }
	// return err
}

// 按端口规则筛选可用服务器
func (s *OpsService) filterPortRuleHost(hosts *[]model.Host, user *model.User, task *model.TaskTemplate, sshReq *[]api.SSHClientConfigReq, args *map[string][]string, memSize float32) (hostList *[]model.Host, err error) {
	// 做一个用来计算的host切片，避免影响到host表(防止后面突然有人Save)
	tmpHosts := make([]model.Host, len(*hosts))
	copy(tmpHosts, *hosts)

	// 解析出每个端口规则
	var portRule map[string]string
	if err = json.Unmarshal([]byte(task.PortRule), &portRule); err != nil {
		return nil, errors.New("端口规则进行json解析失败")
	}
	// 获取每个服的字符串标识
	flags, err := s.getFlag(task.Args, args)
	if err != nil {
		return nil, err
	}
	// 添加到后续的模板变量映射
	flagsString := utils.IntSliceToStringSlice(flags)
	(*args)["flag"] = flagsString

	// 接收可用host切片
	var availHost []model.Host
	// 遍历每个服
	for _, flag := range flags {
		// 遍历前面进行条件筛选和内存排序的机器
		for i := range tmpHosts {
			host := &tmpHosts[i]
			if host.CurrMem < memSize {
				continue
			}
			// 获取基于flag的端口
			var portList []float64
			portList, err = utils.GenerateExprResult(portRule, flag)
			if len(portList) != len(portRule) {
				return nil, errors.New("取出端口数量不等于端口规则数量")
			}
			// 添加到后续的模板变量映射
			var p int
			for key := range portRule {
				portString := strconv.FormatFloat(portList[p], 'f', -1, 64)
				(*args)[key] = append((*args)[key], portString)
				p += 1
			}
			if err != nil {
				return nil, fmt.Errorf("端口规则\n%v\n基于flag %d 生成端口失败: %v", portRule, flag, err)
			}
			// 计算有多少个端口需要检查占用
			var c int
			for _, port := range portList {
				cmdShell := fmt.Sprintf(`
				if [[ -z $(netstat -plan | grep %d) ]];then
					echo "success"
				fi`, int(port))
				// 兼容多个端口多个命令，这样子就允许单主机多命令
				n := len(*sshReq)
				(*sshReq)[n].HostIp = host.Ipv4.String
				(*sshReq)[n].SSHPort = host.Port
				(*sshReq)[n].Username = host.User
				(*sshReq)[n].Key = user.PriKey
				(*sshReq)[n].Passphrase = user.Passphrase
				(*sshReq)[n].Cmd = cmdShell
				c += 1
			}
			var sshResult *[]api.SSHResultRes
			// 判断成功次数
			var successNum int
			sshResult, err = service.SSH().RunSSHCmdAsync(sshReq)
			if err != nil {
				return nil, fmt.Errorf("端口检测命令执行失败: %v", err)
			}
			// 判断每个端口占用
			for _, res := range *sshResult {
				if strings.Contains(res.Response, "success") {
					successNum += 1
				}
			}
			// 查看当前循环的服务器能否装服
			if successNum == c {
				availHost = append(availHost, *host)
				host.CurrMem = host.CurrMem - memSize
				// 开始循环下一个预备单服
				break
			}
			continue
		}
	}
	return &availHost, err
}

func (s *OpsService) writingTaskRecord(sshReq *[]api.SSHClientConfigReq, sftpReq *[]api.SFTPClientConfigReq, user *model.User, task *model.TaskTemplate, auditorIds []uint) (taskRecord *model.TaskRecord, err error) {
	// sshReq编码JSON
	var data []byte
	data, err = json.Marshal(*sshReq)
	if err != nil {
		return nil, fmt.Errorf("map转换json失败: %v", err)
	}

	taskRecord = &model.TaskRecord{
		TaskName:   task.TaskName,
		TemplateId: task.ID,
		OperatorId: user.ID,
		SSHJson:    string(data),
	}
	if len(*sftpReq) != 0 {
		data, err = json.Marshal(*sftpReq)
		if err != nil {
			return nil, fmt.Errorf("map转换json失败: %v", err)
		}
		taskRecord.SFTPJson = string(data)
	}
	tx := model.DB.Begin()
	if err = tx.Create(taskRecord).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("TaskRecord存储失败: %v", err)
	}
	var auditor model.User
	if err = tx.Model(&model.User{}).Find(&auditor, auditorIds).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("查询工单审批用户失败: %v", err)
	}
	if err = tx.Model(taskRecord).Association("Auditor").Replace(&auditor); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("工单任务 关联 审批用户 失败: %v", err)
	}
	tx.Commit()

	nonApprover := make(map[string][]uint)
	nonApprover["ids"] = auditorIds

	data, err = json.Marshal(nonApprover)
	if err != nil {
		return nil, fmt.Errorf("map转换json失败: %v", err)
	}
	taskRecord.NonApprover = string(data)
	// 接入微信小程序之类的请求,向第一个审批用户发送
	// ......
	fmt.Println("==========首次写入,接入微信小程序之类的请求,向第一个审批用户发送===========")
	model.DB.Save(&taskRecord)

	return taskRecord, err
}
