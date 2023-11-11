package service

import (
	"context"
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils/ssh"
	"sync"
)

type SSHServer struct {
}

var (
	insSSH = &SSHServer{}
)

func SSH() *SSHServer {
	return insSSH
}

func (s *SSHServer) TestSSH(param api.TestSSHReq) (*[]api.SSHResultRes, error) {
	var result *[]api.SSHResultRes
	var user model.User
	var hosts []model.Host
	if err := model.DB.First(&user, param.UserId).Error; err != nil {
		return nil, fmt.Errorf("GORM用户未找到: %v", err)
	}
	if err := model.DB.Find(&hosts, param.HostId).Error; err != nil {
		return nil, fmt.Errorf("GORM服务器未找到: %v", err)
	}
	var hostIp []string
	var hostUsers []string
	var ports []string
	for _, host := range hosts {
		hostIp = append(hostIp, host.Ipv4.String)
		hostUsers = append(hostUsers, host.User)
		ports = append(ports, host.Port)
	}
	sshReq := &api.RunSSHCmdAsyncReq{
		HostIp:     hostIp,
		Username:   hostUsers,
		SSHPort:    ports,
		Key:        user.PriKey,
		Passphrase: user.KeyPasswd,
	}
	hostInfo, err := Host().GetHostCurrData(sshReq)
	if err != nil {
		logger.Log().Error("Host", "机器数据采集——数据结构有错误", err)
		return nil, fmt.Errorf("机器数据采集——数据结构有错误: %v", err)
	}
	if err := Host().WritieToDatabase(hostInfo); err != nil {
		logger.Log().Error("Host", "机器数据采集——数据写入数据库失败", err)
		return nil, fmt.Errorf("机器数据采集——数据写入数据库失败: %v", err)
	}
	sshReq.Cmd = []string{`ifconfig eth0 | grep inet`}
	result, err = s.RunSSHCmdAsync(sshReq)
	if err != nil {
		return nil, fmt.Errorf("测试执行失败: %v", err)
	}
	return result, err
}

func (s *SSHServer) RunSSHCmdAsync(param *api.RunSSHCmdAsyncReq) (*[]api.SSHResultRes, error) {
	if err := s.CheckSSHParam(param); err != nil {
		return nil, err
	}

	channel := make(chan *api.SSHResultRes, len(param.HostIp))
	wg := sync.WaitGroup{}
	var err error
	var result []api.SSHResultRes
	// data := make(map[string]string)
	for i := 0; i < len(param.HostIp); i++ {
		if err = configs.Sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("获取信号失败，错误为: %v", err)
		}
		wg.Add(1)
		password := param.Password[param.HostIp[i]]
		sshParam := &api.SSHClientConfigReq{
			HostIp:     param.HostIp[i],
			Username:   param.Username[i],
			SSHPort:    param.SSHPort[i],
			Password:   password,
			Key:        param.Key,
			Passphrase: param.Passphrase,
			Cmd:        param.Cmd[0],
		}
		if len(param.Cmd) > 1 {
			sshParam.Cmd = param.Cmd[i]
		}
		go s.RunSSHCmd(sshParam, channel, &wg)
	}
	wg.Wait()
	close(channel)
	for res := range channel {
		result = append(result, *res)
	}
	return &result, err
}

func (s *SSHServer) RunSSHCmd(param *api.SSHClientConfigReq, ch chan *api.SSHResultRes, wg *sync.WaitGroup) {
	result := &api.SSHResultRes{
		HostIp: param.HostIp,
		Status: true,
	}

	client, err := ssh.SSHNewClient(param)
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("建立SSH客户端错误: %s", err.Error())
		ch <- result
		wg.Done()
		configs.Sem.Release(1)
		return
	}
	defer client.Close()
	session, err := ssh.SSHNewSession(client)
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("建立SSH会话错误: %s", err.Error())
		ch <- result
		wg.Done()
		configs.Sem.Release(1)
		return
	}
	defer session.Close()
	output, err := session.CombinedOutput(param.Cmd)
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("Failed to execute command: %s %s", string(output), err.Error())
		ch <- result
		wg.Done()
		configs.Sem.Release(1)
		return
	}

	result.Response = string(output)
	ch <- result
	wg.Done()
	configs.Sem.Release(1)
}

// 检查是否符合执行条件
func (s *SSHServer) CheckSSHParam(param *api.RunSSHCmdAsyncReq) error {
	fmt.Println(param.HostIp)
	fmt.Println(param.Username)
	fmt.Println(param.SSHPort)
	fmt.Println(param.Cmd)
	for _, value := range param.HostIp {
		if value == "" {
			// 可能path大于hosts数量
			return fmt.Errorf("IP切片中有值为空字符串")
		}
	}
	for _, value := range param.SSHPort {
		if value == "" {
			return fmt.Errorf("端口切片中有值为空字符串")
		}
	}
	for _, value := range param.Username {
		if value == "" {
			return fmt.Errorf("用户切片中有值为空字符串")
		}
	}
	for _, value := range param.Cmd {
		if value == "" {
			return fmt.Errorf("CMD切片中有值为空字符串")
		}
	}
	if len(param.HostIp) != len(param.SSHPort) || len(param.HostIp) != len(param.Username) || len(param.Cmd) > 1 && len(param.Cmd) != len(param.HostIp) {
		return errors.New("请检查: IP、端口、用户名这些切片是否一一对应")
	}
	return nil
}
