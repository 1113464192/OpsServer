package globalFunc

import (
	"errors"
	"sync"
)

// webssh
var (
	mu sync.Mutex
)

func IncreaseWebSSHConn() error {
	mu.Lock()
	defer mu.Unlock()
	if WebSSHcounter >= MaxWebSSH {
		return errors.New("已达到最大webssh数量")
	}

	WebSSHcounter++

	return nil
}

func ReduceWebSSHConn() {
	mu.Lock()
	defer mu.Unlock()

	if WebSSHcounter > 0 {
		WebSSHcounter--
	}
}
