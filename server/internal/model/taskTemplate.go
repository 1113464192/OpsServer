package model

type TaskTemplate struct {
	Global `gorm:"embedded"`
	// Status   uint     `json:"status" gorm:"comment: 状态(0: 审核中 1: 执行成功 2: 执行失败 3: 已驳回)"`
	TypeName   string       `json:"type_name" gorm:"type:varchar(10);comment: 类型名称, 最长10字符"`                           // 如 更新类型
	TaskName   string       `json:"task_name" gorm:"type:varchar(30);comment: 用户执行任务名"`                                // 如 热更新
	CmdTem     string       `json:"cmd_tem" gorm:"type:longtext;comment: 用户执行任务内容,限Shell语言, 待渲染变量参数格式:双大括号间隔空格包含.var"` // 如 bash /tmp/install_server.sh a服 服类型
	ConfigTem  string       `json:"config_tem" gorm:"type:longtext;comment: 配置文件模板, 待渲染变量参数格式:双大括号间隔空格包含.var"`         // 允许一次安装多个服，代码自动筛选机器，不足则自动购买机器
	Comment    string       `json:"comment" gorm:"type:varchar(50);comment: 任务备注, 限长50"`
	Pid        uint         `json:"pid" gorm:"index;comment: 项目ID"`
	Condition  string       `json:"condition" gorm:"type:text;comment:如: mem=5 代表1个服最少占用5G"`                                      // 如 mem=5 代表1个服最少占用5G，进行机器筛选
	PortRule   string       `json:"port_rule" gorm:"type:text;comment: 端口规则,如: serverPort=10000 + flag % 1000 用模板变量名"`            // 端口规则，判断机器是否端口冲突，冲突则循环下一个机器
	Args       string       `json:"args" gorm:"type: text;comment: 任意变量, 如: path=/data/a_b_c,sftpPath=/data/a_b_c/server/config"` // 变量以及值 如 path=/data/a, path=/data/b 等
	Hosts      []Host       `gorm:"many2many:template_host"`
	Auditor    []User       `gorm:"many2many:auditor_task"` // 审批人
	Project    Project      `gorm:"foreignKey:Pid"`
	TaskRecord []TaskRecord `gorm:"foreignKey:TemplateId;references:ID"`
}
