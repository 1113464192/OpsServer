package controller

import (
	"fqhWeb/internal/service/ops"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GetExecParam
// @Tags Ops相关
// @title 提取任务模板内容执行时的参数
// @description 传入模板id，返回ssh执行所需参数
// @Summary 提取任务模板内容执行时的参数
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetExecParamReq true "传入所需id"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ops/getExecParam [get]
func GetExecParam(c *gin.Context) {
	var param api.GetExecParamReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	params, sftpParams, err := ops.Ops().GetExecParam(param)
	if err != nil {
		logger.Log().Error("Task", "获取ssh执行参数", err)
		c.JSON(500, api.Err("获取Ops任务执行参数失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: map[string]any{
			"sshReq":  *params,
			"sftpReq": *sftpParams,
		},
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// ①先从项目中获取操作的机器和命令模板
// 不定长参数接收参数

// ①先从项目中获取操作的机机器
// clientConfig := &sshService.ClientConfigService{
// 	Host:      ,
// 	Port:      22,
// 	Username:  "root",
// 	Password:  clientPassword,
// 	Key:       clientKey,
// 	KeyPasswd: clientKeyPasswd,
// }

// 	// ②走工单审批
// 	fmt.Println("\n ②走工单审批 \n")

// 	// ③执行操作

// }
