package gitWebhook

import "time"

type UpdateGitWebhookStatusReq struct {
	Id     uint  `json:"id" form:"id" binding:"required"`         // 行id
	Status uint8 `json:"status" form:"status" binding:"required"` // 状态(0: 待审核 1: 待执行 2: 执行成功 3: 执行失败 4: 审核中 5: 已驳回)
}

type UpdateGitWebhookReq struct {
	Id                 uint      `json:"id" form:"id" binding:"required"` // ID
	FullName           string    `json:"full_name"`                       // git仓库全名
	ProjectId          uint      `json:"project_id"`                      // 项目ID
	HostId             uint      `json:"host_id"`                         // 服务器ID
	Status             uint8     `json:"status"`                          // 状态(0: 获取hook 1: 拉取包成功 2: 编译成功 3: 测试成功 4: 推送包到存储机器成功 5: 失败执行)
	GitWebhookUpdateAt time.Time `json:"webhook_update_at"`               // 仓库接收对应信号的更新时间
	SSHUrl             string    `json:"ssh_url"`                         // git的ssh_url
	RecData            []byte    `json:"rec_data"`                        // 从git仓库接受的webhook json
}
