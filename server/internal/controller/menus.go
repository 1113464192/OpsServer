package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateMenu
// @Tags 菜单相关
// @title 新增/修改菜单信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改菜单信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateMenuReq true "创建成功，data返回菜单信息"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/menu/menu [post]
func UpdateMenu(c *gin.Context) {
	var menuReq api.UpdateMenuReq
	if err := c.ShouldBind(&menuReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	menu, err := service.Menu().UpdateMenu(&menuReq)
	if err != nil {
		logger.Log().Error("Menu", "创建/修改菜单", err)
		c.JSON(500, api.Err("创建/修改菜单失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: menu,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetMenuList
// @Tags 菜单相关
// @title 获取菜单信息
// @description 返回菜单信息
// @Summary 获取菜单信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.SearchIdStringReq false "所需参数,输入了ids则不再需要输入其他参数；全部留空则全部返回"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/menu/menus [get]
func GetMenuList(c *gin.Context) {
	var param api.SearchIdStringReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	menus, total, err := service.Menu().GetMenuList(param)
	if err != nil {
		logger.Log().Error("Menu", "获取菜单信息", err)
		c.JSON(500, api.Err("获取菜单信息失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data:  menus,
		Total: total,
		Meta: api.Meta{
			Msg: "Success",
		},
		Page:     param.PageInfo.Page,
		PageSize: param.PageInfo.PageSize,
	})
}

// DeleteMenu
// @Tags 菜单相关
// @title 删除菜单
// @description 删除成功返回success
// @Summary 删除菜单
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "菜单ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/menu/menus [delete]
func DeleteMenu(c *gin.Context) {
	var mid api.IdsReq
	if err := c.ShouldBind(&mid); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Menu().DeleteMenu(mid.Ids)
	if err != nil {
		logger.Log().Error("Menu", "删除菜单", err)
		c.JSON(500, api.Err("删除菜单失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}
