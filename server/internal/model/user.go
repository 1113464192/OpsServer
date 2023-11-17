package model

type User struct {
	Global     `gorm:"embedded"`
	Username   string       `json:"username" gorm:"type:varchar(30);index"`
	Password   string       `json:"password" gorm:"type:varchar(60)"`
	Name       string       `json:"name" gorm:"type:varchar(10)"`
	Status     uint8        `json:"status" gorm:"comment: 1是正常，2是禁用;default:1"`
	Email      string       `json:"email" gorm:"type:varchar(255)"`
	Mobile     string       `json:"mobile" gorm:"type:varchar(15)"`
	LoginTime  int64        `json:"login_time"`
	Expiration uint64       `json:"expiration"`
	IsAdmin    uint8        `json:"is_admin" gorm:"default:0"`
	PriKey     []byte       `json:"pri_key" gorm:"type:blob;comment: 用户私钥，不传为NULL"`
	Passphrase []byte       `json:"passphrase" gorm:"type:blob;comment: 用户私钥通行密码，不传为NULL"`
	UserGroups []UserGroup  `json:"user_group" gorm:"many2many:permit_users;comment: 管理用户填1"`
	Project    []Project    `gorm:"foreignKey:UserId;references:ID"`
	TaskRecord []TaskRecord `gorm:"foreignKey:OperatorId;references:ID"`
}
