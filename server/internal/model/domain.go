package model

type Domain struct {
	Global `gorm:"embedded"`
	Value  string `gorm:"type: varchar(255)"`
	Hosts  []Host `gorm:"many2many:host_domain"`
}
