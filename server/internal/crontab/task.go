package crontab

import (
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils"
	"os"
	"time"
)

func CronMysqlLogRename() {
	now := time.Now().Local()
	previousDay := now.AddDate(0, 0, -1) // 获取前一天的日期
	logFileName := fmt.Sprintf(utils.GetRootPath()+"/logs/mysql/%s.log", previousDay.Format("20060102"))

	model.LogFile.Close() // 关闭之前的日志文件句柄

	// 重命名日志文件
	if err := os.Rename(utils.GetRootPath()+"/logs/mysql/mysql.log", logFileName); err != nil {
		logger.Log().Error("Mysql", "重命名mysql日志失败", err)
		return
	}

}

type RunSSHCmdAsyncReq struct {
	HostIp     []string `json:"host_ip"`
	Username   []string `json:"username"`
	SSHPort    []string `json:"ssh_port"`
	Password   []string `json:"password"`
	Key        []byte   `json:"key"`
	Passphrase []byte   `json:"passphrase"`
}

func CronWrittenHostInfo() {
	// 设定指定的用户，一般设置为高权限用户的私钥来执行全机器数据采集，这里设置为1
	var opsUser model.User
	if err := model.DB.First(&opsUser, consts.SSHOpsUserId).Error; err != nil {
		logger.Log().Error("User", "机器数据采集——获取OPS用户权限失败", err)
	}
	key := utils.XorDecrypt(opsUser.PriKey, consts.XorKey)
	passPhrase := utils.XorDecrypt(opsUser.KeyPasswd, consts.XorKey)
	var hosts []model.Host
	if err := model.DB.Find(&hosts).Error; err != nil {
		logger.Log().Error("Host", "机器数据采集——获取主机对象失败", err)
	}

	for _, host := range hosts {
		ipVar := host.Ipv4
		if !ipVar.Valid {
			ipVar = host.Ipv6
		}
		ip := ipVar.String
		username := host.User
		sshPort := host.Port

		// 获取服务器信息
		// 这里你可以根据自己的代码来处理服务器信息
		// 例如打印、存储到变量、调用其他函数等
		fmt.Printf("Server Name: %s\n", host.Name)
		fmt.Printf("Server IP: %s\n", ip)
		fmt.Println("----------------------")
	}
}
