package webhook

import "fqhWeb/configs"

type WebhookService struct {
	GithubSecret string
	GitlabSecret string
}

var (
	insWebhook = &WebhookService{}
)

func Webhook() *WebhookService {
	insWebhook = &WebhookService{
		GithubSecret: configs.Conf.WebhookSecret.Github,
		GitlabSecret: configs.Conf.WebhookSecret.Gitlab,
	}
	return insWebhook
}
