package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateProject
// @Tags 项目相关
// @title 新增/修改项目
// @description 返回新增/修改的指定项目
// @Summary 新增/修改项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateProjectReq true "新增/修改project所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/project [post]
func UpdateProject(c *gin.Context) {
	var projectReq api.UpdateProjectReq
	if err := c.ShouldBind(&projectReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	project, err := service.Project().UpdateProject(&projectReq)
	if err != nil {
		logger.Log().Error("Project", "创建/修改项目", err)
		c.JSON(500, api.Err("创建/修改项目失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: project,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// DeleteProject
// @Tags 项目相关
// @title 删除项目
// @description 返回success
// @Summary 删除项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "删除project的IDs切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/project [delete]
func DeleteProject(c *gin.Context) {
	var projectReq api.IdsReq
	if err := c.ShouldBind(&projectReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Project().DeleteProject(projectReq.Ids)
	if err != nil {
		logger.Log().Error("Project", "删除项目", err)
		c.JSON(500, api.Err("删除项目失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// UpdateHostAss
// @Tags 项目相关
// @title 关联服务器
// @description 服务器ID[多选]
// @Summary 关联服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateProjectAssHostReq true "关联传入参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/ass-host [put]
func UpdateHostAss(c *gin.Context) {
	var ProjectReq api.UpdateProjectAssHostReq
	if err := c.ShouldBind(&ProjectReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Project().UpdateHostAss(&ProjectReq)
	if err != nil {
		logger.Log().Error("Project", "关联服务器", err)
		c.JSON(500, api.Err("关联服务器失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetProject
// @Tags 项目相关
// @title 获取项目
// @description 返回指定项目
// @Summary 获取项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetProjectReq false "输入项目ID，获取项目,不输入返回所有项目"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/project [get]
func GetProject(c *gin.Context) {
	var projectReq api.GetProjectReq
	if err := c.ShouldBind(&projectReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	project, total, err := service.Project().GetProject(&projectReq)
	if err != nil {
		logger.Log().Error("Project", "获取项目", err)
		c.JSON(500, api.Err("获取项目失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: project,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     projectReq.PageInfo.Page,
		PageSize: projectReq.PageInfo.PageSize,
	})
}

// GetHostAss
// @Tags 项目相关
// @title 获取项目对应服务器
// @description 返回服务器切片
// @Summary 获取项目对应的服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetHostAssReq true "获取关联host的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/ass-host [get]
func GetHostAss(c *gin.Context) {
	var projectReq api.GetHostAssReq
	if err := c.ShouldBind(&projectReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}

	hostList, total, err := service.Project().GetHostAss(&projectReq)
	if err != nil {
		logger.Log().Error("Project", "获取项目拥有机器", err)
		c.JSON(500, api.Err("获取项目拥有机器失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: hostList,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     projectReq.PageInfo.Page,
		PageSize: projectReq.PageInfo.PageSize,
	})
}
