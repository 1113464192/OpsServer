package api

type UpdateHostReq struct {
	ID         uint    `form:"id" json:"id"`
	Ipv4       string  `form:"ipv4" json:"ipv4" binding:"required"`
	Ipv6       string  `form:"ipv6" json:"ipv6"`
	Pid        uint    `form:"pid" json:"pid" binding:"required"` // 项目ID
	Name       string  `form:"name" json:"name" binding:"required"`
	User       string  `form:"user" json:"user" binding:"required"`
	Password   []byte  `form:"password" json:"password"`
	Port       string  `form:"port" json:"port" binding:"required"`
	Zone       string  `form:"zone" json:"zone" binding:"required"`           // 所在地，用英文小写，如guangzhou
	ZoneTime   int8    `form:"zone_time" json:"zone_time" binding:"required"` // 时区，如东八区填8
	Cost       float32 `form:"cost" json:"cost"`                              // 下次续费金额, 人民币为单位
	Cloud      string  `form:"cloud" json:"cloud" binding:"required"`
	System     string  `form:"system" json:"system" binding:"required"`
	Mbps       uint32  `form:"mbps" json:"mbps" binding:"required"`
	Type       uint8   `form:"type" json:"type" binding:"required"`               // 1 单服机器, 2 中央服机器, 3 CDN机器, 4 业务服机器  ...后续有需要再加
	Cores      uint16  `form:"cores" json:"cores" binding:"required"`             // 四核输入4
	SystemDisk uint32  `form:"system_disk" json:"system_disk" binding:"required"` // 磁盘单位为G
	DataDisk   uint32  `form:"data_disk" json:"data_disk" binding:"required"`     // 磁盘单位为G
	Mem        uint32  `form:"mem" json:"mem" binding:"required"`                 // 内存单位为G
}

type GetHostReq struct {
	Ip       string `json:"ip" form:"ip"` // 查询IP则输入IP，v4或v6都可以
	PageInfo `form:"page_info" json:"page_info"`
}

type HostInfoRes struct {
	CurrSystemDisk []SSHResultRes
	CurrDataDisk   []SSHResultRes
	CurrMem        []SSHResultRes
	CurrIowait     []SSHResultRes
	CurrIdle       []SSHResultRes
	CurrLoad       []SSHResultRes
}

type HostRes struct {
	ID       uint
	Ipv4     string
	Ipv6     string
	Pid      uint
	Name     string
	Port     string
	Zone     string
	ZoneTime int8
	//	BillingType    uint8
	Cost           float32
	Cloud          string
	System         string
	Iops           uint32
	Mbps           uint32
	Type           uint8
	Cores          uint16
	SystemDisk     uint32
	DataDisk       uint32
	Mem            uint32
	CurrSystemDisk float32
	CurrDataDisk   float32
	CurrMem        float32
	CurrIowait     float32
	CurrIdle       float32
	CurrLoad       float32
}
