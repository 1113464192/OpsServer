package configs

type Config struct {
	Mysql        Mysql        `json:"mysql"`
	Logger       Logger       `json:"logger"`
	SshConfig    SshConfig    `json:"ssh_timeout"`
	Webssh       Webssh       `json:"webssh"`
	System       System       `json:"system"`
	Concurrency  Concurrency  `json:"concurrency"`
	GitWebhook   GitWebhook   `json:"git_webhook"`
	SecurityVars SecurityVars `json:"security_vars"`
	ProjectWeb   ProjectWeb   `json:"project_web"`
}

type Mysql struct {
	Conf            string
	CreateBatchSize int
	TablePrefix     string
}

type Logger struct {
	Level string
}

type SshConfig struct {
	SshClientTimeout string
}

type Webssh struct {
	ReadBufferSize   int
	WriteBufferSize  int
	HandshakeTimeout string
	SshEcho          uint32
	SshTtyOpIspeed   uint32
	SshTtyOpOspeed   uint32
	MaxConnNumber    uint32
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

type ProjectWeb struct {
	RootPath string
}

var Conf = new(Config)
