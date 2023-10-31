package model

type TaskTemplate struct {
	Global `gorm:"embedded"`
	// Status   uint     `json:"status" gorm:"comment: 状态(0: 审核中 1: 执行成功 2: 执行失败 3: 已驳回)"`
	TypeName  string  `json:"type_name" gorm:"type:varchar(10);comment: 类型名称, 最长10字符"`
	TaskName  string  `json:"task_name" gorm:"type:varchar(30);comment: 用户执行任务名"`
	Task      string  `json:"task" gorm:"type:longtext;comment: 用户执行任务内容,限Shell语言"`
	Comment   string  `json:"comment" gorm:"type:varchar(50);comment: 任务备注, 限长50"`
	Pid       uint    `json:"pid" gorm:"index;comment: 项目ID"`
	Condition string  `json:"condition" gorm:"type:text;comment:如: mem=5 代表1个服最少占用5G"`
	PortRule  string  `json:"port_rule" gorm:"type:text;comment: 端口规则, 如: 10000 + flag % 1000"`
	Args      string  `json:"args" gorm:"type: text;comment: 任意变量, 如: path=/data/a_b_c"`
	Hosts     []Host  `gorm:"many2many:template_host"`
	Project   Project `gorm:"foreignKey:Pid"`
}