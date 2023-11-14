package model

type TaskRecord struct {
	Global      `gorm:"embedded"`
	TaskName    string       `json:"task_name" gorm:"type:varchar(20);comment:  最长20字符"`
	TemplateId  uint         `json:"template_id" gorm:"comment: 对应模板id"`
	OperatorId  uint         `json:"type_name" gorm:"index;comment: 操作人ID"`
	Status      uint8        `json:"status" gorm:"type:varchar(10);comment: 状态(0: 待审核 1: 待执行 2: 执行成功 3: 执行失败 4: 已驳回)"` // 状态(0: 待审核 1: 待执行 2: 执行成功 3: 执行失败 4: 已驳回)
	Response    string       `json:"response" gorm:"type:longtext;comment: 返回值"`
	User        User         `gorm:"foreignKey:OperatorId"`
	Template    TaskTemplate `gorm:"foreignKey:TemplateId"`
	SSHJson     string       `json:"ssh_json" gorm:"type:longtext;comment: 包含ssh信息"`
	NonApprover string       `json:"non_approver" gorm:"type:longtext;comment: 待批准的审核员"`
	Auditor     []User       `gorm:"many2many:auditor_task"`
}
