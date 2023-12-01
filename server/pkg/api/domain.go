package api

type UpdateDomainReq struct {
	Id    uint   `json:"id" form:"id"`
	Value string `json:"domain" form:"domain" binding:"required"`
}

type UpdateDomainAssHostReq struct {
	Did  uint   `form:"did" json:"did" binding:"required"`
	Hids []uint `form:"hids" json:"hids"`
}
