package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/service/globalFunc"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util/auth"
	"net/http"
	"strings"
	"sync"
)

type CSOperationService struct {
}

var (
	insCSOper = &CSOperationService{}
)

func CSOper() *CSOperationService {
	return insCSOper
}

func (s *CSOperationService) RunCSCmdAsync(param *[]api.SSHExecReq, args *string) (*[]api.CSCmdRes, error) {
	if err := s.CheckCSCMDParam(param); err != nil {
		return nil, err
	}

	channel := make(chan *api.CSCmdRes, len(*param))
	wg := sync.WaitGroup{}
	var err error
	var result []api.CSCmdRes
	// data := make(map[string]string)
	for i := 0; i < len(*param); i++ {
		if err = globalFunc.Sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("获取信号失败，错误为: %v", err)
		}
		wg.Add(1)
		go s.RunCSCmd(&(*param)[i], channel, &wg)
	}
	wg.Wait()
	close(channel)
	var argsMap map[string][]string
	if err = json.Unmarshal([]byte(*args), &argsMap); err != nil {
		return nil, fmt.Errorf("参数解析失败: %v", err)
	}
	// 返回结果格式化
	for res := range channel {
		for _, path := range argsMap["path"] {
			if strings.Contains(res.Response, path) {
				res.ServerDir = path
			}
		}
		result = append(result, *res)
	}

	return &result, err
}

func (s *CSOperationService) RunCSCmd(param *api.SSHExecReq, ch chan *api.CSCmdRes, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log().Error("Groutine", "RunCSCmd执行失败", r)
			fmt.Println("Groutine", "\n", "RunCSCmd执行失败", "\n", r)
			result := &api.CSCmdRes{
				HostIp:   param.HostIp,
				Status:   9999,
				Response: fmt.Sprintf("触发了recover(): %v", r),
			}
			ch <- result
			wg.Done()
			globalFunc.Sem.Release(1)
		}
	}()
	sslStr := "http"
	if configs.Conf.ClientSide.IsSSL == "true" {
		sslStr = "https"
	}
	url := fmt.Sprintf("%s://%s:%s/%s", sslStr, param.HostIp, configs.Conf.ClientSide.Port, consts.ClientExecApiPath)
	result, err := s.SendPostExecSignal(url, param.Cmd, param.HostIp)
	if err != nil {
		result.Status = 9999
		result.Response = fmt.Sprintf("请求失败: %v", err)
	}
	ch <- &result
	wg.Done()
	globalFunc.Sem.Release(1)
}

// 向Client服务器发送执行信号
// ClientExecApiPath
func (s *CSOperationService) SendPostExecSignal(url string, cmd string, ip string) (res api.CSCmdRes, err error) {
	// Create a new request
	cmdString := map[string]string{
		"string": cmd,
	}
	byteData, err := json.Marshal(cmdString)
	if err != nil {
		return res, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteData))
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add the ClientAuth header
	sign, err := auth.Md5EncryptSign(ip, configs.Conf.SecurityVars.ClientReqMd5Key)
	if err != nil {
		return res, err
	}
	req.Header.Set("ClientAuth", sign)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return res, errors.New("请求失败, 状态码: " + resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, err
	}
	res.HostIp = ip
	return res, err
}

// 检查是否符合执行条件
func (s *CSOperationService) CheckCSCMDParam(param *[]api.SSHExecReq) error {
	for _, p := range *param {
		if p.HostIp == "" || p.Cmd == "" {
			return fmt.Errorf("HostIp or CMD为空", p)
		}
	}
	return nil
}
