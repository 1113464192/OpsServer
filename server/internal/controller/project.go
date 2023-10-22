package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils/jwt"

	"github.com/gin-gonic/gin"
)

// UpdateProject
// @Tags 项目相关
// @title 新增/修改项目
// @description 返回新增/修改的指定项目
// @Summary 新增/修改项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateProjectReq true "更新project所需参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/update [post]
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

// UpdateHostAss
// @Tags 项目相关
// @title 关联服务器
// @description 服务器ID[多选]
// @Summary 关联服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateProjectAssHostReq true "关联传入参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/association [put]
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
// @Param data body api.GetProjectReq false "输入项目ID，获取项目,不输入返回所有项目"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/getProject [post]
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

// GetSelfProjectList
// @Tags 项目相关
// @title 获取自身所属项目
// @description 返回自身所属项目
// @Summary 获取自身所属项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.PageInfo true "页码"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/getSelfProject [get]
func GetSelfProjectList(c *gin.Context) {
	var err error
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
		return
	}
	var pageReq api.PageInfo
	if err := c.ShouldBind(&pageReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	groupList := claims.User.UserGroups
	project, total, err := service.Project().GetSelfProjectList(&groupList, &pageReq)
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
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
	})
}

// GetHostAss
// @Tags 项目相关
// @title 获取项目对应服务器
// @description 返回服务器切片
// @Summary 获取项目对应的服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.GetHostAssReq true "获取关联host的参数"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/project/getHost [post]
func GetHostAss(c *gin.Context) {
	var projectReq api.GetHostAssReq
	if err := c.ShouldBind(&projectReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	// 后面补total
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
