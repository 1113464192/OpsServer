package ops

import (
	"context"
	"fmt"
	"fqhWeb/internal/service/globalFunc"
	"fqhWeb/pkg/api/ops"
	"fqhWeb/pkg/logger"
	"os"
	"os/exec"
	"sync"
)

// 并发执行本地shell
func AsyncRunLocalShell(param *[]ops.RunLocalShellReq) (*[]ops.RunLocalShellRes, error) {
	channel := make(chan *ops.RunLocalShellRes, len(*param))
	wg := sync.WaitGroup{}
	var (
		result []ops.RunLocalShellRes
		err    error
	)
	for i := 0; i < len(*param); i++ {
		if err = globalFunc.Sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("获取信号失败，错误为: %v", err)
		}
		wg.Add(1)
		go RunLocalShell(&(*param)[i], channel, &wg)
	}
	wg.Wait()
	close(channel)
	for res := range channel {
		result = append(result, *res)
	}
	return &result, err
}

// 执行本地shell
func RunLocalShell(param *ops.RunLocalShellReq, ch chan *ops.RunLocalShellRes, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log().Error("Groutine", "RunLocalShell执行失败", r)
			fmt.Println("Groutine", "\n", "RunLocalShell执行", "\n", r)
			result := &ops.RunLocalShellRes{
				Mark:     param.Mark,
				Status:   9999,
				Response: fmt.Sprintf("触发了recover(): %v", r),
			}
			ch <- result
			wg.Done()
			globalFunc.Sem.Release(1)
		}
	}()
	result := &ops.RunLocalShellRes{
		Mark: param.Mark,
	}
	cmd := exec.Command("bash", "-c", param.CmdStr) // 连贯可以用"bash", "-c", "command"
	cmd.Env = append(os.Environ(), param.Env...)
	output, err := cmd.CombinedOutput()
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		result.Status = cmd.ProcessState.ExitCode()
		result.Response = fmt.Sprintf("错误返回: %s\n%s", string(output), err.Error())
		ch <- result
		wg.Done()
		globalFunc.Sem.Release(1)
		return
	}
	result.Response = string(output)
	result.Status = cmd.ProcessState.ExitCode()
	ch <- result
	wg.Done()
	globalFunc.Sem.Release(1)
}
