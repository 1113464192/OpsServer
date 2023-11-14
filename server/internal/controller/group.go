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
// @Param data formData api.UpdateGroupReq true "创建成功，data返回用户组信息"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/update [post]
func UpdateGroup(c *gin.Context) {
	var groupReq api.UpdateGroupReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, err := service.Group().UpdateGroup(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "创建/修改组", err)
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

// UpdateUserAss
// @Tags 用户组相关
// @title 关联用户
// @description 关联用户，关联成功data返回组和关联用户信息
// @Summary 关联用户
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateUserAssReq true "输入组ID和对应用户IDs"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/association [put]
func UpdateUserAss(c *gin.Context) {
	var assReq api.UpdateUserAssReq
	if err := c.ShouldBind(&assReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}

	err := service.Group().UpdateUserAss(&assReq)
	if err != nil {
		logger.Log().Error("Group", "更改用户组关联", err)
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
// @Param ids body api.IdsReq true "要删除的组ID"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/delete [delete]
func DeleteUserGroup(c *gin.Context) {
	var ids api.IdsReq
	if err := c.ShouldBind(&ids); err != nil {
		c.JSON(500, err)
		return
	}

	err := service.Group().DeleteUserGroup(ids.Ids)
	if err != nil {
		logger.Log().Error("Group", "删除用户组", err)
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
// @description 返回指定用户组/不传Name返回所有用户组
// @Summary 获取用户组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param id query api.GetGroupReq true "输入组名，不输入则全部返回"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/getGroups [get]
func GetGroup(c *gin.Context) {
	var groupReq api.GetGroupReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, total, err := service.Group().GetGroupList(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "获取用户组用户", err)
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
// @description 返回组关联的用户
// @Summary 获取组关联用户
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetPagingByIdReq true "传参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/getUserAss [get]
func GetAssUser(c *gin.Context) {
	var groupReq api.GetPagingByIdReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, total, err := service.Group().GetAssUser(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "获取用户组用户", err)
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
// @description 返回组关联的项目
// @Summary 获取组关联项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetPagingByIdReq true "传参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/group/getProjectAss [get]
func GetAssProject(c *gin.Context) {
	var groupReq api.GetPagingByIdReq
	if err := c.ShouldBind(&groupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	group, total, err := service.Group().GetAssProject(&groupReq)
	if err != nil {
		logger.Log().Error("Group", "获取用户组项目", err)
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
