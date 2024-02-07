package model

// 按照对应Cloud官方文档的配置，定义CloudInstanceConfig结构体
type CloudInstanceConfig struct {
	Global                  `gorm:"embedded"`
	Region                  string `gorm:"type: varchar(50)"`
	InstanceChargeType      string `gorm:"type: varchar(50)"`
	Period                  int64  `gorm:"type: smallint"`
	RenewFlag               string `gorm:"type: varchar(50)"`
	ProjectId               int64  `gorm:"type: bigint"`
	InstanceFamily          string `gorm:"type: varchar(50)"`
	CpuCores                int    `gorm:"type: smallint"`
	MemorySize              int    `gorm:"type: smallint;comment:单位GB"`
	Fpga                    int    `gorm:"type: smallint"`
	GpuCores                int    `gorm:"type: smallint"`
	ImageId                 string `gorm:"type: varchar(50)"`
	SystemDiskType          string `gorm:"type: varchar(50)"`
	SystemDiskSize          int64  `gorm:"type: smallint"`
	DataDiskType            string `gorm:"type: varchar(50)"`
	DataDiskSize            int64  `gorm:"type: smallint"`
	VpcId                   string `gorm:"type: varchar(50)"`
	SubnetId                string `gorm:"type: varchar(50)"`
	InternetChargeType      string `gorm:"type: varchar(50)"`
	InternetMaxBandwidthOut int64  `gorm:"type: smallint"`
	InstanceNamePrefix      string `gorm:"type: varchar(50)"`
	SecurityGroupId         string `gorm:"type: varchar(50)"`
	HostNamePrefix          string `gorm:"type: varchar(50)"`
}
