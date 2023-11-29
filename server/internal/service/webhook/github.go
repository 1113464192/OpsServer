package webhook

import (
	"encoding/json"
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/pkg/webhook"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGithubWebhook(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("获取数据的错误处理(入库)")
		return
	}

	// 判断sign是否正确
	sign := c.GetHeader(consts.GITHUB_SECRET_SIGN)
	if !webhook.ValidatePrefix(data, []byte(Webhook().GithubSecret), sign) {
		fmt.Println("验证数据的错误处理(入库)")
		return
	}

	// 判断类型执行命令————最终结果入库处理
	eventType := c.GetHeader(consts.GITHUB_EVENT)
	switch eventType {
	case consts.GITHUB_EVENT_PUSH:
		res, err := Webhook().handleGithubPushReq(data)
		fmt.Printf("\n github的push处理 \n %s \n %v \n", *res, err)
	case consts.GITHUB_EVENT_PR:
		res, err := Webhook().handleGithubPRReq(data)
		fmt.Printf("\n github的pull-request处理 \n %s \n %v \n", *res, err)
	default:
		fmt.Println("处理范围外的请求: " + eventType)
	}
	c.Status(http.StatusOK)
}

func (s *WebhookService) handleGithubPushReq(data []byte) (res *webhook.HandleGithubPushJson, err error) {
	res = &webhook.HandleGithubPushJson{}
	if err = json.Unmarshal(data, res); err != nil {
		return nil, fmt.Errorf("webhook json解析失败: %v", err)
	}
	return res, err
}

func (s *WebhookService) handleGithubPRReq(data []byte) (res *webhook.HandleGithubPRJson, err error) {
	res = &webhook.HandleGithubPRJson{}
	if err = json.Unmarshal(data, res); err != nil {
		return nil, fmt.Errorf("webhook json解析失败: %v", err)
	}
	return res, err
}
