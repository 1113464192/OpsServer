package consts

import "time"

const (
	SSHOpsUserId   = 1
	SSHDefaultPort = "22"
)

// webssh
const (
	WebsshLinuxTerminal = "linux"
	// 终端窗口
	WebsshXTerminal = "xterm"
	PingPeriod      = 1 * time.Minute
	PongPeriod      = 1 * time.Minute
)
