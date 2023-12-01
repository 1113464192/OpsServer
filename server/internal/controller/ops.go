package controller

import (
	"fqhWeb/internal/service/ops"
	"fqhWeb/pkg/api"
	apiOps "fqhWeb/pkg/api/ops"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util/jwt"

	"github.com/gin-gonic/gin"
)

// SubmitTask
// @Tags Ops相关
// @title 提交执行工单
// @description 传入模板id，返回ssh执行所需参数并自动写入任务工单库
// @Summary 提交执行工单
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body apiOps.SubmitTaskReq true "注意Auditor参数: 最先审批的放第一个,因为接入后从第一个到最后一个依次发送信息审批"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ops/submit-task [post]
func SubmitTask(c *gin.Context) {
	var param apiOps.SubmitTaskReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	taskRecord, err := ops.Ops().SubmitTask(param)
	if err != nil {
		logger.Log().Error("Task", "提交执行工单", err)
		c.JSON(500, api.Err("提交执行工单失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: taskRecord,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetTask
// @Tags Ops相关
// @title 查看任务工单
// @description 传入查询所需参数,输了ID就不用name和页码
// @Summary 查看任务工单
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.SearchIdStringReq true "传入所需参数,输了ID就不用name和页码"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ops/task [get]
func GetTask(c *gin.Context) {
	var param api.SearchIdStringReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	data, total, err := ops.Ops().GetTask(&param)
	if err != nil {
		logger.Log().Error("Task", "查看任务工单", err)
		c.JSON(500, api.Err("查看任务工单失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: data,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     param.PageInfo.Page,
		PageSize: param.PageInfo.PageSize,
	})
}

// GetSSHExecParam
// @Tags Ops相关
// @title 提取SSH执行参数
// @description 返回sftp和ssh的执行参数
// @Summary 提取SSH执行参数
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.IdReq true "传入所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ops/ssh-exec-param [get]
func GetSSHExecParam(c *gin.Context) {
	var param api.IdReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	sshReq, sftpReq, err := ops.Ops().GetSSHExecParam(param.Id)
	if err != nil {
		logger.Log().Error("Task", "获取ssh执行参数", err)
		c.JSON(500, api.Err("获取Ops任务执行参数失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: map[string]any{
			"SSHReq":  *sshReq,
			"SFTPReq": *sftpReq,
		},
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// ApproveTask
// @Tags Ops相关
// @title 用户审批工单
// @description 传入工单的ID
// @Summary 用户审批工单
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData apiOps.ApproveTaskReq true "传入工单的ID和是否成功"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ops/approve-task [put]
func ApproveTask(c *gin.Context) {
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
	}
	var param apiOps.ApproveTaskReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	userId := claims.User.ID
	res, err := ops.Ops().ApproveTask(param, userId)
	if err != nil {
		logger.Log().Error("Task", "提交执行工单", err)
		c.JSON(500, api.Err("提交执行工单失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: res,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// DeleteTask
// @Tags Ops相关
// @title 工单删除
// @description 传入工单的ID
// @Summary 工单删除
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "传入工单的ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ops/task [delete]
func DeleteTask(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := ops.Ops().DeleteTask(param.Ids)
	if err != nil {
		logger.Log().Error("Task", "删除工单", err)
		c.JSON(500, api.Err("删除工单失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// OpsExecSSHTask
// @Tags Ops相关
// @title 执行人执行工单操作
// @description 返回执行结果
// @Summary 执行人执行工单操作
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.IdReq true "传入工单的ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ops/exec-ssh-task [post]
func OpsExecSSHTask(c *gin.Context) {
	var param api.IdReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	data, err := ops.Ops().OpsExecSSHTask(param.Id)
	if err != nil {
		logger.Log().Error("Task", "执行工单操作", err)
		c.JSON(500, api.Err("执行工单操作失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: data,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}
