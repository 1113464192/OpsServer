package ops

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
)

func (s *OpsService) getInstallServer(hostList *[]model.Host, task *model.TaskTemplate, user *model.User, pathCount int, args *map[string][]string) (sshReq *[]api.SSHClientConfigReq, sftpReq *[]api.SFTPClientConfigReq, err error) {
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
	for i := 0; i < len(*hostList); i++ {
		(*sshReq)[i].HostIp = (*hostList)[i].Ipv4.String
		(*sshReq)[i].Username = (*hostList)[i].User
		(*sshReq)[i].SSHPort = (*hostList)[i].Port
		(*sshReq)[i].Key = user.PriKey
		(*sshReq)[i].Passphrase = user.Passphrase
		(*sftpReq)[i].HostIp = (*hostList)[i].Ipv4.String
		(*sftpReq)[i].Username = (*hostList)[i].User
		(*sftpReq)[i].SSHPort = (*hostList)[i].Port
		(*sftpReq)[i].Key = user.PriKey
		(*sftpReq)[i].Passphrase = user.Passphrase
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
