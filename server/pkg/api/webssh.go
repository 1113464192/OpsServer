package api

type WebsshConnReq struct {
	Hid        uint `json:"hid" binding:"required"` // 服务器id
	WindowSize      // 屏幕大小
}

type WindowSize struct {
	Hight  int `json:"hight" binding:"required"`  // 单位为字符
	Weight int `json:"weight" binding:"required"` // 单位为字符
}
