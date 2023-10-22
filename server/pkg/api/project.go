package api

type UpdateProjectReq struct {
	ID      uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name    string `form:"name" json:"name" binding:"required"`
	Status  uint   `form:"status" json:"status" binding:"required"`
	UserId  uint   `form:"user_id" json:"user_id" binding:"required"`
	GroupId uint   `form:"group_id" json:"group_id" binding:"required"`
}

type UpdateProjectAssReq struct {
	ProjectID uint `form:"project_id" json:"project_id" binding:"required"`
	GroupID   uint `form:"group_id" json:"group_id" binding:"required"`
}

type UpdateProjectAssHostReq struct {
	Pid  uint   `form:"pid" json:"pid"`
	Hids []uint `form:"hid" json:"hid" binding:"required"`
}

type GetProjectReq struct {
	Name string `form:"name" json:"name" binding:"required"`
	PageInfo
}

type GetHostAssReq struct {
	ProjectId uint `form:"project_id" json:"project_id" binding:"required"`
	PageInfo
}

type ProjectRes struct {
	ID      uint
	Name    string
	Status  uint
	UserId  uint
	GroupId uint
}
