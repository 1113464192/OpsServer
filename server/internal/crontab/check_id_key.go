package crontab

import (
	"fmt"
	"fqhWeb/internal/service/ops"
	opsApi "fqhWeb/pkg/api/ops"
	"fqhWeb/pkg/logger"
	"time"
)

func CronExistCheckIdKey() {
	var (
		err             error
		localShellParam []opsApi.RunLocalShellReq
		result          *[]opsApi.RunLocalShellRes
	)
	localShellParam = append(localShellParam, opsApi.RunLocalShellReq{
		CmdStr: `ls /tmp/ | grep -E "[0-9]+_key"`,
		Mark:   "check_id_key",
	})
	// 正常返回代表有id_key文件，需要观察30s看看是不是程序刚生成的，如果是，就不用管，如果不是，就通知运维
	if result, err = ops.AsyncRunLocalShell(&localShellParam); err != nil || (*result)[0].Status == 0 {
		time.Sleep(30 * time.Second)
		// 再次检测，如果还有id_key文件，就通知运维
		if result, err = ops.AsyncRunLocalShell(&localShellParam); err != nil || (*result)[0].Status == 0 {
			logger.Log().Error("CronCheckIdKey", "发现id_key文件,可能函数中删除错误", err)
			// 接入微信小程序之类的请求, 向运维发送处理id_key问题
			fmt.Println("微信小程序=====向运维发送,处理id_key问题")
			return
		}
	}
}
