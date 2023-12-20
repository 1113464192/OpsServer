package webssh

import (
	"fmt"
	"net"
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
	"golang.org/x/crypto/ssh/agent"
)

type WebSshService struct {
}

var (
	insSSH = &WebSshService{}
)

func WebSsh() *WebSshService {
	return insSSH
}

func (s *WebSshService) WebSshHandle(c *gin.Context, user *model.User, param api.WebsshConnReq) (wsRes string, err error) {
	var (
		host            model.Host
		conn            *websocket.Conn
		client          *ssh.Client
		session         *ssh.Session
		sshConn         *SSHConnect
		sshAgentPointer *agent.ExtendedAgent
		pid             int
	)

	if err = model.DB.First(&host, param.Hid).Error; err != nil {
		return "", fmt.Errorf("服务器 %d 查询失败: %v", param.Hid, err)
	}
	sshParam := &api.SSHExecReq{
		HostIp:     host.Ipv4.String,
		SSHPort:    host.Port,
		Username:   host.User,
		Password:   string(host.Password),
		Key:        user.PriKey,
		Passphrase: user.Passphrase,
	}

	if conn, err = s.upgraderWebSocket(c); err != nil {
		return fmt.Sprintf("用户名: %s, 机器IP: %s", sshParam.Username, sshParam.HostIp), fmt.Errorf("websocket连接失败: %v", err)
	}
	defer conn.Close()

	// Aes解密SSH密钥和passphrase
	if err = utilssh.DecryptAesSSHKey(&sshParam.Key, &sshParam.Passphrase); err != nil {
		return "", fmt.Errorf("用户私钥解密失败: %v", err)
	}

	// Create ssh client
	// 将uid_HostIP组成一个字符串赋值给wsRes
	wsRes = fmt.Sprintf("%d_%s", user.ID, sshParam.HostIp)

	sockPath := fmt.Sprintf("/tmp/agent.%d", user.ID)

	// 生成ssh agent sock
	if pid, err = s.generateLocalSSHAgentSocket(sockPath, user.ID, wsRes, sshParam.Key, sshParam.Passphrase); err != nil {
		s.WebSshSendText(conn, []byte("生成ssh agent socket时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("生成ssh agent socket时发生错误: %v", err)
	}

	// 生成sshClient
	var netConn net.Conn
	if client, netConn, sshAgentPointer, err = utilssh.SSHNewClient(sshParam.HostIp, sshParam.Username, sshParam.SSHPort, sshParam.Password, nil, nil, sockPath); err != nil {
		s.WebSshSendText(conn, []byte("生成ssh.Client时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("生成ssh.Client时发生错误: %v", err)
	}
	defer netConn.Close()
	defer client.Close()

	// 生成sshSession
	if session, err = utilssh.SSHNewSession(client); err != nil {
		s.WebSshSendText(conn, []byte("生成ssh.Session时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("生成ssh.Session时发生错误: %v", err)
	}
	defer session.Close()

	// 生成sshConn
	if sshConn, err = SSHNewConnect(client, session, sshAgentPointer, param.WindowSize); err != nil {
		s.WebSshSendText(conn, []byte("创建ssh连接时发生错误: "+err.Error()))
		return wsRes, fmt.Errorf("创建ssh连接时发生错误: %v", err)
	}

	quit := make(chan int)
	go s.Output(conn, sshConn, quit)
	go s.Recv(conn, sshConn, quit)
	<-quit

	// 清除用户ssh agent socket与ssh agent process
	if err = s.removeSSHAgentSocket(pid, sockPath); err != nil {
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

	// 生成SSH Agent Socket
	id_key_path := fmt.Sprintf("/tmp/%d_key", uid)

	// 写入私钥到机器中
	if idKeyFile, err = os.OpenFile(id_key_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400); err != nil {
		return 0, fmt.Errorf("用户私钥文件打开失败: %v", err)
	}
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
	if strings.Contains((*result)[0].Response, "Success") {
		resStr := strings.Split((*result)[0].Response, "\n")
		pidStr := resStr[len(resStr)-2]
		if pid, err = strconv.Atoi(pidStr); err != nil {
			return 0, fmt.Errorf("pid转换 int 失败: %v", err)
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

func (s *WebSshService) removeSSHAgentSocket(pid int, sockPath string) (err error) {
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

	// 删除socket文件
	if err = os.Remove(sockPath); err != nil {
		logger.Log().Error("RemoveSSHAgentSocket", "删除ssh_agent_socket文件错误", err)
		// 接入微信小程序之类的请求, 向运维发送处理ssh_agent.sock文件问题
		fmt.Println("微信小程序=====向运维发送,处理ssh_agent_socket文件问题")
		return fmt.Errorf("删除ssh_agent_socket文件失败: %v\n请通知相关运维手动删除", err)
	}
	return nil
}
