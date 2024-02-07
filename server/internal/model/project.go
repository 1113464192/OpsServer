package model

type Project struct {
	Global                `gorm:"embedded"`
	Name                  string              `json:"name" gorm:"type:varchar(30);index;comment:项目名;unique"`
	Cloud                 string              `gorm:"comment: 云平台所属，用中文"`
	Status                uint                `json:"status" gorm:"comment:状态 1 正常 2 停摆"`
	UserId                uint                `json:"user_id" gorm:"comment:负责人用户ID;index"`
	GroupId               uint                `json:"group_id" gorm:"comment:关联组ID;index"`
	CloudInstanceConfigId uint                `json:"cloud_instance_config_id" gorm:"comment:云平台默认配置ID;index"`
	User                  User                `gorm:"foreignKey:UserId"`
	UserGroup             UserGroup           `gorm:"foreignKey:GroupId"`
	CloudInstanceConfig   CloudInstanceConfig `gorm:"foreignKey:CloudInstanceConfigId"`
	Hosts                 []Host              `gorm:"many2many:project_host"`
}
