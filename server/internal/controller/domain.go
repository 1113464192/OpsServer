package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateDomain
// @Tags 域名相关
// @title 新增/修改域名
// @description 返回新增或修改的指定域名
// @Summary 新增/修改域名
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateDomainReq true "没有ID则新增域名，有ID则修改域名"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/domain/domain [post]
func UpdateDomain(c *gin.Context) {
	var domainReq api.UpdateDomainReq
	if err := c.ShouldBind(&domainReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	domain, err := service.Domain().UpdateDomain(&domainReq)
	if err != nil {
		logger.Log().Error("Domain", "新增/修改域名", err)
		c.JSON(500, api.Err("新增/修改域名失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: domain,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// UpdateDomainAss
// @Tags 域名相关
// @title 域名关联服务器
// @description 域名关联服务器
// @Summary 域名关联服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateDomainAssHostReq true "关联传入参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/domain/ass-host [put]
func UpdateDomainAss(c *gin.Context) {
	var domainReq api.UpdateDomainAssHostReq
	if err := c.ShouldBind(&domainReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Domain().UpdateDomainAss(&domainReq)
	if err != nil {
		logger.Log().Error("Domain", "关联服务器", err)
		c.JSON(500, api.Err("关联服务器失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// DeleteDomain
// @Tags 域名相关
// @title 删除域名
// @description 返回success
// @Summary 删除域名
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "删除域名ID切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/domain/domain [delete]
func DeleteDomain(c *gin.Context) {
	var domainReq api.IdsReq
	if err := c.ShouldBind(&domainReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Domain().DeleteDomain(domainReq.Ids)
	if err != nil {
		logger.Log().Error("Domain", "删除域名", err)
		c.JSON(500, api.Err("删除域名失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetDomainAssHost
// @Tags 域名相关
// @title 查询域名对应服务器
// @description 返回服务器切片
// @Summary 查询域名对应服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetPagingMustByIdReq true "传所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/domain/ass-host [get]
func GetDomainAssHost(c *gin.Context) {
	var domainReq api.GetPagingMustByIdReq
	if err := c.ShouldBind(&domainReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	domain, total, err := service.Domain().GetDomainAssHost(&domainReq)
	if err != nil {
		logger.Log().Error("Domain", "查询域名关联服务器", err)
		c.JSON(500, api.Err("查询域名关联服务器失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: domain,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     domainReq.PageInfo.Page,
		PageSize: domainReq.PageInfo.PageSize,
	})
}
