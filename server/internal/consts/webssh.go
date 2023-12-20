package consts

import "time"

// webssh
const (
	WebsshLinuxTerminal = "linux"
	// 终端窗口
	WebsshXTerminal           = "xterm"
	PingPeriod                = 1 * time.Minute
	PongPeriod                = 1 * time.Minute
	ReadMessageTickerDuration = time.Millisecond * time.Duration(40)
)
