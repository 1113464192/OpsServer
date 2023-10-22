package api

import "gorm.io/gorm"

type SearchReq struct {
	Condition *gorm.DB // GORM查询规则条件
	Table     any
	PageInfo  PageInfo
}

type AssQueryReq struct {
	Condition *gorm.DB // GORM查询规则条件
	Table     any
	PageInfo  PageInfo
}
