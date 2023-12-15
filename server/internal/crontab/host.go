package crontab

import (
	"fmt"
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
		// 接入微信小程序之类的请求, 向运维发送，马上处理定时主机信息获取失败的问题
		fmt.Println("微信小程序=====向运维发送，马上处理定时主机信息获取失败的问题")
		return
	}
	var hosts []model.Host
	if err := model.DB.Find(&hosts).Error; err != nil {
		logger.Log().Error("Host", "机器数据采集——获取主机对象失败", err)
		// 接入微信小程序之类的请求, 向运维发送，马上处理定时主机信息获取失败的问题
		fmt.Println("微信小程序=====向运维发送，马上处理定时主机信息获取失败的问题")
		return
	}

	var sshParam []api.SSHExecReq
	var req api.SSHExecReq
	for i := 0; i < len(hosts); i++ {
		req = api.SSHExecReq{
			HostIp:     hosts[i].Ipv4.String,
			Username:   hosts[i].User,
			SSHPort:    hosts[i].Port,
			Key:        opsUser.PriKey,
			Passphrase: opsUser.Passphrase,
			Cmd:        `ifconfig eth0 | grep inet`,
		}
		sshParam = append(sshParam, req)
	}
	hostInfo, err := service.Host().GetHostCurrData(&sshParam)
	if err != nil {
		logger.Log().Error("Host", "机器数据采集——数据结构有错误", err)
		// 接入微信小程序之类的请求, 向运维发送，马上处理定时主机信息获取失败的问题
		fmt.Println("微信小程序=====向运维发送，马上处理定时主机信息获取失败的问题")
		return
	}
	if err := service.Host().WritieToDatabase(hostInfo); err != nil {
		logger.Log().Error("Host", "机器数据采集——数据写入数据库失败", err)
		// 接入微信小程序之类的请求, 向运维发送，马上处理定时主机信息获取失败的问题
		fmt.Println("微信小程序=====向运维发送，马上处理定时主机信息获取失败的问题")
		return
	}
}
