package consts

import (
	"time"
)

// webssh
const (
	WebsshLinuxTerminal = "linux"
	// 终端窗口
	WebsshXTerminal                   = "xterm"
	WebsshPingPeriod                  = 20 * time.Second
	WebsshPongWait                    = WebsshPingPeriod * 2
	WebsshWriteWait                   = 10 * time.Second
	WebsshReadMessageTickerDuration   = time.Millisecond * time.Duration(40)
	WebsshSockPath                    = `/tmp/agent.%d`
	WebsshIdKeyPath                   = `/tmp/%d_key`
	WebsshMaxRecordLength             = 2048
	WebsshGenerateLocalSSHAgentSocket = "cd %s/server/shellScript && ./ssh_agent.sh"
)
