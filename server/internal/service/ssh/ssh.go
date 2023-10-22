package ssh

import (
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/pkg/logger"
	"time"

	"golang.org/x/crypto/ssh"
)

type Service struct {
}

var (
	insSSH = &Service{}
)

func SSH() *Service {
	return insSSH
}

type ClientConfigService struct {
	Host      string // ip
	Port      int64  // 端口
	Username  string // 用户名
	Password  string // 密码，填密码优先走密码，走公私钥不用传
	Key       []byte // 私钥字符串
	KeyPasswd []byte // 私钥密码(有就需要输入，没有不用传)
}

func (s *Service) RunCmd(params *ClientConfigService, cmd string) {

}

// 创建*ssh.Client
func (clientC *ClientConfigService) createClient() (client *ssh.Client, err error) {
	if len(clientC.Key) == 0 && clientC.Password == "" {
		logger.Log().Error("SSH", "密钥和密码都为空", err)
		return nil, errors.New("密钥和密码都为空")
	}
	config := &ssh.ClientConfig{
		User: clientC.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(clientC.Password),
		},
		Timeout: 8 * time.Second,
	}
	if len(clientC.Key) != 0 && clientC.Password == "" {
		var signer ssh.Signer
		if len(clientC.KeyPasswd) != 0 {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(clientC.Key, clientC.KeyPasswd)
		} else {
			signer, err = ssh.ParsePrivateKey(clientC.Key)
		}

		if err != nil {
			logger.Log().Error("SSH", "密钥解析失败", err)
			return nil, errors.New("密钥解析失败")
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}

	// InsecureIgnoreHostKey() 方法来忽略主机密钥验证，这在测试或开发环境中可以接受，但在生产环境中应该谨慎使用。建议使用 ssh.FixedHostKey() 方法来验证主机密钥。
	if configs.Conf.System.Mode != "product" {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", clientC.Host, clientC.Port), config)
	if err != nil {
		logger.Log().Error("SSH", "创建ssh的client失败", err)
		return nil, errors.New("创建ssh的client失败")
	}

	return client, nil

}
