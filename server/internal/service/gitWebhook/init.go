package gitWebhook

import "fqhWeb/configs"

type GitWebhookService struct {
	GithubSecret string
	GitlabSecret string
}

var (
	insGitWebhook = &GitWebhookService{}
)

func GitWebhook() *GitWebhookService {
	insGitWebhook = &GitWebhookService{
		GithubSecret: configs.Conf.GitWebhook.GithubSecret,
		GitlabSecret: configs.Conf.GitWebhook.GitlabSecret,
	}
	return insGitWebhook
}
