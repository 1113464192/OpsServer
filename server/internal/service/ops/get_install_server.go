package ops

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
)

func (s *OpsService) getInstallServer(hostList *[]model.Host, task *model.TaskTemplate, user *model.User, pathCount int, args *map[string][]string, configParam *api.RunSFTPAsyncReq) (sshReq *api.RunSSHCmdAsyncReq, sftpReq *api.RunSFTPAsyncReq, err error) {
	if task.ConfigTem == "" {
		return nil, nil, errors.New("装服报错: 配置文件模板为空")
	}
	// 如果没有端口规则, 那么控制一下hosts数量
	if len(*hostList) > pathCount {
		*hostList = (*hostList)[:pathCount]
	}
	// 走端口规则返回符合条件的服务器,因此要清空sshReq的数据
	sshReq = &api.RunSSHCmdAsyncReq{}
	sshReq.Key = user.PriKey
	sshReq.Passphrase = user.Passphrase

	for _, host := range *hostList {
		sshReq.HostIp = append(sshReq.HostIp, host.Ipv4.String)
		sshReq.Username = append(sshReq.Username, host.User)
		sshReq.SSHPort = append(sshReq.SSHPort, host.Port)

		configParam.HostIp = append(configParam.HostIp, host.Ipv4.String)
		configParam.Username = append(configParam.Username, host.User)
		configParam.SSHPort = append(configParam.SSHPort, host.Port)

	}

	if len(*hostList) != pathCount {
		return nil, nil, errors.New("可用服务器数量 和 path参数的总数 不等，请检查服务器资源是否足够")
	}

	// 对模板进行渲染
	cmd, config, err := s.templateRender(task, args)
	if err != nil {
		return nil, nil, fmt.Errorf("cmdTem/configTem 渲染变量失败: %v", err)
	}
	sshReq.Cmd = cmd
	configParam.FileContent = config
	configParam.Path = (*args)["path"]

	if err = service.SSH().CheckSSHParam(sshReq); err != nil {
		return nil, nil, err
	}
	return sshReq, configParam, err
}
