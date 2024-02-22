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
	Cloud        Cloud        `json:"cloud"`
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
	ClientReqMd5Key     string
}

type Cloud struct {
	AllowConsecutiveCreateTimes int
	TencentCloud                struct {
		Ak string
		Sk string
	} `json:"tencent_cloud"`
	AliyunCloud struct {
		Ak string
		Sk string
	} `json:"aliyun_cloud"`
	VolcengineCloud struct {
		Ak string
		Sk string
	} `json:"volcengine_cloud"`
}

var Conf = new(Config)
