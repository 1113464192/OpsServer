package api

type UpdateGroupReq struct {
	ID       uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name     string `form:"name" json:"name" binding:"required"`
	ParentId uint16 `form:"parent_id" json:"parent_id"` // 工作室不用传，项目组传工作室的ID
	Mark     string `form:"mark" json:"mark"`
}

type UpdateUserAssReq struct {
	GroupID uint   `form:"group_id" json:"group_id" binding:"required"`
	UserIDs []uint `form:"user_id" json:"user_id" binding:"required"`
}

// 用户组结果返回
type GroupRes struct {
	ID       uint
	Name     string
	ParentId uint16
	Mark     string
	User     []any
}

type GetGroupReq struct {
	Name string `json:"name" form:"name"`
	PageInfo
}
