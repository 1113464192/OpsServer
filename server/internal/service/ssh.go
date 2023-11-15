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

type clientGroup struct {
	clientMap      map[string]*gossh.Client
	clientMapMutex sync.Mutex
}

func (s *SSHServer) RunSSHCmdAsync(param *api.RunSSHCmdAsyncReq) (*[]api.SSHResultRes, error) {
	if err := s.CheckSSHParam(param); err != nil {
		return nil, err
	}

	insClientGroup := clientGroup{
		clientMap:      make(map[string]*gossh.Client),
		clientMapMutex: sync.Mutex{},
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
		// 如果是一个主机一个命令则多条
		if len(param.Cmd) > 1 {
			sshParam.Cmd = param.Cmd[i]
		}
		go s.RunSSHCmd(sshParam, channel, &wg, &insClientGroup)
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

func (s *SSHServer) RunSFTPAsync(param *api.RunSFTPAsyncReq) (*[]api.SSHResultRes, error) {
	if err := s.CheckSFTPParam(param); err != nil {
		return nil, err
	}

	insClientGroup := clientGroup{
		clientMap:      make(map[string]*gossh.Client),
		clientMapMutex: sync.Mutex{},
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
		sftpParam := &api.SFTPClientConfigReq{
			HostIp:      param.HostIp[i],
			Username:    param.Username[i],
			SSHPort:     param.SSHPort[i],
			Password:    password,
			Key:         param.Key,
			Passphrase:  param.Passphrase,
			Path:        param.Path[0],
			FileContent: param.FileContent[0],
		}
		// 如果是一个主机一个命令则多条
		if len(param.FileContent) > 1 {
			sftpParam.FileContent = param.FileContent[i]
		}
		if len(param.Path) > 1 {
			sftpParam.Path = param.Path[i]
		}
		go s.RunSFTPTransfer(sftpParam, channel, &wg, &insClientGroup)
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
func (s *SSHServer) CheckSFTPParam(param *api.RunSFTPAsyncReq) error {
	fmt.Println(param.HostIp)
	fmt.Println(param.Username)
	fmt.Println(param.SSHPort)
	fmt.Println(param.Path)
	fmt.Println(param.FileContent)
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
	for _, value := range param.Path {
		if value == "" {
			return fmt.Errorf("PATH切片中有值为空字符串")
		}
	}
	for _, value := range param.FileContent {
		if value == "" {
			return fmt.Errorf("FileContent切片中有值为空字符串")
		}
	}
	if len(param.HostIp) != len(param.SSHPort) || len(param.HostIp) != len(param.Username) || len(param.FileContent) > 1 && len(param.FileContent) != len(param.HostIp) || len(param.Path) > 1 && len(param.Path) != len(param.HostIp) {
		return errors.New("请检查: IP、端口、用户名、路径、文件内容这些切片是否一一对应")
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
