package ops

import (
	"errors"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"strconv"
)

func (s *OpsService) getInstallServerParam(hostList *[]model.Host, task *model.TaskTemplate, user *model.User, pathCount int, args *map[string][]string) (sshReq *[]api.SSHClientConfigReq, sftpReq *[]api.SFTPClientConfigReq, err error) {
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
	sshReq = &[]api.SSHClientConfigReq{}
	sftpReq = &[]api.SFTPClientConfigReq{}
	var sReq api.SSHClientConfigReq
	var fReq api.SFTPClientConfigReq
	for i := 0; i < len(*hostList); i++ {
		sReq = api.SSHClientConfigReq{
			HostIp:     (*hostList)[i].Ipv4.String,
			Username:   (*hostList)[i].User,
			SSHPort:    (*hostList)[i].Port,
			Key:        user.PriKey,
			Passphrase: user.Passphrase,
		}
		*sshReq = append(*sshReq, sReq)
		fReq = api.SFTPClientConfigReq{
			HostIp:     (*hostList)[i].Ipv4.String,
			Username:   (*hostList)[i].User,
			SSHPort:    (*hostList)[i].Port,
			Key:        user.PriKey,
			Passphrase: user.Passphrase,
		}
		*sftpReq = append(*sftpReq, fReq)
	}

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
			(*sftpReq)[i].Path = (*args)["path"][i]
		}
	} else {
		for i := 0; i < len(*sftpReq); i++ {
			(*sftpReq)[i].FileContent = config[0]
			(*sftpReq)[i].Path = (*args)["path"][i]
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
func (s *OpsService) opsInstallServer(pathCount int, task *model.TaskTemplate, hosts *[]model.Host, user *model.User, args *map[string][]string, sshReq *[]api.SSHClientConfigReq) (*[]api.SSHClientConfigReq, *[]api.SFTPClientConfigReq, error) {
	var err error
	if pathCount == 0 {
		return nil, nil, errors.New("path参数数量为0")
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
		hostList, err = s.filterPortRuleHost(hosts, user, task, sshReq, args, memSize)
		if err != nil {
			return nil, nil, fmt.Errorf("端口筛选报错: %v", err)
		}
		*hosts = *hostList
	}
	// 记录所有hostid
	var hostIds []string
	var idString string
	for i := 0; i < len(*hosts); i++ {
		idString = strconv.Itoa(int((*hosts)[i].ID))
		hostIds = append(hostIds, idString)
	}
	(*args)["hostId"] = hostIds

	var sftpReq *[]api.SFTPClientConfigReq
	sshReq, sftpReq, err = s.getInstallServerParam(hosts, task, user, pathCount, args)
	if err != nil {
		return nil, nil, fmt.Errorf("获取%s参数报错: %v", consts.OperationInstallServerType, err)
	}
	return sshReq, sftpReq, err
}
