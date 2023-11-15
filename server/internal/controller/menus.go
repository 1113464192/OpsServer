package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils/jwt"

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
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/menu/update [post]
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

// UpdateMenuAss
// @Tags 菜单相关
// @title 更改菜单与用户组关联信息
// @description 关联用户组，关联成功data返回组和关联用户信息
// @Summary 关联用户组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateMenuAssReq true "输入菜单ID和对应用户组IDs"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/menu/association [put]
func UpdateMenuAss(c *gin.Context) {
	var menuReq api.UpdateMenuAssReq
	if err := c.ShouldBind(&menuReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	menu, err := service.Menu().UpdateMenuAss(&menuReq)
	if err != nil {
		logger.Log().Error("Menu", "菜单与用户组关联", err)
		c.JSON(500, api.Err("菜单与用户组关联失败", err))
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
// @title 获取用户组对应菜单信息
// @description 返回关联的菜单
// @Summary 获取用户组对应菜单信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetMenuListReq false "用户组ID，获取菜单和关联用户组,不输入返回所有菜单"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/menu/getMenus [get]
func GetMenuList(c *gin.Context) {
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
	}
	var gid api.GetMenuListReq
	if err := c.ShouldBind(&gid); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	isAdmin := claims.User.IsAdmin
	menu, err := service.Menu().GetMenuList(gid.Id, isAdmin)
	if err != nil {
		logger.Log().Error("Menu", "获取菜单对应用户组", err)
		c.JSON(500, api.Err("获取菜单对应用户组失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: menu,
		Meta: api.Meta{
			Msg: "Success",
		},
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
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/menu/delete [delete]
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
