package api

type UpdateUserReq struct {
	ID         uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Username   string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Name       string `form:"name" json:"name" binding:"required"`
	Expiration uint64 `form:"expiration" json:"expiration" binding:"required"`
	Mobile     string `form:"mobile" json:"mobile" binding:"required"`
	Email      string `form:"email" json:"email" binding:"required"`
	IsAdmin    uint8  `from:"is_admin" json:"is_admin"` // 管理员传1，不然不用传
}

type PasswordReq struct {
	ID       uint   `form:"id" json:"id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=20"` // 要求密码长度不小于8不大于20
}

type StatusReq struct {
	ID     uint  `form:"id" json:"id" binding:"required"`
	Status uint8 `form:"status" json:"status" binding:"required"`
}

type AuthLoginReq struct {
	Username string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=20"`
}

// 登录返回
type AuthLoginRes struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	GroupIds []uint `json:"group_id"`
	Name     string `json:"name"`
	Token    string `json:"token"`
}

// 用户结果返回
type UserRes struct {
	ID         uint
	Username   string
	Name       string
	Status     uint8
	Email      string
	Mobile     string
	LoginTime  uint64
	Expiration uint64
	IsAdmin    uint8
}
