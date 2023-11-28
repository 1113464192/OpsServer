package api

type GetServerRecordReq struct {
	Id         uint   `form:"id" json:"id"`
	Pid        uint   `form:"pid" json:"pid"` // 如果不输入ID则必填
	Flag       string `form:"flag" json:"flag"`
	ServerName string `form:"server_name" json:"server_name"`
	PageInfo
}

type UpdateServerRecordReq struct {
	Id         uint   `form:"id" json:"id" binding:"required"`
	Flag       string `form:"flag" json:"flag" binding:"required"`
	Path       string `form:"path" json:"path" binding:"required"`
	ServerName string `form:"server_name" json:"server_name" binding:"required"`
	HostId     uint   `form:"host_id" json:"host_id" binding:"required"`
	ProjectId  uint   `form:"project_id" json:"project_id" binding:"required"`
}

type ServerRecordRes struct {
	Id         uint   `json:"id"`
	Flag       string `json:"flag"`
	Path       string `json:"path"`
	ServerName string `json:"server_name"`
	HostId     uint   `json:"host_id"`
	ProjectId  uint   `json:"project_id"`
}
