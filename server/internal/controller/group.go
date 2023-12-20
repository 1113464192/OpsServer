package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateGroup
// @Tags 用户组相关
// @title 新增/修改用户组信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改用户组信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateGroupReq true "传新增/修改用户组的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/group [post]
func UpdateGroup(c *gin.Context) {
	var groupReq api.UpdateGroupReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, err := service.Group().UpdateGroup(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "创建/修改组失败", err)
		c.JSON(500, api.Err("创建/修改组失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: group,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// UpdateGroupAssUser
// @Tags 用户组相关
// @title 关联用户
// @description 关联用户，关联成功data返回组和关联用户信息
// @Summary 关联用户
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateGroupAssUserReq true "传组ID和对应用户切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/ass-user [put]
func UpdateGroupAssUser(c *gin.Context) {
	var assReq api.UpdateGroupAssUserReq
	if err := c.ShouldBind(&assReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}

	err := service.Group().UpdateGroupAssUser(&assReq)
	if err != nil {
		logger.Log().Error("Group", "更改用户组关联失败", err)
		c.JSON(500, api.Err("更改用户组关联失败", err))
		return
	}

	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// DeleteUserGroup
// @Tags 用户组相关
// @title 删除用户组
// @description 删除成功返回sucess
// @Summary 删除用户组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "传待删除的组ID切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/groups [delete]
func DeleteUserGroup(c *gin.Context) {
	var ids api.IdsReq
	if err := c.ShouldBind(&ids); err != nil {
		c.JSON(500, err)
		return
	}

	err := service.Group().DeleteUserGroup(ids.Ids)
	if err != nil {
		logger.Log().Error("Group", "删除用户组失败", err)
		c.JSON(500, api.Err("删除用户组失败", err))
		return
	}

	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetGroup
// @Tags 用户组相关
// @title 获取用户组
// @description 获取用户组列表(ID直接取用户组无需其他参数，否则需要name和pageinfo)
// @Summary 获取用户组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param id query api.SearchIdStringReq true "所需参数,输入了ids则不再需要输入其他参数；全部留空则全部返回"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/groups [get]
func GetGroup(c *gin.Context) {
	var groupReq api.SearchIdStringReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, total, err := service.Group().GetGroupList(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "获取用户组用户失败", err)
		c.JSON(500, api.Err("获取用户组用户失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: group,
		Meta: api.Meta{
			Msg: "Success",
		},
		Page:     groupReq.PageInfo.Page,
		PageSize: groupReq.PageInfo.PageSize,
		Total:    total,
	})
}

// GetAssUser
// @Tags 用户组相关
// @title 获取组关联用户
// @description 由于swagger本身的限制，get请求的切片会报错，并非接口本身问题，请换个方式，如http://127.0.0.1:9081/api/v1/group/apis?ids=3&ids=4
// @Summary 获取组关联用户
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetPagingMustByIdsReq true "传参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/ass-user [get]
func GetAssUser(c *gin.Context) {
	var groupReq api.GetPagingMustByIdsReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, total, err := service.Group().GetAssUser(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "获取用户组用户失败", err)
		c.JSON(500, api.Err("获取用户组用户失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: group,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     groupReq.PageInfo.Page,
		PageSize: groupReq.PageInfo.PageSize,
	})
}

// GetAssProject
// @Tags 用户组相关
// @title 获取组关联项目
// @description 由于swagger本身的限制，get请求的切片会报错，并非接口本身问题，请换个方式，如http://127.0.0.1:9081/api/v1/group/apis?ids=3&ids=4
// @Summary 获取组关联项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetPagingMustByIdsReq true "传参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/ass-project [get]
func GetAssProject(c *gin.Context) {
	var groupReq api.GetPagingMustByIdsReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, total, err := service.Group().GetAssProject(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "获取用户组项目失败", err)
		c.JSON(500, api.Err("获取用户组项目失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: group,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     groupReq.PageInfo.Page,
		PageSize: groupReq.PageInfo.PageSize,
	})
}

// UpdateGroupAssMenus
// @Tags 用户组相关
// @title 关联菜单
// @description 关联菜单
// @Summary 关联菜单
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateGroupAssMenusReq true "输入菜单IDs和对应用户组ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/ass-menus [put]
func UpdateGroupAssMenus(c *gin.Context) {
	var assReq api.UpdateGroupAssMenusReq
	if err := c.ShouldBind(&assReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	data, err := service.Menu().UpdateGroupAssMenus(&assReq)
	if err != nil {
		logger.Log().Error("Group", "菜单与用户组关联失败", err)
		c.JSON(500, api.Err("菜单与用户组关联失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: data,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetGroupAssMenus
// @Tags 用户组相关
// @title 获取组关联菜单
// @description 由于swagger本身的限制，get请求的切片会报错，并非接口本身问题，请换个方式，如http://127.0.0.1:9081/api/v1/group/apis?ids=3&ids=4
// @Summary 获取组关联菜单
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetPagingMustByIdsReq true "传组id，空页码则全部返回"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/ass-menus [get]
func GetGroupAssMenus(c *gin.Context) {
	var groupReq api.GetPagingMustByIdsReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	menu, total, err := service.Group().GetGroupAssMenus(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "获取用户组菜单失败", err)
		c.JSON(500, api.Err("获取用户组项目菜单", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: menu,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     groupReq.PageInfo.Page,
		PageSize: groupReq.PageInfo.PageSize,
	})
}

// GetCasbinList
// @Tags 用户组相关
// @title 获取用户组的API权限列表
// @description 由于swagger本身的限制，get请求的切片会报错，并非接口本身问题，请换个方式，如http://127.0.0.1:9081/api/v1/group/apis?ids=3&ids=4
// @Summary 获取用户组的API权限列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param data query api.IdsReq true "传组ID切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router    /api/v1/group/apis [get]
func GetCasbinList(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	paths, err := service.CasbinServiceApp().GetPolicyPathByGroupIds(param.Ids)
	if err != nil {
		logger.Log().Error("Group", "获取用户API列表失败", err)
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
// @Tags 用户组相关
// @title 为用户组分配API权限
// @description 为用户组分配API权限
// @Summary 为用户组分配API权限
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param     data  body      api.UpdateCasbinReq  true  "传为用户组分配的API切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router    /api/v1/group/apis [put]
func UpdateCasbin(c *gin.Context) {
	var param api.UpdateCasbinReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.CasbinServiceApp().UpdateCasbin(param)
	if err != nil {
		logger.Log().Error("API", "分配API权限错误失败", err)
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
