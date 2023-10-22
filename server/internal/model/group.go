package model

import "database/sql"

type UserGroup struct {
	Global   `gorm:"embedded"`
	Name     string         `json:"name,omitempty" gorm:"type:varchar(10);unique;index"`
	ParentId uint16         `json:"parent_id" gorm:"comment:工作室下分项目组"`
	Mark     sql.NullString `json:"mark" gorm:"type:varchar(50);comment:标识"`
	Users    []User         `json:"users" gorm:"many2many:permit_users;"`
	Menus    []Menus        `json:"menus" gorm:"many2many:permit_menus;"`
	Project  []Project      `gorm:"foreignKey:GroupId;references:ID"`
}
