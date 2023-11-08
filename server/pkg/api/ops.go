package api

type GetExecParamReq struct {
	Tid uint `form:"tid" json:"tid"`
	Uid uint `form:"uid" json:"uid"`
}
