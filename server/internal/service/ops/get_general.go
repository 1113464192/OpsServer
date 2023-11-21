package ops

// import (
// 	"errors"
// 	"fmt"
// 	"fqhWeb/internal/model"
// 	"fqhWeb/internal/service"
// 	"fqhWeb/pkg/api"
// )

// func (s *OpsService) getGeneral(hosts *[]model.Host, user *model.User, task *model.TaskTemplate, args *map[string][]string, pathCount int) (sshReq *[]api.SSHClientConfigReq, sftpReq *[]api.SFTPClientConfigReq, err error) {
// 	// 如果没有端口规则, 那么控制一下hosts数量
// 	if len(*hosts) > pathCount {
// 		*hosts = (*hosts)[:pathCount]
// 	}
// 	// 判断两者是否相等
// 	if len(*hosts) != pathCount {
// 		return nil, nil, errors.New("可用服务器数量 和 path参数的总数 不等，请检查服务器资源是否足够")
// 	}
// 	// 走端口规则返回符合条件的服务器,因此要清空sshReq的数据
// 	if task.CmdTem == "" || task.ConfigTem == "" {
// 		return nil, nil, errors.New("执行命令与配置文件模板都为空")
// 	}
// 	for _, host := range *hosts {
// 		if task.CmdTem == "" && task.ConfigTem == "" {
// 			return nil, nil, errors.New("任务的命令和传输文件内容都为空")
// 		}
// 		if task.CmdTem != "" {
// 			(*sshReq)[i].HostIp = (*hostList)[i].Ipv4.String
// 			(*sshReq)[i].Username = (*hostList)[i].User
// 			(*sshReq)[i].SSHPort = (*hostList)[i].Port
// 			(*sshReq)[i].Key = user.PriKey
// 			(*sshReq)[i].Passphrase = user.Passphrase
// 		}
// 		if task.ConfigTem != "" {
// 			configParam.HostIp = append(configParam.HostIp, host.Ipv4.String)
// 			configParam.Username = append(configParam.Username, host.User)
// 			configParam.SSHPort = append(configParam.SSHPort, host.Port)
// 		}
// 	}

// 	cmd, config, err := s.templateRender(task, args)

// 	if err != nil {
// 		return nil, nil, fmt.Errorf("cmdTem/configTem 渲染变量失败: %v", err)
// 	}
// 	sshParam.Cmd = cmd
// 	if config != nil {
// 		configParam.FileContent = config
// 		configParam.Path = (*args)["path"]
// 	}

// 	if err = service.SSH().CheckSSHParam(sshParam); err != nil {
// 		return nil, nil, err
// 	}
// 	return sshParam, configParam, err
// }
