package gitWebhook

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/api/gitWebhook"
	"fqhWeb/pkg/util2"
	"strings"
)

func (s *GitWebhookService) ExecServerCustomCi(whId uint, sshurl string, name string, hid uint) (err error) {
	// 获取最高权限用户
	var opsUser model.User
	if err = model.DB.First(&opsUser, consts.SSHOpsUserId).Error; err != nil {
		return fmt.Errorf("获取ops权限用户失败: %v", err)
	}
	fmt.Println("============", sshurl, name, hid, configs.Conf.GitWebhook.GitCiScriptDir, configs.Conf.GitWebhook.GitCiRepo, "============")

	// 获取执行主机
	var host model.Host
	if err = model.DB.First(&host, hid).Error; err != nil {
		return fmt.Errorf("获取CI服务器失败: %v", err)
	}

	// 执行CI操作
	cmd := fmt.Sprintf(`bash %s/%s.sh %d %s %s`, configs.Conf.GitWebhook.GitCiScriptDir, name, whId, sshurl, configs.Conf.GitWebhook.GitCiRepo+"/"+name)

	var sshClientConfigParam []api.SSHClientConfigReq
	sshClientConfigParam = append(sshClientConfigParam,
		api.SSHClientConfigReq{
			HostIp:     host.Ipv4.String,
			Username:   host.User,
			SSHPort:    host.Port,
			Key:        opsUser.PriKey,
			Passphrase: opsUser.Passphrase,
			Cmd:        cmd,
		})
	var sshResult *[]api.SSHResultRes
	if sshResult, err = service.SSH().RunSSHCmdAsync(&sshClientConfigParam); err != nil && (*sshResult)[0].Status != 0 {
		return fmt.Errorf("执行CI脚本报错: %v \n %s", err, (*sshResult)[0].Response)
	}
	return err
}

func (s *GitWebhookService) UpdateStatusAuth(clientSign string, clientIp string) (err error) {
	if clientSign == "" {
		return errors.New("未从header获取签名到CiAuthSign的值")
	} else if clientIp == "" {
		return errors.New("未从header获取签名到clientIp的值")
	}
	sign, err := s.md5EncryptSign(clientIp, configs.Conf.SecurityVars.CiMd5Key)
	if err != nil {
		return fmt.Errorf("sign生成报错: %v", err)
	}
	if sign != clientSign {
		return errors.New(`认证码错误，请确认认证码生成方式`)
	}
	return err
}

func (s *GitWebhookService) md5EncryptSign(clientIp string, md5Key string) (sign string, err error) {
	builder := strings.Builder{}
	builder.WriteString(clientIp)
	builder.WriteString(md5Key)
	md5Hash := md5.Sum([]byte(builder.String()))
	sign = hex.EncodeToString(md5Hash[:])
	return sign, err
}

func (s *GitWebhookService) UpdateGitWebhookStatus(param gitWebhook.UpdateGitWebhookStatusReq) (err error) {
	if err = model.DB.Model(&model.GitWebhookRecord{}).Where("id = ?", param.Id).Update("status", param.Status).Error; err != nil {
		return fmt.Errorf("更改GitWebhook的Status失败: %v", err)
	}
	return err
}

func (s *GitWebhookService) UpdateGitWebhook(param gitWebhook.UpdateGitWebhookReq) (*model.GitWebhookRecord, error) {
	var gitWebhookRecord model.GitWebhookRecord
	var err error
	if !util2.CheckIdExists(&gitWebhookRecord, param.Id) {
		return nil, errors.New("记录不存在")
	}

	if err := model.DB.Where("id = ?", param.Id).First(&gitWebhookRecord).Error; err != nil {
		return &gitWebhookRecord, errors.New("GitWebhook记录数据库查询失败")
	}

	gitWebhookRecord.FullName = param.FullName
	gitWebhookRecord.ProjectId = param.ProjectId
	gitWebhookRecord.HostId = param.HostId
	gitWebhookRecord.Status = param.Status
	gitWebhookRecord.GitWebhookUpdateAt = param.GitWebhookUpdateAt
	gitWebhookRecord.SSHUrl = param.SSHUrl
	gitWebhookRecord.RecData = param.RecData
	if err = model.DB.Save(&gitWebhookRecord).Error; err != nil {
		return &gitWebhookRecord, fmt.Errorf("数据保存失败: %v", err)
	}
	return &gitWebhookRecord, err
}

func (s *GitWebhookService) GetGitWebhook(param api.SearchIdStringReq) (*[]model.GitWebhookRecord, int64, error) {
	var gitWebhookRecords []model.GitWebhookRecord
	var err error
	var total int64
	db := model.DB.Model(&model.GitWebhookRecord{})
	if param.Id != 0 {
		if err = db.Where("id = ?", param.Id).Find(&gitWebhookRecords).Error; err != nil {
			return nil, 0, fmt.Errorf("查询id错误: %v", err)
		}
	} else {
		searchReq := &api.SearchReq{
			Condition: db,
			Table:     &gitWebhookRecords,
			PageInfo:  param.PageInfo,
		}
		if param.String != "" {
			name := "%" + strings.ToUpper(param.String) + "%"
			db = model.DB.Where("UPPER(full_name) LIKE ?", name)
			searchReq.Condition = db
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
		} else {
			if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
				return nil, 0, err
			}
		}
	}
	return &gitWebhookRecords, total, err
}

func (s *GitWebhookService) DeleteGitWebhook(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.GitWebhookRecord{}, ids); err != nil {
		return err
	}
	if err = model.DB.Where("id IN (?)", ids).Delete(&model.GitWebhookRecord{}).Error; err != nil {
		return fmt.Errorf("记录删除失败: %v", err)
	}
	return err
}
