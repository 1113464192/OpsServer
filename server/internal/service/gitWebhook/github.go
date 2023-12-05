package gitWebhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api/gitWebhook"
	"fqhWeb/pkg/util"
	"github.com/gin-gonic/gin"
	"io"
)

func (s *GitWebhookService) HandleGithubWebhook(pid uint, hid uint, c *gin.Context) (err error) {
	var data []byte
	if data, err = io.ReadAll(c.Request.Body); err != nil {
		return fmt.Errorf("获取数据的错误处理(入库): %v", err)
	}

	// 判断sign是否正确
	sign := c.GetHeader(consts.GITHUB_SECRET_SIGN)
	if !util.ValidatePrefix(data, []byte(GitWebhook().GithubSecret), sign) {
		return errors.New("验证数据的错误处理(入库)")
	}

	// 判断类型执行命令————最终结果入库处理
	eventType := c.GetHeader(consts.GITHUB_EVENT)
	switch eventType {
	case consts.GITHUB_EVENT_PUSH:
		var res *gitWebhook.HandleGithubPushJson
		if err = GitWebhook().handleGithubPushReq(data, pid, hid); err != nil {
			return fmt.Errorf("执行webhook处理函数报错: %v", err)
		}

		fmt.Printf("\n github的push处理 \n %s \n %v \n", *res, err)
	case consts.GITHUB_EVENT_PR:
		var res *gitWebhook.HandleGithubPRJson
		err = GitWebhook().handleGithubPRReq(data)
		fmt.Printf("\n github的pull-request处理 \n %s \n %v \n", *res, err)
	default:
		fmt.Println("处理范围外的请求: " + eventType)
	}
	return err
}

func (s *GitWebhookService) writeGithubPushDataToDb(nData *gitWebhook.HandleGithubPushJson, allData []byte, pid uint, hid uint) (uint, error) {
	wh := &model.GitWebhookRecord{
		FullName:           nData.Repository.FullName,
		ProjectId:          pid,
		HostId:             hid,
		GitWebhookUpdateAt: nData.Repository.UpdatedAt,
		SSHUrl:             nData.Repository.SshUrl,
		RecData:            allData,
	}
	var err error
	if err = model.DB.Model(&model.GitWebhookRecord{}).Create(wh).Error; err != nil {
		return 0, fmt.Errorf("写入gitWebhook数据库时报错: %v", err)
	}
	return wh.ID, err
}

func (s *GitWebhookService) handleGithubPushReq(data []byte, pid uint, hid uint) (err error) {
	res := &gitWebhook.HandleGithubPushJson{}
	if err = json.Unmarshal(data, res); err != nil {
		return fmt.Errorf("gitWebhook json解析失败: %v", err)
	}
	var whId uint
	if whId, err = s.writeGithubPushDataToDb(res, data, pid, hid); err != nil {
		return err
	}
	//fmt.Println("============", whId, res.Repository.SshUrl, res.Repository.Name, hid, configs.Conf.GitWebhook.GitCiScriptDir, configs.Conf.GitWebhook.GitCiRepo, "============")
	if err = s.ExecServerCustomCi(whId, res.Repository.SshUrl, res.Repository.Name, hid); err != nil {
		return err
	}
	return err
}

func (s *GitWebhookService) handleGithubPRReq(data []byte) (err error) {
	res := &gitWebhook.HandleGithubPRJson{}
	if err = json.Unmarshal(data, res); err != nil {
		return fmt.Errorf("gitWebhook json解析失败: %v", err)
	}
	return err
}
