package api

type CSCmdRes struct {
	HostIp    string `json:"host_ip"`
	ServerDir string `json:"server_dir,omitempty"`
	Status    int    `json:"status"`
	Response  string `json:"response"`
}
