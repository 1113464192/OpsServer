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
}
