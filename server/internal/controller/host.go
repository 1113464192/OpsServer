package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateHost
// @Tags 服务器相关
// @title 新增/修改服务器
// @description 返回新增/修改的指定服务器
// @Summary 新增/修改服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateHostReq true "host传入参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/host [post]
func UpdateHost(c *gin.Context) {
	var hostReq api.UpdateHostReq
	if err := c.ShouldBind(&hostReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	host, err := service.Host().UpdateHost(&hostReq)
	if err != nil {
		logger.Log().Error("Host", "创建/修改服务器失败", err)
		c.JSON(500, api.Err("创建/修改服务器失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: host,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetHostPasswd
// @Tags 服务器相关
// @title 返回服务器的密码
// @description 返回服务器的密码
// @Summary 返回服务器的密码
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.IdReq true "传hostid"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/password [get]
func GetHostPasswd(c *gin.Context) {
	var param api.IdReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	passwd, err := service.Host().GetHostPasswd(param.Id)
	if err != nil {
		logger.Log().Error("Host", "获取服务器密码失败", err)
		c.JSON(500, api.Err("获取服务器密码失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: passwd,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// DeleteHost
// @Tags 服务器相关
// @title 删除服务器
// @description 返回success
// @Summary 删除服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "删除服务器ID切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/host [delete]
func DeleteHost(c *gin.Context) {
	var hostReq api.IdsReq
	if err := c.ShouldBind(&hostReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Host().DeleteHost(hostReq.Ids)
	if err != nil {
		logger.Log().Error("Host", "删除服务器失败", err)
		c.JSON(500, api.Err("删除服务器失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetHost
// @Tags 服务器相关
// @title 查询服务器
// @description 全部查询不传Ip
// @Summary 查询服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.SearchIdStringReq true "填写了ID则无需填其他参数，填写了IP模糊查询IP，都不填查询全部"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/host [get]
func GetHost(c *gin.Context) {
	var hostReq api.SearchIdStringReq
	if err := c.ShouldBind(&hostReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	host, count, err := service.Host().GetHost(&hostReq)
	if err != nil {
		logger.Log().Error("Host", "查询服务器失败", err)
		c.JSON(500, api.Err("查询服务器失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data:     host,
		Total:    count,
		Page:     hostReq.PageInfo.Page,
		PageSize: hostReq.PageInfo.PageSize,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}
