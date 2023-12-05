package gitWebhook

import (
	"errors"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/pkg/util"
	"github.com/gin-gonic/gin"
	"io"
)

func (s *GitWebhookService) HandleGitlabWebhook(pid uint, hid uint, c *gin.Context) (err error) {
	fmt.Println(pid, hid)
	var data []byte
	if data, err = io.ReadAll(c.Request.Body); err != nil {
		return fmt.Errorf("获取数据的错误处理(入库): %v", err)
	}

	// 判断sign是否正确
	sign := c.GetHeader(consts.GITHUB_SECRET_SIGN)
	if !util.ValidatePrefix(data, []byte(GitWebhook().GitlabSecret), sign) {
		return errors.New("验证数据的错误处理(入库)")
	}
	// 判断类型执行命令————最终结果入库处理
	eventType := c.GetHeader(consts.GITHUB_EVENT)
	switch eventType {
	case consts.GITHUB_EVENT_PUSH:
		fmt.Println("github的push处理")
	case consts.GITHUB_EVENT_PR:
		fmt.Println("github的pull-request处理")
	default:
		fmt.Println("处理范围外的请求: " + eventType)
	}
	return err
}
