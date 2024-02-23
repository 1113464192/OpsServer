# OpsServer
swagger: http://106.52.66.254:9081/swagger/index.html#
> 仅供查看
> 注意：目前尚未编写完成，许多功能还有小报错需要处理，仅供学习以及面试官查看代码
>
> 前端能力薄弱，也没有时间编写，仅编写了后端接口作面试用

## 搭建条件

> go1.21也正常运行

go1.19+

mariadb:11.2

> 先在数据库创建一个admin用户再由他进行接口操作

shellScript下脚本给予755权限



## 功能讲解

> 代码中ssh模式部分注释，更推荐cs模式，更少性能占用以及更稳定，后期可以拓展为全http短链，返回从server端做个接口接收入库即可
>
> 云平台部分未测试(要烧钱)，都是根据官方文档结合本代码作改动，注释后功能正常运行

基本的RBAC，运维部份涵盖从CI到CD的全流程，结合了CD支持ssh模式或cs模式，以及webssh

### CI

1. ci通过各git平台或gitlab等私有平台自带的webhook功能，通过接口接收后执行后端人员自行编写的编译测试脚本
2. 本次请求入库保存，方便后续获取gitsshurl拉取编译后的代码
3. 后端人员自行编写自己代码的编译测试脚本，成功/失败返回对应json到/api/v1/git-gitWebhook/project-update-status接口进行编译状态更改
4. 接口自行更改数据库展示CI进度

### CD
> group(工作室)区分用户权限，project(项目)区分机器以及云商关联安全组等
1. 先由运维编写装服脚本
2. 在template模块进行操作编写 (可扩展json，允许带参数渲染cmd)  并提交
3. 自动筛选host表的机器，根据参数判断内存CPU等是否足够、端口是否冲突，逐个服判断并绑定
4. 设置多个审批人以及一个执行人
5. 审批人逐个审批完后，发送给php接口微信小程序通知执行人
6. 执行人通过get接口获取执行参数与命令和机器，确认无误后点击确认执行
7. SSH模式：goroutine并发执行ssh          C/S模式: goroutine并发发送命令给client端(https://github.com/1113464192/OpsClient)执行操作
8. 返回结果到任务记录表，方便后续用户查看结果
9. 如template的TypeName是"装服类型"，则将单服入库

> 此处两表结构举例，可以安装多个服多个条件，不足机器自动购买，可设置上限
>
> 一个是模块表一个是记录任务表

```golang
package model

type TaskTemplate struct {
	Global `gorm:"embedded"`
	// Status   uint     `json:"status" gorm:"comment: 状态(0: 审核中 1: 执行成功 2: 执行失败 3: 已驳回)"`
	TypeName   string       `json:"type_name" gorm:"type:varchar(10);comment: 类型名称, 最长10字符"`	// 如 更新类型
	TaskName   string       `json:"task_name" gorm:"type:varchar(30);comment: 用户执行任务名"`	// 如 热更新
	CmdTem     string       `json:"cmd_tem" gorm:"type:longtext;comment: 用户执行任务内容,限Shell语言, 待渲染变量参数格式:双大括号间隔空格包含.var"`	// 如 bash /tmp/install_server.sh a服 服类型
	ConfigTem  string       `json:"config_tem" gorm:"type:longtext;comment: 配置文件模板, 待渲染变量参数格式:双大括号间隔空格包含.var"`	// 允许一次安装多个服，代码自动筛选机器，不足则自动购买机器
	Comment    string       `json:"comment" gorm:"type:varchar(50);comment: 任务备注, 限长50"`
	Pid        uint         `json:"pid" gorm:"index;comment: 项目ID"`
	Condition  string       `json:"condition" gorm:"type:text;comment:如: mem=5 代表1个服最少占用5G"`	// 如 mem=5 代表1个服最少占用5G，进行机器筛选
	PortRule   string       `json:"port_rule" gorm:"type:text;comment: 端口规则,如: serverPort=10000 + flag % 1000 用模板变量名"`	// 端口规则，判断机器是否端口冲突，冲突则循环下一个机器
	Args       string       `json:"args" gorm:"type: text;comment: 任意变量, 如: path=/data/a_b_c,sftpPath=/data/a_b_c/server/config"`	// 变量以及值 如 path=/data/a, path=/data/b 等, 代码格式化为json
	Hosts      []Host       `gorm:"many2many:template_host"`
	Auditor    []User       `gorm:"many2many:auditor_task"`	// 审批人
	Project    Project      `gorm:"foreignKey:Pid"`
	TaskRecord []TaskRecord `gorm:"foreignKey:TemplateId;references:ID"`
}


package model

type TaskRecord struct {
	Global      `gorm:"embedded"`
	TaskName    string       `json:"task_name" gorm:"index;type:varchar(20);comment:  最长20字符"`
	TemplateId  uint         `json:"template_id" gorm:"comment: 对应模板id"`
	OperatorId  uint         `json:"type_name" gorm:"comment: 操作人ID"`
	Status      uint8        `json:"status" gorm:"comment: 状态(0: 待审核 1: 待执行 2: 执行成功 3: 执行失败 4: 审核中 5: 已驳回)"` // 状态(0: 待审核 1: 待执行 2: 执行成功 3: 执行失败 4: 审核中 5: 已驳回)
	Response    string       `json:"response" gorm:"type:longtext;comment: 返回值"`
	Args        string       `json:"args" gorm:"type:longtext;comment: 参数展示"`
	SSHJson     string       `json:"ssh_json" gorm:"type:longtext;comment: 包含ssh信息"`
	SFTPJson    string       `json:"sftp_json" gorm:"type:longtext;comment: 包含sftp信息"`
	NonApprover string       `json:"non_approver" gorm:"type:longtext;comment: 待批准的审核员"`
	User        User         `gorm:"foreignKey:OperatorId"`
	Template    TaskTemplate `gorm:"foreignKey:TemplateId"`
	Auditor     []User       `gorm:"many2many:auditor_task"`
}

```







入行新人，欢迎交流、指导、批评：fqh1113464192@163.com





