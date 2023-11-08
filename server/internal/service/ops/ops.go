package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

type OpsService struct {
}

var (
	insOps = &OpsService{}
)

func Ops() *OpsService {
	return insOps
}

func (s *OpsService) getFlag(param string) (flags []int, err error) {
	var args map[string][]string
	if err = json.Unmarshal([]byte(param), &args); err != nil {
		return flags, errors.New("参数字段进行json解析失败")
	}
	re := regexp.MustCompile(`\d+`)
	if len(args["path"]) != 0 {
		for _, path := range args["path"] {
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

func (s *OpsService) generatePortRuleSet(portRule map[int]string, flag int) ([]string, error) {
	var resultSet []string
	for _, rule := range portRule {
		// 判断端口规则是否规范
		if !strings.Contains(rule, "flag") {
			return nil, errors.New(rule + " 不包含 flag 字符串")
		}
		// 创建端口规则表达式
		expr, err := govaluate.NewEvaluableExpression(rule)
		if err != nil {
			return nil, fmt.Errorf("创建表达式解析器报错: %v", err)
		}
		vars := map[string]interface{}{
			"flag": flag,
		}
		// 获取出游戏服占用端口
		port, err := expr.Evaluate(vars)
		if err != nil {
			return nil, fmt.Errorf("表达式计算报错: %v", err)
		}
		cmdShell := fmt.Sprintf(`
		if [[ -z $(netstat -plan | grep %d) ]];then
			echo "success"
		fi`, port.(int))
		resultSet = append(resultSet, cmdShell)
	}
	return resultSet, nil
}

func (s *OpsService) filterPortRuleHost(hosts *[]model.Host, task *model.TaskTemplate, sshReq *api.RunSSHCmdAsyncReq, memSize float32) (hostList []model.Host, err error) {
	tmpHosts := make([]model.Host, len(*hosts))
	copy(tmpHosts, *hosts)

	var portRule map[int]string
	if err = json.Unmarshal([]byte(task.PortRule), &portRule); err != nil {
		return nil, errors.New("端口规则进行json解析失败")
	}
	flags, err := s.getFlag(task.Args)
	if err != nil {
		return nil, err
	}
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
			var cmdSet []string
			cmdSet, err = s.generatePortRuleSet(portRule, flag)
			if err != nil {
				continue
			}
			var sshResult *[]api.SSHResultRes
			count := len(cmdSet)
			num := 0
			for _, cmd := range cmdSet {
				sshResult, err = service.SSH().RunSSHCmdAsync(sshReq, cmd)
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

func (s *OpsService) filterConditionHost(hosts *[]model.Host, task *model.TaskTemplate, resParam *api.RunSSHCmdAsyncReq) (memSize float32, err error) {
	for _, host := range *hosts {
		resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
		resParam.Username = append(resParam.Username, host.User)
		resParam.SSHPort = append(resParam.SSHPort, host.Port)
	}
	hostInfo, err := service.Host().GetHostCurrData(resParam)
	if err != nil {
		logger.Log().Error("Host", "机器数据采集——数据结构有错误", err)
		return memSize, fmt.Errorf("机器数据采集——数据结构有错误: %v", err)
	}
	if err := service.Host().WritieToDatabase(hostInfo); err != nil {
		logger.Log().Error("Host", "机器数据采集——数据写入数据库失败", err)
		return memSize, fmt.Errorf("机器数据采集——数据写入数据库失败: %v", err)
	}

	var condition map[string][]string
	if err = json.Unmarshal([]byte(task.Condition), &condition); err != nil {
		return memSize, errors.New("筛选机器条件规则进行json解析失败")
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
				return 0, fmt.Errorf(value[0]+"转换为浮点数失败: %v", err)
			}
			memSize = float32(memFloat)

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
			return memSize, fmt.Errorf("%s 不属于ConditionSet中的任何一个", key)
		}
	}

	// 使用单个查询筛选符合条件的主机
	if len(fields) > 0 {
		conditions := strings.Join(fields, " AND ")
		if err = model.DB.Where(conditions, values...).Order("curr_mem").Find(hosts).Error; err != nil {
			return memSize, fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
		}
	}
	return memSize, err
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

func (s *OpsService) GetExecParam(param api.GetExecParamReq) (resParam *api.RunSSHCmdAsyncReq, err error) {
	var task model.TaskTemplate
	var user model.User
	if err = model.DB.First(&task, param.Tid).Error; err != nil {
		return nil, errors.New("根据id查询任务失败")
	}
	if err = model.DB.First(&user, param.Uid).Error; err != nil {
		return nil, errors.New("根据id查询用户失败")
	}
	resParam.Key = user.PriKey
	resParam.Passphrase = user.KeyPasswd
	var hosts []model.Host
	if err = model.DB.Model(&task).Association("Hosts").Find(&hosts); err != nil {
		return nil, errors.New("查询task关联主机失败")
	}

	var memSize float32
	if task.Condition != "" {
		if memSize, err = s.filterConditionHost(&hosts, &task, resParam); err != nil {
			return nil, fmt.Errorf("筛选符合条件的主机失败: %v", err)
		}
	}
	if task.PortRule != "" {
		if task.Args == "" {
			return nil, errors.New("有端口规则请传path, 否则无标识判断")
		}
		var hostList []model.Host
		if hostList, err = s.filterPortRuleHost(&hosts, &task, resParam, memSize); err != nil {
			return nil, fmt.Errorf("筛选符合端口空余的主机失败: %v", err)
		}
		// 走端口规则返回符合条件的服务器
		resParam = nil
		resParam.Key = user.PriKey
		resParam.Passphrase = user.KeyPasswd
		for _, host := range hostList {
			resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
			resParam.Username = append(resParam.Username, host.User)
			resParam.SSHPort = append(resParam.SSHPort, host.Port)
		}
		return resParam, err
	}
	// 不走端口规则返回全部符合条件的服务器
	for _, host := range hosts {
		resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
		resParam.Username = append(resParam.Username, host.User)
		resParam.SSHPort = append(resParam.SSHPort, host.Port)
	}
	return resParam, err

}
