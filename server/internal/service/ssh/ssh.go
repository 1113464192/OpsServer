package ssh

import (
	"context"
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/utils"
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

func (s *SSHServer) RunSSHCmdAsync(param api.RunSSHCmdAsyncReq, cmd string) (result *api.SSHResultRes, err error) {
	if len(param.HostIp) != len(param.SSHPort) || len(param.HostIp) != len(param.Username) || len(param.HostIp) != len(param.Password) {
		return nil, errors.New("请检查: IP、端口、用户名是否一一对应, 不使用密码情况下是否传递空字符串")
	}
	channel := make(chan *api.GetSSHRes, len(param.HostIp))
	wg := sync.WaitGroup{}
	// data := make(map[string]string)
	var data []string
	for i := 0; i < len(param.HostIp); i++ {
		if err = configs.Sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("获取信号失败，错误为: %v", err)
		}
		wg.Add(1)
		sshParam := &api.SSHClientConfigReq{
			HostIp:     param.HostIp[i],
			Username:   param.Username[i],
			SSHPort:    param.SSHPort[i],
			Password:   param.Password[i],
			Key:        param.Key,
			Passphrase: param.Passphrase,
		}

		go s.RunSSHCmd(sshParam, cmd, channel, &wg)

	}
	wg.Wait()
	close(channel)
	for res := range channel {
		value, ok := res.Response.Data.(string)
		if ok {
			// data[res.HostIps] = value
			data = append(data, res.HostIps+"="+value)
		} else {
			return nil, errors.New("channel结果提取返回转换为string时触发报错")
		}
	}
	dataJson, err := utils.ConvertToJsonPair(data)
	if err != nil {
		return nil, err
	}
	result = &api.SSHResultRes{
		HostIps: param.HostIp,
		Status:  true,
		Response: api.Response{
			Data: dataJson,
		},
	}
	return result, err
}

func (s *SSHServer) RunSSHCmd(param *api.SSHClientConfigReq, cmd string, ch chan *api.GetSSHRes, wg *sync.WaitGroup) error {
	result := &api.GetSSHRes{
		HostIps: param.HostIp,
	}

	client, err := ssh.SSHNewClient(param)
	if err != nil {
		result.Response.Data = err.Error()
		ch <- result
		return err
	}
	defer client.Close()
	session, err := ssh.SSHNewSession(client)
	if err != nil {
		result.Response.Data = err.Error()
		ch <- result
		return err
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		result.Response.Data = fmt.Sprintf("Failed to execute command: %s %s", string(output), err.Error())
		ch <- result
		return fmt.Errorf("failed to execute command: %s %v", string(output), err)
	}
	result.Response.Data = string(output)
	ch <- result
	wg.Done()
	configs.Sem.Release(1)
	return nil
}
