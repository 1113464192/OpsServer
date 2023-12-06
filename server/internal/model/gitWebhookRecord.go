package model

import "time"

type GitWebhookRecord struct {
	Global             `gorm:"embedded"`
	FullName           string    `json:"full_name" gorm:"comment: 完整仓库名"`
	ProjectId          uint      `json:"project_id" gorm:"comment: 对应项目id"`
	HostId             uint      `json:"host_id" gorm:"comment: 对应服务器id"`
	Status             uint8     `json:"status" gorm:"comment: 状态(0: 获取hook 1: 拉取包成功 2: 编译成功 3: 测试成功 4: 推送包到存储机器成功 5: 失败执行)"` // 状态(0: 获取hook 1: 拉取包成功 2: 编译成功 3: 测试成功 4: 推送包到存储机器成功 5: 失败执行)
	GitWebhookUpdateAt time.Time `json:"webhook_update_at" gorm:"comment: 仓库接收对应信号的更新时间"`                                     // 仓库接收对应信号的更新时间
	SSHUrl             string    `json:"ssh_url" gorm:"type:text;column:ssh_url;comment:git的ssh_url"`                         // git的ssh_url
	RecData            []byte    `json:"rec_data" gorm:"type:mediumblob;comment:由于webhook数据太多, 不同仓库也有不同返回, 此处包含接收的几乎所有信息"`    // 由于webhook数据太多, 不同仓库也有不同返回, 为防止以后过度更改表结构, 因此用长文本字段存储后续需要增加减少json
	ErrResponse        string    `json:"err_response" gorm:"type:text;comment:错误返回"`                                          // 错误返回
	Project            Project   `json:"project,omitempty" gorm:"foreignKey:ProjectId"`
	Host               Host      `json:"host,omitempty" gorm:"foreignKey:HostId"`
}
