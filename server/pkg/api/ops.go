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
	SSHReqs     []TaskRecordSSHRes
	NonApprover []uint
	Auditor     []uint
}

type TaskRecordSSHRes struct {
	HostIp      string `json:"host_ip"`
	Username    string `json:"username"`
	SSHPort     string `json:"ssh_port"`
	Cmd         string `json:"cmd"`
	Path        string `json:"path"`
	FileContent string `json:"file_content"`
}

type ApproveTaskReq struct {
	Id     uint  `json:"id" form:"id" binding:"required"`
	Status uint8 `json:"status" form:"status" binding:"required"` // 1:通过 4:驳回
}
