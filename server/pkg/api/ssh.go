package api

type TestSSHReq struct {
	HostId []uint `form:"host_id" json:"host_id"`
	UserId uint   `form:"user_id" json:"user_id"`
}

type RunSSHCmdAsyncReq struct {
	HostIp     []string          `json:"host_ip"`
	Username   []string          `json:"username"`
	SSHPort    []string          `json:"ssh_port"`
	Password   map[string]string `json:"password"`
	Key        []byte            `json:"key"`
	Passphrase []byte            `json:"passphrase"`
}

type SSHClientConfigReq struct {
	HostIp     string `json:"host_ip"`
	Username   string `json:"username"`
	SSHPort    string `json:"ssh_port"`
	Password   string `json:"password"`
	Key        []byte `json:"key"`
	Passphrase []byte `json:"passphrase"`
}

type SftpReq struct {
	DestFile string `json:"dest_file"`
}

// 返回更改
type SSHResultRes struct {
	HostIp   string
	Status   bool
	Response string
}
