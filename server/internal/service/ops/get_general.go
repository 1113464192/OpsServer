package ops

import (
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
)

func (s *OpsService) getGeneral(hosts *[]model.Host, task *model.TaskTemplate, args *map[string][]string, sshParam *api.RunSSHCmdAsyncReq, configParam *api.SftpReq) (resParam *api.RunSSHCmdAsyncReq, resConfig *api.SftpReq, err error) {
	// 不走端口规则，但有path参数，过滤至path总数的可用服务器
	// if pathCount != 0 {
	// 	needHosts := (*hosts)[:pathCount]
	// 	// 清空变量
	// 	resParam = &api.RunSSHCmdAsyncReq{}
	// 	resParam.Key = user.PriKey
	// 	resParam.Passphrase = user.KeyPasswd
	// 	for _, host := range needHosts {
	// 		resParam.HostIp = append(resParam.HostIp, host.Ipv4.String)
	// 		resParam.Username = append(resParam.Username, host.User)
	// 		resParam.SSHPort = append(resParam.SSHPort, host.Port)
	// 		if task.ConfigTem != "" {
	// 			resConfig.HostIp = append(resConfig.HostIp, host.Ipv4.String)
	// 			resConfig.Username = append(resConfig.Username, host.User)
	// 			resConfig.SSHPort = append(resConfig.SSHPort, host.Port)
	// 		}
	// 	}
	// 	// 不走端口规则返回全部符合条件的服务器
	// }
	for _, host := range *hosts {
		sshParam.HostIp = append(sshParam.HostIp, host.Ipv4.String)
		sshParam.Username = append(sshParam.Username, host.User)
		sshParam.SSHPort = append(sshParam.SSHPort, host.Port)
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
