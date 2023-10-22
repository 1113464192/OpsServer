package model

import "database/sql"

type Host struct {
	Global         `gorm:"embedded"`
	Ipv4           sql.NullString `gorm:"type: varchar(30);index;comment: 如: 11.129.212.42"`
	Ipv6           sql.NullString `gorm:"type: varchar(100);comment: 如: 241d:c000:2022:601c:0:91aa:274c:e7ac/64"`
	Password       string         `gorm:"type: varchar(60);comment: 服务器密码加密后的字符串，一般机器都会禁止密码登录"`
	Zone           string         `gorm:"type: varchar(100);comment: 服务器所在地区"`
	ZoneTime       uint8          `gorm:"comment: 时区"`
	BillingType    uint8          `gorm:"comment: 1 按量收费, 2 包月收费, 3 包年收费 ...后续有需要再加"`
	Cost           float32        `gorm:"comment: 下次续费金额, 人民币为单位"`
	Cloud          string         `gorm:"comment: 云平台所属，用中文"`
	System         string         `gorm:"type: varchar(30)"`
	Iops           uint32         `gorm:"comment: 服务器IOPS"`
	Mbps           uint32         `gorm:"comment: 带宽, 单位为M"`
	Type           uint8          `gorm:"comment: 1 单服机器, 2 中央服机器, 3 CDN机器, 4 业务服机器  ...后续有需要再加"`
	Cores          uint16
	SystemDisk     uint32 `gorm:"comment: 系统盘, 单位为G"`
	DataDisk       uint32 `gorm:"comment: 数据盘, 单位为G"`
	Mem            uint32 `gorm:"comment: 单位为G"`
	CurrSystemDisk float32
	CurrDataDisk   float32
	CurrMem        float32
	CurrIowait     float32
	CurrIdle       float32
	CurrLoad       float32

	Domains  []Domain  `gorm:"many2many:host_domain"`
	Projects []Project `gorm:"many2many:project_host"`
}
