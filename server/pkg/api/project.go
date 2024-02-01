package api

type UpdateProjectReq struct {
	ID      uint   `form:"id" json:"id"`                                // 修改才需要传，没有传算新增
	Name    string `form:"name" json:"name" binding:"required"`         // 项目名
	Cloud   string `form:"cloud" json:"cloud" binding:"required"`       // 云平台所属，用中文
	Status  uint   `form:"status" json:"status" binding:"required"`     // 状态：1 正常 2 停摆
	UserId  uint   `form:"user_id" json:"user_id" binding:"required"`   // 负责人用户ID
	GroupId uint   `form:"group_id" json:"group_id" binding:"required"` // 关联组ID
}

type UpdateProjectAssReq struct {
	ProjectID uint `form:"project_id" json:"project_id" binding:"required"`
	GroupID   uint `form:"group_id" json:"group_id" binding:"required"`
}

type UpdateProjectAssHostReq struct {
	Pid  uint   `form:"pid" json:"pid"`
	Hids []uint `form:"hid" json:"hid" binding:"required"`
}

type GetHostAssReq struct {
	ProjectId uint `form:"project_id" json:"project_id" binding:"required"`
	PageInfo
}

type ProjectRes struct {
	ID      uint
	Name    string
	Cloud   string
	Status  uint
	UserId  uint
	GroupId uint
}
