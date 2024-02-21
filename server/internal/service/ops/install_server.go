package ops

import (
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/cloudScript/tencentCloud"
	"strconv"
	"strings"
)

func (s *OpsService) getInstallServerParam(hostList *[]model.Host, task *model.TaskTemplate, user *model.User, pathCount int, args *map[string][]string) (sshReq *[]api.SSHExecReq, sftpReq *[]api.SFTPExecReq, err error) {
	// 如果没有端口规则, 那么控制一下hosts数量
	if len(*hostList) > pathCount {
		*hostList = (*hostList)[:pathCount]
	}
	// 判断两者是否相等
	if len(*hostList) != pathCount {
		return nil, nil, errors.New("可用服务器数量 和 path参数的总数 不等，请检查服务器资源是否足够")
	}
	// 走端口规则返回符合条件的服务器,因此要清空sshReq的数据
	if task.CmdTem == "" || task.ConfigTem == "" {
		return nil, nil, errors.New("执行命令与配置文件模板都为空")
	}
	sshReq = &[]api.SSHExecReq{}
	sftpReq = &[]api.SFTPExecReq{}
	var sReq api.SSHExecReq
	var fReq api.SFTPExecReq
	for i := 0; i < len(*hostList); i++ {
		sReq = api.SSHExecReq{
			HostIp:     (*hostList)[i].Ipv4.String,
			Username:   (*hostList)[i].User,
			SSHPort:    (*hostList)[i].Port,
			Key:        user.PriKey,
			Passphrase: user.Passphrase,
		}
		*sshReq = append(*sshReq, sReq)
		fReq = api.SFTPExecReq{
			HostIp:     (*hostList)[i].Ipv4.String,
			Username:   (*hostList)[i].User,
			SSHPort:    (*hostList)[i].Port,
			Key:        user.PriKey,
			Passphrase: user.Passphrase,
		}
		*sftpReq = append(*sftpReq, fReq)
	}
	// 加入tid到cmd的渲染中，方便后续client更改task状态
	(*args)["tid"] = []string{strconv.Itoa(int(task.ID))}

	// 对模板进行渲染
	cmd, config, err := s.templateRender(task, args, pathCount)
	if err != nil {
		return nil, nil, fmt.Errorf("cmdTem/configTem 渲染变量失败: %v", err)
	}
	// 判断是每个单服一套命令还是不同命令
	if len(cmd) == pathCount {
		for i := 0; i < len(*sshReq); i++ {
			(*sshReq)[i].Cmd = cmd[i]
		}
	} else {
		for i := 0; i < len(*sshReq); i++ {
			(*sshReq)[i].Cmd = cmd[0]
		}
	}
	// 判断是每个单服一套配置还是不同配置
	if len(config) == pathCount {
		for i := 0; i < len(*sftpReq); i++ {
			(*sftpReq)[i].FileContent = config[i]
			(*sftpReq)[i].Path = (*args)["sftpPath"][i]
		}
	} else {
		for i := 0; i < len(*sftpReq); i++ {
			(*sftpReq)[i].FileContent = config[0]
			(*sftpReq)[i].Path = (*args)["sftpPath"][i]
		}
	}

	if err = service.SSH().CheckSSHParam(sshReq); err != nil {
		return nil, nil, err
	}
	if err = service.SSH().CheckSFTPParam(sftpReq); err != nil {
		return nil, nil, err
	}
	return sshReq, sftpReq, err
}

