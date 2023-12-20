package globalFunc

import (
	"fqhWeb/configs"

	"golang.org/x/sync/semaphore"
)

// semaphore
var (
	Sem           *semaphore.Weighted
	MaxWebSSH     uint64
	WebSSHcounter uint64
)

func DeclareGlobalVar() {
	// 设置总并发数
	Sem = semaphore.NewWeighted(configs.Conf.Concurrency.Number)
	MaxWebSSH = uint64(configs.Conf.Webssh.MaxConnNumber)
	WebSSHcounter = 0
}
