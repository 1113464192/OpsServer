package model

type Menus struct {
	Global     `gorm:"embedded"`
	Name       string      `json:"name,omitempty" gorm:"unique;type:varchar(30)"`
	ParentId   uint        `json:"parent_id" gorm:"comment:子菜单"`
	UserGroups []UserGroup `json:"group" gorm:"many2many:permit_menus;"`
	Mark       string      `json:"mark" gorm:"type:varchar(255);comment:简称;type:varchar(60)"`
	Type       string      `json:"type" gorm:"type:varchar(30);default:center;type:varchar(30)"`
	Title      string      `json:"title" gorm:"type:varchar(50);comment:菜单名;type:varchar(30)"` // 菜单名
	Url        string      `json:"path" gorm:"comment:路由;type:varchar(255)"`                   // 路由
	Sort       int         `json:"sort" gorm:"comment:排序标记"`                                   // 排序标记
	Icon       string      `json:"icon" gorm:"type:text;comment:菜单图标;type:varchar(255)"`       // 菜单图标
	Author     string      `json:"author" gorm:"comment:创建人;type:varchar(30)"`
	Component  string      `json:"component" gorm:"comment:组件;type:varchar(30)"` // 组件
}
