package api

type WebsshConnReq struct {
	Hid        uint `json:"hid" form:"hid" binding:"required"` // 服务器id
	WindowSize      // 屏幕大小
}

type WindowSize struct {
	Hight  int `json:"hight" form:"hight" binding:"required"`   // 单位为字符
	Weight int `json:"weight" form:"weight" binding:"required"` // 单位为字符
}

// 查询用户操作记录
type GetWebsshRecordReq struct {
	Id   uint   `json:"id" form:"id"`                            // 用户ID
	Date string `json:"string" form:"string" binding:"required"` // 年月，如：2006_01
	PageInfo
}
