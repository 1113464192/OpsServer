package api

type RunCmdtemRes struct {
	Cmd       string // 操作指令
	Host      string // ip
	Port      int64  // 端口
	Username  string // 用户名
	Password  string // 密码，填密码优先走密码，走公私钥不用传
	Key       []byte // 私钥字符串
	KeyPasswd []byte // 私钥密码(有就需要输入，没有不用传)
}

type GetResultRes struct {
	TaskId     uint
	OperatorId uint
	Response
}

type GetTemplateParamReq struct {
	Tid uint `form:"tid" json:"tid"`
	Uid uint `form:"uid" json:"uid"`
}
