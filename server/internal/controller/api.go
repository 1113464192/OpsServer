package controller

import (
	"fmt"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GetApiList
// @Tags Api相关
// @title 获取Api列表
// @description 获取Api列表 可分页
// @Summary 获取Api列表
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param     data  formData      api.PageInfo   false  "页码, 每页大小"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/api/getApiList [post]
func GetApiList(c *gin.Context) {
	var param api.PageInfo
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	apiList, total, err := service.Api().GetApiList(param)
	if err != nil {
		logger.Log().Error("API", "获取失败", err)
		c.JSON(500, api.Err("获取失败", nil))
		return
	} else {
		c.JSON(200, api.PageResult{
			Meta: api.Meta{
				Msg: "Success",
			},
			Data:     apiList,
			Page:     param.Page,
			PageSize: param.PageSize,
			Total:    total,
		})
	}
}

// UpdateApi
// @Tags Api相关
// @title 新增或者修改Api
// @description 新增或者修改Api
// @Summary 新增或者修改Api
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param     data  formData      api.UpdateApiReq  true  "新增或者修改Api"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router    /api/v1/api/updateApi [post]
func UpdateApi(c *gin.Context) {
	var param api.UpdateApiReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	apiRes, err := service.Api().UpdateApi(&param)
	if err != nil {
		logger.Log().Error("API", "Api新增或修改错误", err)
		c.JSON(500, api.Err("操作失败", nil))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
			Data: fmt.Sprintf("ID: %d", apiRes.ID),
		})
	}
}

// DeleteApi
// @Tags Api相关
// @title 删除API
// @description 删除API
// @Summary 删除API
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param     data  body      api.IdsReq   true  "id"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/api/delApi [delete]
func DeleteApi(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	if err := service.Api().DeleteApi(param.Ids); err != nil {
		logger.Log().Error("API", "删除失败", err)
		c.JSON(500, api.Err("删除失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}

// FreshCasbin
// @Tags Api相关
// @title 刷新casbin缓存
// @description 刷新casbin缓存
// @Summary 刷新casbin缓存
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/api/fresh [get]
func FreshCasbin(c *gin.Context) {
	err := service.Api().FreshCasbin()
	if err != nil {
		logger.Log().Error("API", "刷新casbin失败", err)
		c.JSON(500, api.Err("刷新失败", nil))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCasbinList
// @Tags Api相关
// @title 获取用户已有的API权限列表
// @description 获取用户已有的API权限列表
// @Summary 获取用户已有的API权限列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param gid query []uint true "组ID"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router    /api/v1/api/casbinList [get]
func GetCasbinList(c *gin.Context) {
	idStr := c.Query("gid")
	paths, err := service.CasbinServiceApp().GetPolicyPathByGroupId(idStr)
	if err != nil {
		logger.Log().Error("API", "获取用户API列表失败", err)
		c.JSON(500, api.Err("获取用户API列表失败", nil))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
		Data: paths,
	})
}

// UpdateCasbin
// @Tags Api相关
// @title 为用户分配API权限
// @description 为用户分配API权限
// @Summary 为用户分配API权限
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param     data  body      api.CasbinInReceiveReq  true  "为用户分配API权限的请求"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router    /api/v1/api/updateCasbin [post]
func UpdateCasbin(c *gin.Context) {
	var param api.CasbinInReceiveReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.CasbinServiceApp().UpdateCasbin(param.GroupId, param.Ids)
	if err != nil {
		logger.Log().Error("API", "分配API权限错误", err)
		c.JSON(500, api.Err("分配API权限错误", nil))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}
