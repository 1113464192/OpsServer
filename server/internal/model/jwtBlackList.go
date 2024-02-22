package model

type JwtBlacklist struct {
	Global `gorm:"embedded"`
	Jwt    string `gorm:"type:text;comment:auth"`
}
