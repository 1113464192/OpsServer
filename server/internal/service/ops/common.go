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

func (s *OpsService) templateRender(task *model.TaskTemplate, args *map[string][]string) (cmd []string, config []string, err error) {
	if task.CmdTem == "" && task.ConfigTem == "" {
		return nil, nil, errors.New("任务的命令和传输文件内容都为空")
	}
	pathCount := len((*args)["path"])
	var cmdTem *template.Template
	if task.CmdTem != "" {
		cmdTem, err = template.New("cmdTem").Parse(task.CmdTem)
		if err != nil {
			return cmd, config, fmt.Errorf("无法解析CMD模板: %v", err)
		}
	}
	var configTem *template.Template
	if pathCount != 0 && task.ConfigTem != "" {
		configTem, err = template.New("configTem").Parse(task.ConfigTem)
		if err != nil {
			return cmd, config, fmt.Errorf("无法解析config模板: %v", err)
		}
	}

	var buf strings.Builder
	var bufString string
	serverInfo := utils.SplitStringMap(*args)
	for i := 0; i < len(serverInfo); i++ {
		if task.CmdTem != "" {
			if err = cmdTem.Execute(&buf, serverInfo[i]); err != nil {
				return cmd, config, fmt.Errorf("无法渲染cmd模板: %v", err)
			}
			bufString = buf.String()
			if strings.Contains(bufString, "no value") {
				return cmd, config, fmt.Errorf("cmd模板有变量没有获取对应解析 %s", bufString)
			}
			cmd = append(cmd, bufString)
			buf.Reset()
		}
		if pathCount != 0 && task.ConfigTem != "" {
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
	}
	return cmd, config, err
}

func (s *OpsService) filterConditionHost(hosts *[]model.Host, task *model.TaskTemplate, resParam *api.RunSSHCmdAsyncReq, memSize *float32) (err error) {
	for _, host := range *hosts {
		resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
		resParam.Username = append(resParam.Username, host.User)
		resParam.SSHPort = append(resParam.SSHPort, host.Port)
	}
	hostInfo, err := service.Host().GetHostCurrData(resParam)
	if err != nil {
		logger.Log().Error("Host", "机器数据采集——数据结构有错误", err)
		return fmt.Errorf("机器数据采集——数据结构有错误: %v", err)
	}
	if err := service.Host().WritieToDatabase(hostInfo); err != nil {
		logger.Log().Error("Host", "机器数据采集——数据写入数据库失败", err)
		return fmt.Errorf("机器数据采集——数据写入数据库失败: %v", err)
	}

	var condition map[string][]string
	if err = json.Unmarshal([]byte(task.Condition), &condition); err != nil {
		return errors.New("筛选机器条件规则进行json解析失败")
	}
	var fields []string
	// 为了使用不定长参数的解包方法，所以要设置为interface{}
	var values []any

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

	if len(fields) > 0 {
		conditions := strings.Join(fields, " AND ")
		// 从关联的主机中查询
		if err = model.DB.Where("ipv4 IN ?", resParam.HostIp).Where(conditions, values...).Order("curr_mem").Find(hosts).Error; err != nil {
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

func (s *OpsService) filterPortRuleHost(hosts *[]model.Host, task *model.TaskTemplate, sshReq *api.RunSSHCmdAsyncReq, args *map[string][]string, memSize float32) (hostList *[]model.Host, err error) {
	tmpHosts := make([]model.Host, len(*hosts))
	copy(tmpHosts, *hosts)

	var portRule map[int]string
	if err = json.Unmarshal([]byte(task.PortRule), &portRule); err != nil {
		return nil, errors.New("端口规则进行json解析失败")
	}
	flags, err := s.getFlag(task.Args, args)
	if err != nil {
		return nil, err
	}
	flagsString := utils.IntSliceToStringSlice(flags)
	(*args)["flag"] = flagsString

	var availHost []model.Host
	for _, flag := range flags {
		// 遍历前面进行条件筛选和内存排序的机器
		for i := range tmpHosts {
			host := &tmpHosts[i]
			if host.CurrMem < memSize {
				continue
			}
			sshReq.HostIp = []string{host.Ipv4.String}
			sshReq.SSHPort = []string{host.Port}
			sshReq.Username = []string{host.User}
			var cmdList []string
			var portList []float64
			portList, err = utils.GenerateExprResult(portRule, flag)
			for _, port := range portList {
				cmdShell := fmt.Sprintf(`
				if [[ -z $(netstat -plan | grep %d) ]];then
					echo "success"
				fi`, int(port))
				cmdList = append(cmdList, cmdShell)
			}
			if err != nil {
				continue
			}
			portString := utils.Float64SliceToStringSlice(portList)
			(*args)["port"] = portString

			var sshResult *[]api.SSHResultRes
			count := len(cmdList)
			num := 0
			for _, cmd := range cmdList {
				sshReq.Cmd = []string{cmd}
				sshResult, err = service.SSH().RunSSHCmdAsync(sshReq)
				if err != nil {
					return nil, fmt.Errorf("端口检测命令执行失败: %v", err)
				}
				if strings.Contains((*sshResult)[0].Response, "success") {
					num += 1
				}
			}
			if count == num {
				availHost = append(availHost, *host)
				host.CurrMem = host.CurrMem - memSize
				break
			} else {
				continue
			}
		}
	}
	return &availHost, err
}

func (s *OpsService) writingTaskRecord(resParam *api.RunSSHCmdAsyncReq, resConfig *api.RunSFTPAsyncReq, user *model.User, task *model.TaskTemplate, auditorIds []uint) (taskRecord *model.TaskRecord, err error) {
	// sshReq编码JSON
	sshJson := make(map[string][]string)
	sshJson["ipv4"] = resParam.HostIp
	sshJson["sshPort"] = resParam.SSHPort
	sshJson["username"] = resParam.Username
	sshJson["cmd"] = resParam.Cmd
	if resConfig.FileContent != nil {
		sshJson["config"] = resConfig.FileContent
		sshJson["configPath"] = resConfig.Path
	}
	var data []byte
	data, err = json.Marshal(sshJson)
	if err != nil {
		return nil, fmt.Errorf("map转换json失败: %v", err)
	}

	taskRecord = &model.TaskRecord{
		TaskName:   task.TaskName,
		TemplateId: task.ID,
		OperatorId: user.ID,
		SSHJson:    string(data),
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
