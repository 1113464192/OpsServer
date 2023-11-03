package api

type UpdateHostReq struct {
	ID          uint   `form:"id" json:"id"`
	Ipv4        string `form:"ipv4" json:"ipv4" binding:"required"`
	Ipv6        string `form:"ipv6" json:"ipv6"`
	User        string `form:"user" json:"user" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	Port        string `form:"port" json:"port" binding:"required"`
	Zone        string `form:"zone" json:"zone" binding:"required"`           // 所在地
	ZoneTime    uint8  `form:"zone_time" json:"zone_time" binding:"required"` // 时区，如东八区填8
	BillingType uint8  `form:"billing" json:"billing" binding:"required"`     // 1 按量收费, 2 包月收费, 3 包年收费 ...后续有需要再加
	Cloud       string `form:"cloud" json:"cloud" binding:"required"`
	System      string `form:"system" json:"system" binding:"required"`
	Iops        uint32 `form:"iops" json:"iops" binding:"required"`
	Mbps        uint32 `form:"mbps" json:"mbps" binding:"required"`
	Type        uint8  `form:"type" json:"type" binding:"required"`               // 1 单服机器, 2 中央服机器, 3 CDN机器, 4 业务服机器  ...后续有需要再加
	Cores       uint16 `form:"cores" json:"cores" binding:"required"`             // 四核输入4
	SystemDisk  uint32 `form:"system_disk" json:"system_disk" binding:"required"` // 磁盘单位为G
	DataDisk    uint32 `form:"data_disk" json:"data_disk" binding:"required"`     // 磁盘单位为G
	Mem         uint32 `form:"mem" json:"mem" binding:"required"`                 // 内存单位为G
	// CurrDisk    float32 `form:"curr_disk" json:"curr_disk"`
	// CurrMem     float32 `form:"curr_mem" json:"curr_mem"`
	// CurrIowait  float32 `form:"curr_iowait" json:"curr_iowait"`
	// CurrIdle    float32 `form:"curr_idle" json:"curr_idle"`
	// CurrLoad    float32 `form:"curr_load" json:"curr_load"`
}

type UpdateHostAssProjectReq struct {
	Pids []uint `form:"pid" json:"pid"`
	Hid  uint   `form:"hid" json:"hid" binding:"required"`
}

type GetHostReq struct {
	Ip       string `json:"ip" form:"ip"` // 查询IP则输入IP，v4或v6都可以
	PageInfo `form:"page_info" json:"page_info"`
}

type GetHostAssProjectReq struct {
	Id uint `json:"id" form:"id" binding:"required"`
	PageInfo
}

type UpdateDomainReq struct {
	Id    uint   `json:"id" form:"id"`
	Value string `json:"domain" form:"domain" binding:"required"`
}

type UpdateDomainAssHostReq struct {
	Did  uint   `form:"did" json:"did" binding:"required"`
	Hids []uint `form:"hids" json:"hids"`
}

type HostRes struct {
	ID             uint
	Ipv4           string
	Ipv6           string
	Port           string
	Zone           string
	ZoneTime       uint8
	BillingType    uint8
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
