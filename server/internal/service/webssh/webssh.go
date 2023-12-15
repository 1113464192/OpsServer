package webssh

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"fqhWeb/configs"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/ops"
	"fqhWeb/pkg/api"
	opsApi "fqhWeb/pkg/api/ops"
	"fqhWeb/pkg/logger"
	utilssh "fqhWeb/pkg/util/ssh"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type WebSshService struct {
}

var (
	insSSH = &WebSshService{}
)

func WebSsh() *WebSshService {
	return insSSH
}

func (s *WebSshService) WebSshHandle(c *gin.Context, sshParam *api.SSHExecReq) (wsRes string, err error) {
	var (
		conn    *websocket.Conn
		client  *ssh.Client
		session *ssh.Session
		sshConn *api.SSHConnect
		pid     int
	)

	if conn, err = s.upgraderWebSocket(c); err != nil {
		return fmt.Sprintf("用户名: %s, 机器IP: %s", sshParam.Username, sshParam.HostIp), fmt.Errorf("websocket连接失败: %v", err)
	}
	defer conn.Close()

	//Create ssh client
	var uid uint
	if err = model.DB.Model(&model.User{}).Select("id").Where("username = ?", sshParam.Username).First(&uid).Error; err != nil {
		s.WebSshSendText(conn, []byte("查询数据库获取用户ID时发生错误: "+err.Error()))
		return fmt.Sprintf("用户名: %s, 机器IP: %s", sshParam.Username, sshParam.HostIp), fmt.Errorf("查询数据库获取用户ID时发生错误: %v", err)
	}
	// 将uid_HostIP组成一个字符串赋值给wsRes
	wsRes = fmt.Sprintf("%d_%s", uid, sshParam.HostIp)

	sockPath := fmt.Sprintf("/tmp/%s_agent.%d", sshParam.Username, uid)

	// 生成ssh agent sock
	if pid, err = s.generateLocalSSHAgentSocket(sockPath, uid, wsRes, sshParam.Key, sshParam.Passphrase); err != nil {
		s.WebSshSendText(conn, []byte("生成ssh agent socket时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("生成ssh agent socket时发生错误: %v", err)
	}

	// 生成sshClient
	if client, err = utilssh.SSHNewClient(sshParam.HostIp, sshParam.Username, sshParam.SSHPort, sshParam.Password, sshParam.Key, sshParam.Passphrase, sockPath); err != nil {
		s.WebSshSendText(conn, []byte("生成ssh.Client时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("生成ssh.Client时发生错误: %v", err)
	}
	defer client.Close()

	// 生成sshSession
	if session, err = utilssh.SSHNewSession(client); err != nil {
		s.WebSshSendText(conn, []byte("生成ssh.Session时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("生成ssh.Session时发生错误: %v", err)
	}
	defer session.Close()

	// 生成sshConn
	if sshConn, err = utilssh.SSHNewConnect(session); err != nil {
		s.WebSshSendText(conn, []byte("创建ssh连接时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("创建ssh连接时发生错误: %v", err)
	}

	quit := make(chan int)
	go s.Output(conn, sshConn, quit)
	go s.Recv(conn, sshConn, quit)
	<-quit

	// 清除用户ssh agent socket与ssh agent process
	if err = s.removeSSHAgentSocket(sockPath, pid); err != nil {
		return wsRes, fmt.Errorf("清除用户ssh agent socket与ssh agent process时发生错误: %v", err)
	}
	return wsRes, nil
}

func (s *WebSshService) generateLocalSSHAgentSocket(sockPath string, uid uint, mark string, priKey []byte, passphrase []byte) (pid int, err error) {
	var (
		idKeyFile       *os.File
		localShellParam []opsApi.RunLocalShellReq
		result          *[]opsApi.RunLocalShellRes
	)
	// Aes解密SSH密钥和passphrase
	if err = utilssh.DecryptAesSSHKey(&priKey, &passphrase); err != nil {
		return 0, fmt.Errorf("用户私钥解密失败: %v", err)
	}

	// 生成SSH Agent Socket
	id_key_path := fmt.Sprintf("/tmp/%d_key", uid)

	// 写入私钥到机器中
	idKeyFile, err = os.OpenFile(id_key_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400)
	defer idKeyFile.Close()
	if _, err = idKeyFile.Write(priKey); err != nil {
		return 0, fmt.Errorf("用户私钥解密失败: %v", err)
	}

	env := []string{
		fmt.Sprintf("agent_sock_path=%s", sockPath),
		fmt.Sprintf("id_key_path=%s", id_key_path),
		fmt.Sprintf("id_key_passphrase=%s", string(passphrase)),
	}
	cmdStr := fmt.Sprintf("cd %s/server/shellScript && ./ssh_agent.sh", configs.Conf.ProjectWeb.RootPath)
	localShellParam = append(localShellParam, opsApi.RunLocalShellReq{
		CmdStr: cmdStr,
		Env:    env,
		Mark:   mark,
	})
	if result, err = ops.AsyncRunLocalShell(&localShellParam); err != nil || (*result)[0].Status != 0 {
		return 0, fmt.Errorf("ssh agent socket 生成失败: %v", err)
	}
	// 属于安全措施，定时任务也会每分钟检测
	if err = os.Remove(id_key_path); err != nil {
		return 0, fmt.Errorf("删除id_key_path文件失败: %v\n很严重的权限问题, 请立即通知相关运维手动删除", err)
	}
	if strings.Contains((*result)[0].Response, "Success") {
		pidStr := strings.Split((*result)[0].Response, "\n")[2]
		if pid, err = strconv.Atoi(pidStr); err != nil {
			return 0, fmt.Errorf("pid转换 uint64 失败: %v", err)
		}
		return pid, nil
	} else {
		return 0, fmt.Errorf("ssh agent socket 生成失败, Response不包含Success: %v", (*result)[0].Response)
	}
}

func (s *WebSshService) upgraderWebSocket(c *gin.Context) (conn *websocket.Conn, err error) {
	// 获取超时时间
	var duration time.Duration
	duration, err = time.ParseDuration(configs.Conf.Webssh.HandshakeTimeout)
	if err != nil {
		return nil, fmt.Errorf("超时时间获取失败: %v", err)
	}

	var upgrader = websocket.Upgrader{
		HandshakeTimeout: duration,
		// 读写缓冲大小, 这个值越大，一次可以处理的数据就越多，但是也会消耗更多的内存
		// 如果不设置的话它们的值默认是 4096 byte
		ReadBufferSize:  configs.Conf.Webssh.ReadBufferSize,
		WriteBufferSize: configs.Conf.Webssh.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			// 写入状态码
			w.WriteHeader(status)
			// 写入错误信息
			w.Write([]byte("WebSocket upgrade failed: " + reason.Error()))
		},
	}

	if conn, err = upgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		return nil, fmt.Errorf("websocket连接失败: %v", err)
	}
	return conn, nil
}

func (s *WebSshService) WebSshSendText(conn *websocket.Conn, b []byte) error {
	if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
		return err
	}
	return nil
}

func (s *WebSshService) removeSSHAgentSocket(sockPath string, pid int) (err error) {
	if err = os.Remove(sockPath); err != nil {
		logger.Log().Error("RemoveSSHAgentSocket", "ssh_agent.sock文件删除错误", err)
		// 接入微信小程序之类的请求, 向运维发送处理ssh_agent.sock文件问题
		fmt.Println("微信小程序=====向运维发送,处理ssh_agent.sock文件问题")
		return fmt.Errorf("删除ssh_agent.sock文件失败: %v\n很严重的权限问题, 请立即通知相关运维手动删除", err)
	}
	// 获取进程
	var process *os.Process
	process, err = os.FindProcess(pid)
	if err != nil {
		logger.Log().Error("RemoveSSHAgentSocket", "ssh_agent进程查询错误", err)
		// 接入微信小程序之类的请求, 向运维发送处理ssh_agent.sock文件问题
		fmt.Println("微信小程序=====向运维发送,处理ssh_agent进程问题")
		return fmt.Errorf("查询ssh_agent进程失败: %v\n很严重的权限问题, 请立即通知相关运维手动删除", err)
	}

	// 关闭进程
	if err = process.Signal(syscall.SIGKILL); err != nil {
		logger.Log().Error("RemoveSSHAgentSocket", "关闭ssh_agent进程错误", err)
		// 接入微信小程序之类的请求, 向运维发送处理ssh_agent.sock文件问题
		fmt.Println("微信小程序=====向运维发送,处理ssh_agent进程问题")
		return fmt.Errorf("关闭ssh_agent进程失败: %v\n很严重的权限问题, 请立即通知相关运维手动删除", err)
	}
	return nil
}
