package api

type TestSSHReq struct {
	HostId []uint `form:"host_id" json:"host_id"`
	UserId uint   `form:"user_id" json:"user_id"`
}

// type RunSSHCmdAsyncReq struct {
// 	HostIp     []string          `json:"host_ip"`
// 	Username   []string          `json:"username"`
// 	SSHPort    []string          `json:"ssh_port"`
// 	Password   map[string]string `json:"password"`
// 	Key        []byte            `json:"key"`
// 	Passphrase []byte            `json:"passphrase"`
// 	Cmd        []string          `json:"cmd"`
// }

type SSHClientConfigReq struct {
	HostIp     string `json:"host_ip"`
	Username   string `json:"username"`
	SSHPort    string `json:"ssh_port"`
	Password   string `json:"password"`
	Key        []byte `json:"key"`
	Passphrase []byte `json:"passphrase"`
	Cmd        string `json:"cmd"`
}

// type RunSFTPAsyncReq struct {
// 	HostIp      []string          `json:"host_ip"`
// 	Username    []string          `json:"username"`
// 	SSHPort     []string          `json:"ssh_port"`
// 	Password    map[string]string `json:"password"`
// 	Key         []byte            `json:"key"`
// 	Passphrase  []byte            `json:"passphrase"`
// 	Path        []string          `json:"path"`
// 	FileContent []string          `json:"file_content"`
// }

type SFTPClientConfigReq struct {
	HostIp      string `json:"host_ip"`
	Username    string `json:"username"`
	SSHPort     string `json:"ssh_port"`
	Password    string `json:"password"`
	Key         []byte `json:"key"`
	Passphrase  []byte `json:"passphrase"`
	Path        string `json:"path"`
	FileContent string `json:"file_content"`
}

// 返回更改
type SSHResultRes struct {
	HostIp   string
	Status   int
	Response string
}
