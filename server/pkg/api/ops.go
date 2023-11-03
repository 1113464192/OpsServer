package api

type GetTemplateParamReq struct {
	Tid uint `form:"tid" json:"tid"`
	Uid uint `form:"uid" json:"uid"`
}
