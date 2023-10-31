package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api"
	"strings"
)

type SshService struct {
}

var (
	insSSH = &SshService{}
)

func SSH() *SshService {
	return insSSH
}

func (s *SshService) GetFlag(task *model.TaskTemplate) (err error) {
	var args []map[string]string
	if err = json.Unmarshal([]byte(task.Args), &args); err != nil {
		return errors.New("参数字段进行json解析失败")
	}

}

func (s *SshService) filterPortRuleHost(host *[]model.Host, task *model.TaskTemplate) (err error) {
	var portRule []map[int]string
	if err = json.Unmarshal([]byte(task.PortRule), &portRule); err != nil {
		return errors.New("端口规则进行json解析失败")
	}

}

func (s *SshService) filterConditionHost(host *[]model.Host, task *model.TaskTemplate) (err error) {
	var condition []map[string]string
	if err = json.Unmarshal([]byte(task.Condition), &condition); err != nil {
		return errors.New("筛选机器条件规则进行json解析失败")
	}
	var fields []string
	// 为了使用不定长参数的解包方法，所以要设置为interface{}
	var values []any
	for _, c := range condition {
		for key, value := range c {
			switch key {
			case "mem":
				fields = append(fields, "curr_mem > ?")
				values = append(values, value)
			case "data_disk":
				fields = append(fields, "curr_data_disk > ?")
				values = append(values, value)
			case "iowait":
				fields = append(fields, "curr_iowait < ?")
				values = append(values, value)
			case "idle":
				fields = append(fields, "curr_idle > ?")
				values = append(values, value)
			case "load":
				fields = append(fields, "curr_load < ?")
				values = append(values, value)
			default:
				return fmt.Errorf("%s 不属于ConditionSet中的任何一个", key)
			}
		}
	}
	// 使用单个查询筛选符合条件的主机
	if len(fields) > 0 {
		conditions := strings.Join(fields, " AND ")
		if err = model.DB.Where(conditions, values...).Find(host).Error; err != nil {
			return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
		}
	}
	return err
	// for _, c := range condition {
	// 	for key, value := range c {
	// 		switch key {
	// 		case "mem":
	// 			if err = model.DB.Where("curr_mem > ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "data_disk":
	// 			if err = model.DB.Where("curr_data_disk > ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "iowait":
	// 			if err = model.DB.Where("curr_iowait < ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "idle":
	// 			if err = model.DB.Where("curr_idle > ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		case "load":
	// 			if err = model.DB.Where("curr_load < ?", value).Find(host).Error; err != nil {
	// 				return fmt.Errorf("%s%v", "筛选符合条件的主机操作报错: ", err)
	// 			}
	// 		default:
	// 			return fmt.Errorf("%s 不属于ConditionSet中的其中一个", key)
	// 		}
	// 	}
	// }
	// return err
}

func (s *SshService) GetTemplateParam(param api.GetTemplateParamReq) (resParam *api.RunCmdtemRes, err error) {
	var task model.TaskTemplate
	var user model.User
	if err = model.DB.First(&task, param.Tid).Error; err != nil {
		return nil, errors.New("根据id查询任务失败")
	}
	if err = model.DB.First(&user, param.Uid).Error; err != nil {
		return nil, errors.New("根据id查询用户失败")
	}
	var host []model.Host
	if task.Condition != "" {
		if err = s.filterConditionHost(&host, &task); err != nil {
			return nil, errors.New("筛选符合条件的主机失败")
		}
	}
	if task.PortRule != "" {
		if task.Args == "" {
			return nil, errors.New("有端口规则请传flag或path, 否则无标识判断")
		}
		if err = s.filterPortRuleHost(&host, &task); err != nil {
			return nil, errors.New("筛选符合端口空余的主机失败")
		}

	}
	resParam.Cmd = task.Task
	resParam.Key = []byte(user.PriKey.String)
	resParam.KeyPasswd = []byte(user.KeyPasswd.String)
}

// func (s *Service) RunCmd(params *ClientConfigService, cmd string) {

// }

// // 创建*ssh.Client
// func (clientC *ClientConfigService) createClient() (client *ssh.Client, err error) {
// 	if len(clientC.Key) == 0 && clientC.Password == "" {
// 		logger.Log().Error("SSH", "密钥和密码都为空", err)
// 		return nil, errors.New("密钥和密码都为空")
// 	}
// 	config := &ssh.ClientConfig{
// 		User: clientC.Username,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(clientC.Password),
// 		},
// 		Timeout: 8 * time.Second,
// 	}
// 	if len(clientC.Key) != 0 && clientC.Password == "" {
// 		var signer ssh.Signer
// 		if len(clientC.KeyPasswd) != 0 {
// 			signer, err = ssh.ParsePrivateKeyWithPassphrase(clientC.Key, clientC.KeyPasswd)
// 		} else {
// 			signer, err = ssh.ParsePrivateKey(clientC.Key)
// 		}

// 		if err != nil {
// 			logger.Log().Error("SSH", "密钥解析失败", err)
// 			return nil, errors.New("密钥解析失败")
// 		}
// 		config.Auth = []ssh.AuthMethod{
// 			ssh.PublicKeys(signer),
// 		}
// 	}

// 	// InsecureIgnoreHostKey() 方法来忽略主机密钥验证，这在测试或开发环境中可以接受，但在生产环境中应该谨慎使用。建议使用 ssh.FixedHostKey() 方法来验证主机密钥。
// 	if configs.Conf.System.Mode != "product" {
// 		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
// 	}

// 	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", clientC.Host, clientC.Port), config)
// 	if err != nil {
// 		logger.Log().Error("SSH", "创建ssh的client失败", err)
// 		return nil, errors.New("创建ssh的client失败")
// 	}

// 	return client, nil

// }