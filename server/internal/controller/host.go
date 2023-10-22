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
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/update [post]
func UpdateHost(c *gin.Context) {
	var hostReq api.UpdateHostReq
	if err := c.ShouldBind(&hostReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	host, err := service.Host().UpdateHost(&hostReq)
	if err != nil {
		logger.Log().Error("Host", "创建/修改服务器", err)
		c.JSON(500, api.Err("创建/修改服务器", err))
		return
	}
	c.JSON(200, api.Response{
		Data: host,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// UpdateProjectAss
// @Tags 服务器相关
// @title 关联项目
// @description 项目ID[多选]
// @Summary 关联项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateHostAssProjectReq true "关联传入参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/association [put]
func UpdateProjectAss(c *gin.Context) {
	var hostReq api.UpdateHostAssProjectReq
	if err := c.ShouldBind(&hostReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Host().UpdateProjectAss(&hostReq)
	if err != nil {
		logger.Log().Error("Host", "关联项目", err)
		c.JSON(500, api.Err("关联项目失败", err))
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
// @Param data body api.GetHostReq true "获取host的参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/getHost [post]
func GetHost(c *gin.Context) {
	var hostReq api.GetHostReq
	if err := c.ShouldBind(&hostReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	host, count, err := service.Host().GetHost(&hostReq)
	if err != nil {
		logger.Log().Error("Host", "查询服务器", err)
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

// GetProjectAss
// @Tags 服务器相关
// @title 查询服务器对应项目
// @description 返回项目切片
// @Summary 查询服务器对应项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.GetHostAssProject true "传参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/getProject [get]
func GetProjectAss(c *gin.Context) {
	var hostReq api.GetHostAssProject
	if err := c.ShouldBind(&hostReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	project, total, err := service.Host().GetProject(&hostReq)
	if err != nil {
		logger.Log().Error("Host", "查询服务器关联项目", err)
		c.JSON(500, api.Err("查询服务器关联项目失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: project,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     hostReq.PageInfo.Page,
		PageSize: hostReq.PageInfo.PageSize,
	})
}

// DeleteHost
// @Tags 服务器相关
// @title 删除服务器
// @description 传服务器ID
// @Summary 删除服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param id body uint true "需删除的host ID"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/host/delete [delete]
func DeleteHost(c *gin.Context) {
	hidStr := c.PostForm("id")
	err := service.Host().DeleteHost(&hidStr)
	if err != nil {
		logger.Log().Error("Host", "删除服务器", err)
		c.JSON(500, api.Err("删除服务器失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// 新增/修改域名

// 域名关联服务器

// 查询域名绑定的服务器
