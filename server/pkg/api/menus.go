package api

type UpdateMenuReq struct {
	ID        uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name      string `form:"name" json:"name"`
	ParentId  uint   `form:"parent_id" json:"parent_id"` // 对应主菜单的ID
	Mark      string `form:"mark" json:"mark"`           // 前端标志
	Type      string `form:"type" json:"type"`           // 前端类型
	Title     string `form:"title" json:"title"`         // 前端展示菜单名
	Url       string `form:"url" json:"url"`             // 前端路由
	Sort      int    `form:"sort" json:"sort"`           // 排序标记
	Icon      string `form:"icon" json:"icon"`           // 菜单图标
	Author    string `form:"author" json:"author"`       // 创建人
	Component string `form:"component" json:"component"` // 组件
}

type UpdateMenuAssReq struct {
	MenuID   uint   `form:"menu_id" json:"menu_id" binding:"required"`
	GroupIDs []uint `form:"group_id" json:"group_id" binding:"required"`
}
