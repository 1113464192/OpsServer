package model

import "database/sql"

type Host struct {
	Global         `gorm:"embedded"`
	Ipv4           sql.NullString `gorm:"type: varchar(30);index;comment: 如: 11.129.212.42"`
	Ipv6           sql.NullString `gorm:"type: varchar(100);comment: 如: 241d:c000:2022:601c:0:91aa:274c:e7ac/64"`
	Pid            uint           `gorm:"comment: 项目ID;index"`
	Name           string         `gorm:"type: varchar(30);comment: 服务器名"`
	User           string         `gorm:"type: varchar(20)"`
	Password       []byte         `gorm:"type: blob;comment: 服务器密码加密后的字符串，一般机器都会禁止密码登录"`
	Port           string         `gorm:"type: varchar(10);comment: SSH端口"`
	Zone           string         `gorm:"type: varchar(100);comment: 服务器所在地区,用英文小写，如guangzhou"`
	ZoneTime       int8           `gorm:"comment: 时区"`
	Cost           float32        `gorm:"comment: 下次续费金额, 人民币为单位"`
	Cloud          string         `gorm:"comment: 云平台所属，用中文"`
	System         string         `gorm:"type: varchar(30)"`
	Iops           uint32         `gorm:"comment: 服务器IOPS"`
	Mbps           uint32         `gorm:"comment: 带宽, 单位为M"`
	Type           uint8          `gorm:"comment: 1 单服机器, 2 中央服机器, 3 CDN机器, 4 业务服机器  ...后续有需要再加"`
	Cores          uint16
	SystemDisk     uint32 `gorm:"comment: 系统盘, 单位为G"`
	DataDisk       uint32 `gorm:"comment: 数据盘, 单位为G"`
	Mem            uint64 `gorm:"comment: 单位为M"`
	CurrSystemDisk float32
	CurrDataDisk   float32
	CurrMem        float32
	CurrIowait     float32
	CurrIdle       float32
	CurrLoad       float32

	Project Project  `gorm:"foreignKey:Pid"`
	Domains []Domain `gorm:"many2many:host_domain"`
	//Projects     []Project      `gorm:"many2many:project_host"`
	TaskTemplate []TaskTemplate `gorm:"many2many:template_host"`
}
