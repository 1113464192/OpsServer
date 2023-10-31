package api

type TaskTemRes struct {
	ID        uint
	TypeName  string
	TaskName  string
	Task      string
	Comment   string
	Pid       uint
	Condition string
	PortRule  string
	Args      string
}

type TaskInfo struct {
	TaskName string `json:"task_name"`
	ID       uint   `json:"id"`
}

type UpdateTaskTemplateReq struct {
	ID        uint     `form:"id" json:"id"`                                  // 修改才需要传，没有传算新增
	TypeName  string   `form:"type_name" json:"type_name" binding:"required"` // 模板类型名
	TaskName  string   `form:"task_name" json:"task_name" binding:"required"` // 模板名
	Task      string   `form:"task" json:"task" binding:"required"`           // 任务内容，规定用Shell
	Condition []string `json:"condition" form:"condition" binding:"required"` // mem=5就是单服最少5G，还有iowait/idle/load
	Comment   string   `form:"comment" json:"comment"`                        // 模板备注
	Pid       uint     `form:"pid" json:"pid" binding:"required"`             // 对应项目ID
	PortRule  []string `form:"port_rule" json:"port_rule"`                    // 端口规则, 如: 10000 + flag % 1000
	Args      []string `form:"args" json:"args"`                              // 任意变量, 如: path=/data/a_b_c
}

type GetProjectTaskReq struct {
	ID       uint   `form:"id" json:"id"`               // 传Task的ID查询，则无需填其它参数，返回Task的所有内容
	Pid      uint   `form:"pid" json:"pid"`             // 传Pid不传typename，返回对应Type的Name及其ID和 及Type包含的Task的Name和ID
	TypeName string `form:"type_name" json:"type_name"` // 需要精确到类型, 则传项目ID和类型名
	PageInfo
}

type UpdateTemplateAssHostReq struct {
	Tid  uint   `form:"tid" json:"tid"`
	Hids []uint `form:"hid" json:"hid" binding:"required"`
}