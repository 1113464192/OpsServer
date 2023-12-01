package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateServerRecord
// @Tags 服务端相关
// @title 更改单服记录列表
// @description 传入更改所需参数
// @Summary 更改单服记录列表
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.UpdateServerRecordReq true "传入所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/server/record [put]
func UpdateServerRecord(c *gin.Context) {
	var param api.UpdateServerRecordReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	data, err := service.Server().UpdateServerRecord(param)
	if err != nil {
		logger.Log().Error("Server", "更改单服记录列表", err)
		c.JSON(500, api.Err("更改单服记录列表失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: data,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetServerRecord
// @Tags 服务端相关
// @title 查看单服记录列表
// @description 传入查询所需参数,输了ID就不用Flag和页码以及PID，没传ID的话则必填项目ID，然后flag和name二选一加页码
// @Summary 查看单服记录列表
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetServerRecordReq true "传入所需参数,输了ID就不用name和Flag和页码"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/server/record [get]
func GetServerRecord(c *gin.Context) {
	var param api.GetServerRecordReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	data, total, err := service.Server().GetServerRecord(param)
	if err != nil {
		logger.Log().Error("Server", "查询单服记录列表", err)
		c.JSON(500, api.Err("查询单服记录列表失败", err))
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

// DeleteServerRecord
// @Tags 服务端相关
// @title 删除单服记录
// @description 删除成功返回success
// @Summary 删除单服记录
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "单服记录ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/server/record [delete]
func DeleteServerRecord(c *gin.Context) {
	var rid api.IdsReq
	if err := c.ShouldBind(&rid); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Server().DeleteServerRecord(rid.Ids)
	if err != nil {
		logger.Log().Error("Server", "删除单服记录", err)
		c.JSON(500, api.Err("删除单服记录失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}
