package crontab

import (
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
)

func CronWrittenHostInfo() {
	// 设定指定的用户，一般设置为高权限用户的私钥来执行全机器数据采集，这里设置为1
	var opsUser model.User
	if err := model.DB.First(&opsUser, consts.SSHOpsUserId).Error; err != nil {
		logger.Log().Error("User", "机器数据采集——获取OPS用户权限失败", err)
		return
	}
	var hosts []model.Host
	if err := model.DB.Find(&hosts).Error; err != nil {
		logger.Log().Error("Host", "机器数据采集——获取主机对象失败", err)
		return
	}

	var sshParam = api.RunSSHCmdAsyncReq{
		Key:        opsUser.PriKey,
		Passphrase: opsUser.KeyPasswd,
	}
	for _, host := range hosts {
		sshParam.HostIp = append(sshParam.HostIp, host.Ipv4.String)
		sshParam.Username = append(sshParam.Username, host.User)
		sshParam.SSHPort = append(sshParam.SSHPort, host.Port)
	}
	hostInfo, err := service.Host().GetHostCurrData(&sshParam)
	if err != nil {
		logger.Log().Error("Host", "机器数据采集——数据结构有错误", err)
		return
	}
	if err := service.Host().WritieToDatabase(hostInfo); err != nil {
		logger.Log().Error("Host", "机器数据采集——数据写入数据库失败", err)
		return
	}
}