// 装服传参请包含path、serverName、端口规则(key要规则名)
func (s *OpsService) opsInstallServer(pathCount int, task *model.TaskTemplate, hosts *[]model.Host, user *model.User, args *map[string][]string, sshReq *[]api.SSHExecReq) (*[]api.SSHExecReq, *[]api.SFTPExecReq, error) {
	var err error
	if pathCount == 0 {
		return nil, nil, errors.New("path参数数量为0")
	}
	if pathCount != len((*args)["sftpPath"]) {
		return nil, nil, errors.New("path和sftpPath参数数量不对等")
	}
	if task.CmdTem == "" || task.ConfigTem == "" {
		return nil, nil, errors.New("任务的命令和传输文件内容都为空")
	}
	var hostList *[]model.Host
	// 如果有设置条件 则筛选符合条件的主机
	var memSize float32
	if task.Condition != "" {
		if err = s.filterConditionHost(hosts, user, task, sshReq, &memSize); err != nil {
			return nil, nil, fmt.Errorf("筛选符合条件的主机失败: %v", err)
		}
	}
	// 筛选端口规则
	if task.PortRule != "" {
		var flagTotal, flagAssHostSum int
		for i := 0; i < configs.Conf.Cloud.AllowConsecutiveCreateTimes; i++ { // 最多尝试AllowConsecutiveCreateTimes次
			hostList, flagTotal, flagAssHostSum, err = s.filterPortRuleHost(hosts, user, task, sshReq, args, memSize)
			if flagAssHostSum != flagTotal {
				if err = s.autoCreateInstance(task); err != nil {
					return nil, nil, fmt.Errorf("自动购买服务器失败: %v", err)
				}
				// 服务器购买成功，继续下一次循环尝试筛选
				continue
			} else if err != nil {
				return nil, nil, fmt.Errorf("端口筛选报错: %v", err)
			}

			*hosts = *hostList
			break // 筛选成功，跳出循环
		}
		if flagAssHostSum != flagTotal {
			// 接入微信小程序之类的请求,向对应运维发送
			fmt.Printf("微信小程序=====通知运维购买%d次机器仍然无法装服===========", configs.Conf.Cloud.AllowConsecutiveCreateTimes)
		}
	}

	// 记录所有hostid
	var hostIds []string
	var idString string
	for i := 0; i < len(*hosts); i++ {
		idString = strconv.Itoa(int((*hosts)[i].ID))
		hostIds = append(hostIds, idString)
	}
	(*args)["hostId"] = hostIds

	var sftpReq *[]api.SFTPExecReq
	sshReq, sftpReq, err = s.getInstallServerParam(hosts, task, user, pathCount, args)
	if err != nil {
		return nil, nil, fmt.Errorf("获取%s参数报错: %v", consts.OperationInstallServerType, err)
	}
	return sshReq, sftpReq, err
}

// 自动购买服务器
func (s *OpsService) autoCreateInstance(task *model.TaskTemplate) (err error) {
	// 获取配置
	insConfig, err := service.Cloud().GetCloudInstanceConfig(task.Pid)
	if err != nil {
		return fmt.Errorf("获取云服务器配置失败: %v", err)
	}
	insTypeInterface, err := service.Cloud().GetCloudInstanceTypeList(task.Project.Cloud, insConfig.Region, insConfig.InstanceFamily, insConfig.CpuCores, insConfig.MemorySize, insConfig.Fpga, insConfig.GpuCores)
	if err != nil {
		return fmt.Errorf("获取云服务器类型失败: %v", err)
	}
	//var (
	//	insTypeTencentCloud tencentCloud.InstanceConfigRes
	//)
	var (
		projectHosts *[]model.Host
		maxFlag      int
	)
	if err = model.DB.Model(&task.Project).Association("Hosts").Find(projectHosts); err != nil {
		return fmt.Errorf("获取项目下的服务器失败: %v", err)
	}
	// 取出最大的flag
	for _, projectHost := range *projectHosts {
		numberPart := strings.TrimPrefix(projectHost.Name, task.Project.Name+"-")
		number, err := strconv.Atoi(numberPart)
		if err != nil {
			return fmt.Errorf("获取服务器编号失败: %v", err)
		}
		if number > maxFlag {
			maxFlag = number
		}
	}

	instanceName := task.Project.Name + "-" + strconv.Itoa(maxFlag+1)
	switch task.Project.Cloud {
	case "腾讯云":
		insType, ok := insTypeInterface.(tencentCloud.InstanceConfigRes)
		if !ok {
			return fmt.Errorf("获取云服务器类型失败: %v", err)
		}
		// 购买服务器
		//cloudType string, region string, instanceChargeType string, period int64, renewFlag string, zone string, projectId int64,
		//	instanceType string, imageId string, systemDiskType string, systemDiskSize int64, dataDiskType string, dataDiskSize int64, vpcId string,
		//	subnetId string, internetChargeType string, internetMaxBandwidthOut int64, instanceName string, securityGroupId string, hostName string)
		if err = service.Cloud().CreateCloudInstance(task.Project.Cloud, insConfig.Region, insConfig.InstanceChargeType, insConfig.Period, insConfig.RenewFlag, insType.Zone,
			insConfig.ProjectId, insType.InstanceType, insConfig.ImageId, insConfig.SystemDiskType, insConfig.SystemDiskSize, insConfig.DataDiskType, insConfig.DataDiskSize,
			insConfig.VpcId, insConfig.SubnetId, insConfig.InternetChargeType, insConfig.InternetMaxBandwidthOut, instanceName, insConfig.SecurityGroupId, instanceName); err != nil {
			return fmt.Errorf("购买云服务器失败: %v", err)
		}
	default:
		return fmt.Errorf("不支持的云服务器类型: %v", task.Project.Cloud)
	}
	return err
}

