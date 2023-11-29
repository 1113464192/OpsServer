package webhook

import (
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/pkg/webhook"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGitlabWebhook(c *gin.Context) {
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
		fmt.Println("github的push处理")
	case consts.GITHUB_EVENT_PR:
		fmt.Println("github的pull-request处理")
	default:
		fmt.Println("处理范围外的请求: " + eventType)
	}
	c.Status(http.StatusOK)
}
