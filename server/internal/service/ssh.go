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
	"os"
	"sync"

	gossh "golang.org/x/crypto/ssh"
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
	var sshReq []api.SSHClientConfigReq
	for i := 0; i < len(hosts); i++ {
		sshReq[i].HostIp = hosts[i].Ipv4.String
		sshReq[i].Username = hosts[i].User
		sshReq[i].SSHPort = hosts[i].Port
		sshReq[i].Key = user.PriKey
		sshReq[i].Passphrase = user.Passphrase
		sshReq[i].Cmd = `ifconfig eth0 | grep inet`
	}

	hostInfo, err := Host().GetHostCurrData(&sshReq)
	if err != nil {
		logger.Log().Error("Host", "机器数据采集——数据结构有错误", err)
		return nil, fmt.Errorf("机器数据采集——数据结构有错误: %v", err)
	}
	if err := Host().WritieToDatabase(hostInfo); err != nil {
		logger.Log().Error("Host", "机器数据采集——数据写入数据库失败", err)
		return nil, fmt.Errorf("机器数据采集——数据写入数据库失败: %v", err)
	}
	result, err = s.RunSSHCmdAsync(&sshReq)
	if err != nil {
		return nil, fmt.Errorf("测试执行失败: %v", err)
	}
	return result, err
}

type clientGroup struct {
	clientMap      map[string]*gossh.Client
	clientMapMutex sync.Mutex
}

func (s *SSHServer) RunSSHCmdAsync(param *[]api.SSHClientConfigReq) (*[]api.SSHResultRes, error) {
	if err := s.CheckSSHParam(param); err != nil {
		return nil, err
	}

	insClientGroup := clientGroup{
		clientMap:      make(map[string]*gossh.Client),
		clientMapMutex: sync.Mutex{},
	}

	channel := make(chan *api.SSHResultRes, len(*param))
	wg := sync.WaitGroup{}
	var err error
	var result []api.SSHResultRes
	// data := make(map[string]string)
	for i := 0; i < len(*param); i++ {
		if err = configs.Sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("获取信号失败，错误为: %v", err)
		}
		wg.Add(1)
		go s.RunSSHCmd(&(*param)[i], channel, &wg, &insClientGroup)
	}
	wg.Wait()
	close(channel)
	for res := range channel {
		result = append(result, *res)
	}
	return &result, err
}

func (s *SSHServer) RunSSHCmd(param *api.SSHClientConfigReq, ch chan *api.SSHResultRes, wg *sync.WaitGroup, insClientGroup *clientGroup) {
	result := &api.SSHResultRes{
		HostIp: param.HostIp,
		Status: true,
	}

	// client, err := ssh.SSHNewClient(param.HostIp, param.Username, param.SSHPort, param.Password, param.Key, param.Passphrase)
	client, err := s.getSSHClient(param.HostIp, param, insClientGroup)
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("建立/获取SSH客户端错误: %s", err.Error())
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
func (s *SSHServer) CheckSSHParam(param *[]api.SSHClientConfigReq) error {
	for _, p := range *param {
		if p.HostIp == "" || p.Username == "" || p.SSHPort == "" || p.Cmd == "" || p.Key == nil || p.Passphrase == nil {
			return fmt.Errorf("执行参数中存在空值: \n%v", p)
		}
	}
	return nil
}

func (s *SSHServer) RunSFTPAsync(param *[]api.SFTPClientConfigReq) (*[]api.SSHResultRes, error) {
	if err := s.CheckSFTPParam(param); err != nil {
		return nil, err
	}

	insClientGroup := clientGroup{
		clientMap:      make(map[string]*gossh.Client),
		clientMapMutex: sync.Mutex{},
	}

	channel := make(chan *api.SSHResultRes, len(*param))
	wg := sync.WaitGroup{}
	var err error
	var result []api.SSHResultRes
	// data := make(map[string]string)
	for i := 0; i < len(*param); i++ {
		if err = configs.Sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("获取信号失败，错误为: %v", err)
		}
		wg.Add(1)
		go s.RunSFTPTransfer(&(*param)[i], channel, &wg, &insClientGroup)
	}
	wg.Wait()
	close(channel)
	for res := range channel {
		result = append(result, *res)
	}
	return &result, err
}

func (s *SSHServer) RunSFTPTransfer(param *api.SFTPClientConfigReq, ch chan *api.SSHResultRes, wg *sync.WaitGroup, insClientGroup *clientGroup) {
	result := &api.SSHResultRes{
		HostIp: param.HostIp,
		Status: true,
	}

	// client, err := ssh.SSHNewClient(param.HostIp, param.Username, param.SSHPort, param.Password, param.Key, param.Passphrase)
	client, err := s.getSSHClient(param.HostIp, param, insClientGroup)
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("建立/获取SSH客户端错误: %s", err.Error())
		ch <- result
		wg.Done()
		configs.Sem.Release(1)
		return
	}
	defer client.Close()
	sftpClient, err := ssh.CreateSFTPClient(client)
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("建立SSH会话错误: %s", err.Error())
		ch <- result
		wg.Done()
		configs.Sem.Release(1)
		return
	}
	defer sftpClient.Close()
	remoteFile, err := sftpClient.OpenFile(param.Path, os.O_WRONLY|os.O_CREATE)
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("开启文件失败: %s ", err.Error())
		ch <- result
		wg.Done()
		configs.Sem.Release(1)
		return
	}
	defer remoteFile.Close()
	var bytesWritten int
	bytesWritten, err = remoteFile.Write([]byte(param.FileContent))
	if err != nil {
		result.Status = false
		result.Response = fmt.Sprintf("写入文件内容到文件失败: %s", err.Error())
		ch <- result
		wg.Done()
		configs.Sem.Release(1)
		return
	}

	result.Response = fmt.Sprintf("写入文件字节数为: %d", bytesWritten)
	ch <- result
	wg.Done()
	configs.Sem.Release(1)
}

// 检查是否符合执行条件
func (s *SSHServer) CheckSFTPParam(param *[]api.SFTPClientConfigReq) error {
	for _, p := range *param {
		if p.HostIp == "" || p.Username == "" || p.SSHPort == "" || p.FileContent == "" || p.Key == nil || p.Passphrase == nil || p.Path == "" {
			return fmt.Errorf("执行参数中存在空值: \n%v", p)
		}
	}
	return nil
}

func (s *SSHServer) getSSHClient(hostIp string, param any, insClientGroup *clientGroup) (client *gossh.Client, err error) {
	insClientGroup.clientMapMutex.Lock()
	defer insClientGroup.clientMapMutex.Unlock()

	var ok bool
	// 检查Map中是否已经存在对应hostIp的client
	if client, ok = insClientGroup.clientMap[hostIp]; ok {
		// 如果存在，则直接返回已有的client
		return client, err
	}

	if params, ok := param.(*api.SSHClientConfigReq); ok {
		client, err = ssh.SSHNewClient(params.HostIp, params.Username, params.SSHPort, params.Password, params.Key, params.Passphrase)

	}
	if params, ok := param.(*api.SFTPClientConfigReq); ok {
		client, err = ssh.SSHNewClient(params.HostIp, params.Username, params.SSHPort, params.Password, params.Key, params.Passphrase)
	}
	if client == nil {
		return nil, errors.New("未能成功获取到ssh.Client")
	}
	insClientGroup.clientMap[hostIp] = client
	return client, err
}
