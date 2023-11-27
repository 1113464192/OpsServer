package ssh

import (
	"errors"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/pkg/util"
	"net"
	"os"
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

func AuthWithAgent() (ssh.AuthMethod, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, errors.New("agent disabled")
	}
	socks, err := net.Dial("unix", sock)
	if err != nil {
		return nil, err
	}
	// 1. 返回Signers函数的结果
	client := agent.NewClient(socks)
	signers, err := client.Signers()
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signers...), nil
}

// func SSHNewClient(config *api.SSHClientConfigReq) (client *ssh.Client, err error) {
func SSHNewClient(hostIp string, username string, sshPort string, password string, priKey []byte, passphrase []byte) (client *ssh.Client, err error) {
	clientConfig := &ssh.ClientConfig{
		User:            username,
		Timeout:         consts.SSHTimeout * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略public key的安全验证
	}

	if sshPort == "" {
		sshPort = "22"
	}

	// 1. private key bytes
	key := util.XorDecrypt(priKey, consts.XorKey)
	passPhrase := util.XorDecrypt(passphrase, consts.XorKey)
	if priKey != nil {
		if auth, err := AuthWithPrivateKeyBytes(key, passPhrase); err == nil {
			clientConfig.Auth = append(clientConfig.Auth, auth)
		}
	}
	// 2. 密码方式 放在key之后,这样密钥失败之后可以使用Password方式
	if password != "" {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(password))
	}
	// 3. agent 模式放在最后,这样当前两者都不能使用时可以采用Agent模式
	if auth, err := AuthWithAgent(); err == nil {
		clientConfig.Auth = append(clientConfig.Auth, auth)
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

func CreateSFTPClient(client *ssh.Client) (sftpClient *sftp.Client, err error) {
	sftpClient, err = sftp.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("生成sftp.Client失败: %v", err)
	}

	return sftpClient, err
}
