package globalfunc

import (
	"errors"
	"fqhWeb/configs"
	"sync"
)

var (
	maxWebSSH        = configs.Conf.Webssh.MaxConnNumber
	counter   uint64 = 0
	mu        sync.Mutex
)

func IncreaseWebSSHConn() error {
	mu.Lock()
	defer mu.Unlock()

	if counter >= maxWebSSH {
		return errors.New("已达到最大webssh数量")
	}

	counter++
	// 创建webssh的代码在这里

	return nil
}

func ReduceWebSSHConn() {
	mu.Lock()
	defer mu.Unlock()

	if counter > 0 {
		counter--
	}
}
