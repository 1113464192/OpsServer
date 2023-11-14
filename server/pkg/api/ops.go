package api

type SubmitTaskReq struct {
	Tid     uint   `form:"tid" json:"tid"`
	Uid     uint   `form:"uid" json:"uid"`
	Auditor []uint `form:"auditorId" json:"auditorId"`
}

type GetTaskReq struct {
	Tid      uint   `form:"tid" json:"tid"`
	TaskName string `form:"task_name" json:"task_name"`
	PageInfo
}

type TaskRecordRes struct {
	ID          uint
	TaskName    string
	TemplateId  uint
	OperatorId  uint
	Status      uint8
	Response    string
	HostIp      []string
	Username    []string
	SSHPort     []string
	Cmd         []string
	ConfigPath  []string
	FileContent []string
	NonApprover []uint
	Auditor     []uint
}
