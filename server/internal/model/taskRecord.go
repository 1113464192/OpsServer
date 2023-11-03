package model

import "database/sql"

type TaskRecord struct {
	Global     `gorm:"embedded"`
	TaskName   string        `json:"task_name" gorm:"type:varchar(20);comment:  最长20字符"`
	TemplateId sql.NullInt64 `json:"template_id" gorm:"comment: 对应模板id, 没有则NULL"`
	OperatorId uint          `json:"type_name" gorm:"index;comment: 操作人ID"`
	Status     uint8         `json:"status" gorm:"type:varchar(10);comment: 状态(0: 待审核 1: 执行成功 2: 执行失败 3: 已驳回 5: 已确认)"` // 状态(0: 待审核 1: 执行成功 2: 执行失败 3: 已驳回 5: 已确认)
	Response   string        `json:"response" gorm:"type:longtext;comment: 返回值"`
	User       User          `gorm:"foreignKey:OperatorId"`
	Hosts      []Host        `gorm:"many2many:task_host"`
	Auditor    []User        `gorm:"many2many:auditor_task"`
}
