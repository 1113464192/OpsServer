package ssh

import (
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/consts"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/util"
	"net"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func AuthWithPrivateKeyBytes(key []byte, passphrase []byte) (ssh.AuthMethod, error) {
	var signer ssh.Signer
	var err error
	if passphrase == nil {
		signer, err = ssh.ParsePrivateKey(key)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
	}
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func AuthWithAgent(sockPath string) (ssh.AuthMethod, error) {
	socks, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, err
	}
	defer socks.Close()
	// 1. 返回Signers函数的结果
	sshAgent := agent.NewClient(socks)
	return ssh.PublicKeysCallback(sshAgent.Signers), nil
}

// func SSHNewClient(config *api.SSHExecReq) (client *ssh.Client, err error) {
func SSHNewClient(hostIp string, username string, sshPort string, password string, priKey []byte, passphrase []byte, sockPath string) (client *ssh.Client, err error) {
	duration, err := time.ParseDuration(configs.Conf.SshConfig.SshClientTimeout)
	if err != nil {
		return nil, fmt.Errorf("超时时间获取失败: %v", err)
	}

	clientConfig := &ssh.ClientConfig{
		User:            username,
		Timeout:         duration,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略public key的安全验证
	}

	if sshPort == "" {
		sshPort = consts.SSHDefaultPort
	}

	// 1. private key bytes
	// AES解密私钥
	if err = DecryptAesSSHKey(&priKey, &passphrase); err != nil {
		return nil, fmt.Errorf("用户私钥解密失败: %v", err)
	}

	if priKey != nil {
		if auth, err := AuthWithPrivateKeyBytes(priKey, passphrase); err == nil {
			clientConfig.Auth = append(clientConfig.Auth, auth)
		}
	}
	// 2. 密码方式 放在key之后,这样密钥失败之后可以使用Password方式
	if password != "" {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(password))
	}

	// 3. agent 模式放在最后,意味着websocket连接，需要使用 openssh agent forwarding
	if sockPath != "" {
		if auth, err := AuthWithAgent(sockPath); err == nil {
			clientConfig.Auth = append(clientConfig.Auth, auth)
		}
	}

	if clientConfig.Auth == nil {
		return nil, errors.New("未能生成clientConfig.Auth")
	}
	client, err = ssh.Dial("tcp", net.JoinHostPort(hostIp, sshPort), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("生成ssh.Client失败: %v", err)
	}
	return client, err
}

func SSHNewSession(client *ssh.Client) (session *ssh.Session, err error) {
	session, err = client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("生成ssh.Session失败: %v", err)
	}
	return session, err
}

func SSHNewConnect(session *ssh.Session) (*api.SSHConnect, error) {
	modes := ssh.TerminalModes{
		// 设置控制是否在终端上显示输入。设置为0表示不显示输入，这在输入密码等敏感信息时很有用。
		ssh.ECHO: configs.Conf.Webssh.SshEcho,
		// 设置控制终端的输入输出的速度。设置为14400表示输入速度为14400bps。
		// 如果你设置的值超过了你的网络带宽，可能会导致数据传输不稳定，可能会出现丢包、延迟增大等问题
		ssh.TTY_OP_ISPEED: configs.Conf.Webssh.SshTtyOpIspeed,
		ssh.TTY_OP_OSPEED: configs.Conf.Webssh.SshTtyOpOspeed,
	}
	// 伪终端屏幕高、宽大小(单位为字符)
	if err := session.RequestPty(consts.WebsshTerminal, configs.Conf.Webssh.ScreenHight, configs.Conf.Webssh.ScreenWeight, modes); err != nil {
		return nil, fmt.Errorf("生成Session.RequestPty失败: %v", err)
	}
	pipe, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("生成session.StdinPipe失败: %v", err)
	}
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("生成session.StdoutPipe失败: %v", err)
	}

	return &api.SSHConnect{
		Session:    session,
		StdinPipe:  pipe,
		StdoutPipe: stdoutPipe,
	}, nil
}

func CreateSFTPClient(client *ssh.Client) (sftpClient *sftp.Client, err error) {
	sftpClient, err = sftp.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("生成sftp.Client失败: %v", err)
	}

	return sftpClient, err
}

func DecryptAesSSHKey(key *[]byte, passphrase *[]byte) (err error) {
	// AES解密私钥
	*key, err = util.DecryptAESCBC(*key, []byte(configs.Conf.SecurityVars.AesKey), []byte(configs.Conf.SecurityVars.AesIv))
	if err != nil {
		return fmt.Errorf("用户私钥解密失败: %v", err)
	}
	// AES解密passphrase
	if *passphrase != nil {
		*passphrase, err = util.DecryptAESCBC(*passphrase, []byte(configs.Conf.SecurityVars.AesKey), []byte(configs.Conf.SecurityVars.AesIv))
		if err != nil {
			return fmt.Errorf("用户passphrase解密失败: %v", err)
		}
	}
	return nil
}
