package ops

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
)

func (s *OpsService) getGeneral(hosts *[]model.Host, user *model.User, task *model.TaskTemplate, args *map[string][]string, sshParam *api.SSHClientConfigReq, configParam *api.SFTPClientConfigReq) (sshReq *api.RunSSHCmdAsyncReq, sftpReq *api.RunSFTPAsyncReq, err error) {
	// 不走端口规则，但有path参数，过滤至path总数的可用服务器
	// if pathCount != 0 {
	// 	needHosts := (*hosts)[:pathCount]
	// 	// 清空变量
	// 	sshReq = &api.RunSSHCmdAsyncReq{}
	// 	sshReq.Key = user.PriKey
	// 	sshReq.Passphrase = user.Passphrase
	// 	for _, host := range needHosts {
	// 		sshReq.HostIp = append(sshReq.HostIp, host.Ipv4.String)
	// 		sshReq.Username = append(sshReq.Username, host.User)
	// 		sshReq.SSHPort = append(sshReq.SSHPort, host.Port)
	// 		if task.ConfigTem != "" {
	// 			sftpReq.HostIp = append(sftpReq.HostIp, host.Ipv4.String)
	// 			sftpReq.Username = append(sftpReq.Username, host.User)
	// 			sftpReq.SSHPort = append(sftpReq.SSHPort, host.Port)
	// 		}
	// 	}
	// 	// 不走端口规则返回全部符合条件的服务器
	// }
	for _, host := range *hosts {
		if task.CmdTem == "" && task.ConfigTem == "" {
			return nil, nil, errors.New("任务的命令和传输文件内容都为空")
		}
		if task.CmdTem != "" {
			sshParam.HostIp = append(sshParam.HostIp, host.Ipv4.String)
			sshParam.Username = append(sshParam.Username, host.User)
			sshParam.SSHPort = append(sshParam.SSHPort, host.Port)
		}
		if task.ConfigTem != "" {
			configParam.HostIp = append(configParam.HostIp, host.Ipv4.String)
			configParam.Username = append(configParam.Username, host.User)
			configParam.SSHPort = append(configParam.SSHPort, host.Port)
		}
	}

	cmd, config, err := s.templateRender(task, args)

	if err != nil {
		return nil, nil, fmt.Errorf("cmdTem/configTem 渲染变量失败: %v", err)
	}
	sshParam.Cmd = cmd
	if config != nil {
		configParam.FileContent = config
		configParam.Path = (*args)["path"]
	}

	if err = service.SSH().CheckSSHParam(sshParam); err != nil {
		return nil, nil, err
	}
	return sshParam, configParam, err
}
