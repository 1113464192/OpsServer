package api

// UpdateApiReq 新增或修改Api的请求格式
type UpdateApiReq struct {
	ID          uint   `json:"id" form:"id"`                                // 修改才需要传，没有传算新增
	Path        string `json:"path" form:"path"  binding:"required"`        // api路径
	Method      string `json:"method" form:"method" binding:"required"`     // 方法:创建/更新POST(默认)|查看GET|删除DELETE
	ApiGroup    string `json:"apiGroup" form:"apiGroup" binding:"required"` // api组
	Description string `json:"description" form:"description"`              // api中文描述
}

// CasbinInReceiveReq 分配用户API权限的请求格式
type CasbinInReceiveReq struct {
	GroupId string `json:"group_id"  binding:"required"` // 组id
	Ids     []uint `json:"ids"`
}

// CasbinGroupIds 用户组id
type CasbinGroupIds struct {
	GroupIds []uint `form:"group_id" json:"group_id"  binding:"required"`
}
