package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateTemplate
// @Tags 模板相关
// @title 新增/修改任务模板
// @description 运营点击发出工单/运维审批最后确认 都可以修改
// @Summary 新增/修改任务模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateTemplateReq true "更新taskTem所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/template/template [post]
func UpdateTemplate(c *gin.Context) {
	var taskReq api.UpdateTemplateReq
	if err := c.ShouldBind(&taskReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	result, err := service.Template().UpdateTemplate(&taskReq)
	if err != nil {
		logger.Log().Error("Task", "创建/修改任务模板失败", err)
		c.JSON(500, api.Err("创建/修改任务模板失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: result,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetProjectTemplate
// @Tags 模板相关
// @title 获取任务模板
// @description 传ID(Task)返回模板内容/只传项目ID返回包含任务类型/传任务类型和项目ID返回包含模板名
// @Summary 获取任务模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetProjectTemplateReq true "传ID(Template)返回模板内容/只传项目ID返回包含任务类型/传任务类型和项目ID返回包含模板名"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/template/template [get]
func GetProjectTemplate(c *gin.Context) {
	var templateReq api.GetProjectTemplateReq
	if err := c.ShouldBind(&templateReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	result, total, err := service.Template().GetProjectTemplate(&templateReq)
	if err != nil {
		logger.Log().Error("Task", "创建/修改任务模板失败", err)
		c.JSON(500, api.Err("创建/修改任务模板失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: result,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     templateReq.PageInfo.Page,
		PageSize: templateReq.PageInfo.PageSize,
	})
}

// DeleteTemplate
// @Tags 模板相关
// @title 删除任务模板
// @description 删除的任务模板
// @Summary 删除任务模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/template/template [delete]
func DeleteTemplate(c *gin.Context) {
	var taskReq api.IdsReq
	if err := c.ShouldBind(&taskReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Template().DeleteTemplate(taskReq.Ids)
	if err != nil {
		logger.Log().Error("Task", "创建/修改任务模板失败", err)
		c.JSON(500, api.Err("创建/修改任务模板失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// UpdateTemplateAssHost
// @Tags 模板相关
// @title 关联服务器
// @description 服务器ID[多选](如果直接使用对应项目关联主机则无需关联主机)
// @Summary 关联服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateTemplateAssHostReq true "关联传入参数(如果直接使用对应项目关联主机则无需关联主机)"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/template/ass-host [put]
func UpdateTemplateAssHost(c *gin.Context) {
	var TaskReq api.UpdateTemplateAssHostReq
	if err := c.ShouldBind(&TaskReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Template().UpdateTemplateAssHost(TaskReq)
	if err != nil {
		logger.Log().Error("Task", "关联服务器失败", err)
		c.JSON(500, api.Err("关联服务器失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetConditionSet
// @Tags 模板相关
// @title 获取可输入条件集合
// @description 可不选或多选,有需要再让运维从代码中添加功能(opsservice也要添加)
// @Summary 获取可输入条件集合
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/template/condition-set [get]
func GetConditionSet(c *gin.Context) {
	m := map[uint]string{
		1: "data_disk	单位G",
		2: "mem		单位G",
		3: "iowait	单位%",
		4: "idle	单位%",
		5: "load	单位%",
	}
	c.JSON(200, api.Response{
		Data: m,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}
