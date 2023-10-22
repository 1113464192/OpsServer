package model

type CasbinRule struct {
	GroupId string `json:"group_id"`
	Path    string `json:"path"`
	Method  string `json:"method"`
}

func (c CasbinRule) TableName() string {
	return "casbin_rule"
}
