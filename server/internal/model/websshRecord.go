package model

import "database/sql"

type WebsshRecord struct {
	Global `gorm:"embedded"`
	UserId uint           `json:"user_id" gorm:"index"`                                     // userid
	HostId uint           `json:"host_id" gorm:"comment:服务器ID"`                             // 服务器ID
	Ipv4   sql.NullString `json:"ipv4" gorm:"comment:防止机器被删除, 保留一下IP号"`                     // ipv4
	Ipv6   sql.NullString `json:"ipv6" gorm:"comment:防止机器被删除, 保留一下IP号"`                     // ipv6
	Actlog []byte         `json:"act_log" gorm:"type:blob;comment:用户输入与服务器返回记录, 超过2048b截断"` // 服务器返回，超过2048b截断
	User   User           `gorm:"foreignKey:UserId"`
	Host   Host           `gorm:"foreignKey:HostId"`
}