// 写入host表
func (s *OpsService) writeHostTable(insTypeInterface any, insName string, project *model.Project, insConfig *model.CloudInstanceConfig) (err error) {
	// 写入host表
	//ID       uint   `form:"id" json:"id"`
	//Ipv4     string `form:"ipv4" json:"ipv4" binding:"required"`
	//Ipv6     string `form:"ipv6" json:"ipv6"`
	//Name     string `form:"name" json:"name" binding:"required"`
	//User     string `form:"user" json:"user" binding:"required"`
	//Password []byte `form:"password" json:"password"`
	//Port     string `form:"port" json:"port" binding:"required"`
	//Zone     string `form:"zone" json:"zone" binding:"required"`           // 所在地，用英文小写，如guangzhou、seoul
	//ZoneTime uint8  `form:"zone_time" json:"zone_time" binding:"required"` // 时区，如东八区填8
	////BillingType uint8   `form:"billing" json:"billing" binding:"required"`     // 1 按量收费, 2 包月收费, 3 包年收费 ...后续有需要再加
	//Cost       float32 `form:"cost" json:"cost"` // 下次续费金额, 人民币为单位
	//Cloud      string  `form:"cloud" json:"cloud" binding:"required"`
	//System     string  `form:"system" json:"system" binding:"required"`
	//Iops       uint32  `form:"iops" json:"iops" binding:"required"`
	//Mbps       uint32  `form:"mbps" json:"mbps" binding:"required"`
	//Type       uint8   `form:"type" json:"type" binding:"required"`               // 1 单服机器, 2 中央服机器, 3 CDN机器, 4 业务服机器  ...后续有需要再加
	//Cores      uint16  `form:"cores" json:"cores" binding:"required"`             // 四核输入4
	//SystemDisk uint32  `form:"system_disk" json:"system_disk" binding:"required"` // 磁盘单位为G
	//DataDisk   uint32  `form:"data_disk" json:"data_disk" binding:"required"`     // 磁盘单位为G
	//Mem        uint32  `form:"mem" json:"mem" binding:"required"`                 // 内存单位为G

	insInfoInterface, err := service.Cloud().GetCloudInsInfo(project.Cloud, insConfig.Region, "", "", insName, "", 1, 1)
	if err != nil {
		return fmt.Errorf("获取云实例信息失败: %v", err)

	}
	switch project.Cloud {
	case "腾讯云":
		insType, ok := insTypeInterface.(tencentCloud.InstanceConfigRes)
		if !ok {
			return fmt.Errorf("获取云服务器类型失败: %v", err)
		}
		insInfo, ok := insInfoInterface.(tencentCloud.HostResponse)
		if !ok {
			return fmt.Errorf("断言云服务器信息失败: %v", err)
		}
		// 默认root用户并不提供密码
		hostReq := api.UpdateHostReq{
			Ipv4:     insInfo.CloudHostResponse.InstanceSet[0].PrivateIpAddresses[0],
			Ipv6:     insInfo.CloudHostResponse.InstanceSet[0].IPv6Addresses[0],
			Name:     insName,
			User:     consts.DefaultHostUsername,
			Password: []byte(consts.DefaultHostPassword),
			Port:     consts.DefaultHostPort,
			Zone:     insInfo.CloudHostResponse.InstanceSet[0].Placement.Zone,
		}
	default:
		return fmt.Errorf("不支持的云服务器类型: %v", task.Project.Cloud)
	}
	return err
}
