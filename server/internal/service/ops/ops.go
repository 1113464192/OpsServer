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

type OpsService struct {
}

var (
	insOps = &OpsService{}
)

func Ops() *OpsService {
	return insOps
}

func (s *OpsService) getFlag(param string, serverSum *int, args *map[string][]string) (flags []int, err error) {
	re := regexp.MustCompile(`\d+`)
	if len((*args)["path"]) != 0 {
		*serverSum = len((*args)["path"])
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
	}

	return flags, err
}

func (s *OpsService) templateRender(task *model.TaskTemplate, args *map[string][]string) (cmd []string, config []string, err error) {
	pathCount := len((*args)["path"])
	cmdTem, err := template.New("cmdTem").Parse(task.Task)
	if err != nil {
		return cmd, config, fmt.Errorf("无法解析CMD模板: %v", err)
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
		if err = cmdTem.Execute(&buf, serverInfo[i]); err != nil {
			return cmd, config, fmt.Errorf("无法渲染cmd模板: %v", err)
		}
		bufString = buf.String()
		if strings.Contains(bufString, "no value") {
			return cmd, config, fmt.Errorf("cmd模板有变量没有获取对应解析 %s", bufString)
		}
		cmd = append(cmd, bufString)
		buf.Reset()
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

func (s *OpsService) filterPortRuleHost(hosts *[]model.Host, task *model.TaskTemplate, sshReq *api.RunSSHCmdAsyncReq, serverSum *int, args *map[string][]string, memSize float32) (hostList []model.Host, err error) {
	tmpHosts := make([]model.Host, len(*hosts))
	copy(tmpHosts, *hosts)

	var portRule map[int]string
	if err = json.Unmarshal([]byte(task.PortRule), &portRule); err != nil {
		return nil, errors.New("端口规则进行json解析失败")
	}
	flags, err := s.getFlag(task.Args, serverSum, args)
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
	return availHost, err
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
			*memSize = float32(memFloat)

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

	// 使用单个查询筛选符合条件的主机
	// 由于没找到基于hosts进行二次过滤的方法，只能重复操作一次，后续找到办法再进行优化
	if len(fields) > 0 {
		conditions := strings.Join(fields, " AND ")
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

	// resConfig存在的话赋值
	if pathCount != 0 && task.ConfigTem != "" {
		resConfig = new(api.SftpReq)
		resConfig.Key = user.PriKey
		resConfig.Passphrase = user.KeyPasswd
	}

	// 如果有设置条件 则筛选符合条件的主机
	var memSize float32
	if task.Condition != "" {
		if err = s.filterConditionHost(&hosts, &task, resParam, &memSize); err != nil {
			return nil, nil, fmt.Errorf("筛选符合条件的主机失败: %v", err)
		}
	}
	// 如果有设置端口规则 则从符合条件的主机中 筛选符合端口规则的服务器
	if task.PortRule != "" {
		if len(args["path"]) == 0 {
			return nil, nil, errors.New("有端口规则请传path, 否则无标识判断")
		}
		var hostList []model.Host
		var serverSum int
		// 端口检测较为耗时, 需要逐台检测并逐台计算内存损耗，日后做成redis版本可用goroutine解决
		if hostList, err = s.filterPortRuleHost(&hosts, &task, resParam, &serverSum, &args, memSize); err != nil {
			return nil, nil, fmt.Errorf("筛选符合端口空余的主机失败: %v", err)
		}
		// 走端口规则返回符合条件的服务器
		resParam = &api.RunSSHCmdAsyncReq{}
		resParam.Key = user.PriKey
		resParam.Passphrase = user.KeyPasswd
		for _, host := range hostList {
			resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
			resParam.Username = append(resParam.Username, host.User)
			resParam.SSHPort = append(resParam.SSHPort, host.Port)
			if pathCount != 0 && task.ConfigTem != "" {
				resConfig.HostIp = append(resConfig.HostIp, host.Ipv4.String)
				resConfig.Username = append(resConfig.Username, host.User)
				resConfig.SSHPort = append(resConfig.SSHPort, host.Port)
			}

		}

		if len(hostList) != serverSum {
			return nil, nil, errors.New("可用服务器数量 和 path参数的总数 不等，请检查服务器资源是否足够")
		}

		// 对模板进行渲染
		cmd, config, err := s.templateRender(&task, &args)
		if err != nil {
			return nil, nil, fmt.Errorf("cmdTem/configTem 渲染变量失败: %v", err)
		}
		resParam.Cmd = cmd

		if config != nil {
			resConfig.FileContent = config
			resConfig.Path = args["path"]
		}

		if err = service.SSH().CheckSSHParam(resParam); err != nil {
			return nil, nil, err
		}
		return resParam, resConfig, err
	}

	// 不走端口规则，但有path参数，过滤至path总数的可用服务器
	if pathCount != 0 {
		needHosts := hosts[:pathCount]
		// 清空变量
		resParam = &api.RunSSHCmdAsyncReq{}
		resParam.Key = user.PriKey
		resParam.Passphrase = user.KeyPasswd
		for _, host := range needHosts {
			resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
			resParam.Username = append(resParam.Username, host.User)
			resParam.SSHPort = append(resParam.SSHPort, host.Port)
			if task.ConfigTem != "" {
				resConfig.HostIp = append(resConfig.HostIp, host.Ipv4.String)
				resConfig.Username = append(resConfig.Username, host.User)
				resConfig.SSHPort = append(resConfig.SSHPort, host.Port)
			}
		}
		// 不走端口规则返回全部符合条件的服务器
	} else {
		for _, host := range hosts {
			resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
			resParam.Username = append(resParam.Username, host.User)
			resParam.SSHPort = append(resParam.SSHPort, host.Port)
		}
	}

	cmd, config, err := s.templateRender(&task, &args)

	if err != nil {
		return nil, nil, fmt.Errorf("cmdTem/configTem 渲染变量失败: %v", err)
	}
	resParam.Cmd = cmd
	if pathCount != 0 && config != nil {
		resConfig.FileContent = config
		resConfig.Path = args["path"]
	}

	if err = service.SSH().CheckSSHParam(resParam); err != nil {
		return nil, nil, err
	}
	return resParam, resConfig, err
}
