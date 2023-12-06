package configs

type Config struct {
	Mysql        Mysql        `json:"mysql"`
	Logger       Logger       `json:"logger"`
	SshTimeout   SshTimeout   `json:"ssh_timeout"`
	System       System       `json:"system"`
	Concurrency  Concurrency  `json:"concurrency"`
	GitWebhook   GitWebhook   `json:"git_webhook"`
	SecurityVars SecurityVars `json:"security_vars"`
}

type Mysql struct {
	Conf            string
	CreateBatchSize int
}

type Logger struct {
	Level string
}

type SshTimeout struct {
	SshTimeout string
}

type Concurrency struct {
	Number int64
}

type System struct {
	Mode string
}

type GitWebhook struct {
	GithubSecret   string
	GitlabSecret   string
	GitCiScriptDir string
	GitCiRepo      string
}

type SecurityVars struct {
	AesKey              string
	AesIv               string
	TokenExpireDuration string
	TokenKey            string
	CiMd5Key            string
}

var Conf = new(Config)
