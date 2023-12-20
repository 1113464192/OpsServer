package api

type WebsshConnReq struct {
	Hid        uint `json:"hid" form:"hid" binding:"required"` // 服务器id
	WindowSize      // 屏幕大小
}

type WindowSize struct {
	Hight  int `json:"hight" form:"hight" binding:"required"`   // 单位为字符
	Weight int `json:"weight" form:"weight" binding:"required"` // 单位为字符
}
