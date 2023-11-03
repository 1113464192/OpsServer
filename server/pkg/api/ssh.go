package api

type RunSSHCmdAsyncReq struct {
	HostIp     []string `json:"host_ip"`
	Username   []string `json:"username"`
	SSHPort    []string `json:"ssh_port"`
	Password   []string `json:"password"`
	Key        []byte   `json:"key"`
	Passphrase []byte   `json:"passphrase"`
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

type SSHResultRes struct {
	HostIps []string
	Status  bool
	Response
}

type GetSSHRes struct {
	HostIps string
	Response
}
