package model

type ServerList struct {
	Global        `gorm:"embedded"`
	Flag          string  `json:"flag" gorm:"type:varchar(20);index;comment:单服标识"`
	Path          string  `json:"path" gorm:"type:varchar(30);index;comment:单服标识"`
	ServerName    string  `json:"server_name" gorm:"type:varchar(30);index;comment:目录名"`
	HostId        uint    `json:"host_id" gorm:"comment: 对应服务器id"`
	ProjectId     uint    `json:"project_id" gorm:"comment: 对应项目id"`
	ServerHost    Host    `gorm:"foreignKey:HostId"`
	ServerProject Project `gorm:"foreignKey:ProjectId"`
}
