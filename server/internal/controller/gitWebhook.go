package controller

import (
	serviceGitWebhook "fqhWeb/internal/service/gitWebhook"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/api/gitWebhook"
	"fqhWeb/pkg/logger"
	"github.com/gin-gonic/gin"
	"strconv"
)

// HandleGithubWebhook
// @Tags GitWebhook相关
// @title 接收github的Webhook进行处理
// @description 接收github的Webhook进行处理
// @Summary 接收github的Webhook进行处理
// @Produce  application/json
// @Param pid path uint true "项目ID"
// @Param hid path uint true "服务器ID"
// @Success 200 {object} api.Response "{"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/git-gitWebhook/github/{pid}/{hid} [post]
func HandleGithubWebhook(c *gin.Context) {
	pidStr := c.Param("pid")
	var err error
	var pid, hid uint64
	if pid, err = strconv.ParseUint(pidStr, 10, 0); err != nil {
		c.JSON(500, api.ErrorResponse(err))
	}
	hidStr := c.Param("hid")
	if hid, err = strconv.ParseUint(hidStr, 10, 0); err != nil {
		c.JSON(500, api.ErrorResponse(err))
	}
	if err = serviceGitWebhook.GitWebhook().HandleGithubWebhook(uint(pid), uint(hid), c); err != nil {
		logger.Log().Error("GitWebhook", "GithubWebhookPush", err)
		c.JSON(500, api.Err("GithubWebhookPushFail", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// HandleGitlabWebhook
// @Tags GitWebhook相关
// @title 接收gitlab的Webhook进行处理
// @description 接收gitlab的Webhook进行处理
// @Summary 接收gitlab的Webhook进行处理
// @Produce  application/json
// @Param pid path uint true "项目ID"
// @Param hid path uint true "服务器ID"
// @Success 200 {object} api.Response "{"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/git-gitWebhook/gitlab/{pid}/{hid} [post]
func HandleGitlabWebhook(c *gin.Context) {
	pidStr := c.Param("pid")
	var err error
	var pid, hid uint64
	if pid, err = strconv.ParseUint(pidStr, 10, 0); err != nil {
		c.JSON(500, api.ErrorResponse(err))
	}
	hidStr := c.Param("hid")
	if hid, err = strconv.ParseUint(hidStr, 10, 0); err != nil {
		c.JSON(500, api.ErrorResponse(err))
	}
	if err = serviceGitWebhook.GitWebhook().HandleGitlabWebhook(uint(pid), uint(hid), c); err != nil {
		logger.Log().Error("GitWebhook", "GitlabWebhookPush", err)
		c.JSON(500, api.Err("GitlabWebhookPushFail", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// UpdateGitWebhookStatus
// @Tags GitWebhook相关
// @title 更改GitWebhook状态
// @description 更改GitWebhook状态
// @Summary 更改GitWebhook状态
// @Produce  application/json
// @Param CiAuthSign header string true "格式为: 发送机的IP.运维密钥(.不作加密, 两个字符串相连) 再由md5加密"
// @Param data formData gitWebhook.UpdateGitWebhookStatusReq true "填入行ID和状态码"
// @Success 200 {object} api.Response "{"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/git-gitWebhook/project-update-status [patch]
func UpdateGitWebhookStatus(c *gin.Context) {
	// 判断是否运维给的签名
	if err := serviceGitWebhook.GitWebhook().UpdateStatusAuth(c.Request.Header.Get("CiAuthSign"), c.ClientIP()); err != nil {
		c.JSON(403, api.ErrorResponse(err))
		return
	}
	var param gitWebhook.UpdateGitWebhookStatusReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := serviceGitWebhook.GitWebhook().UpdateGitWebhookStatus(param)
	if err != nil {
		logger.Log().Error("GitWebhook", "更改状态码", err)
		c.JSON(500, api.Err("更改状态码失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetGitWebhook
// @Tags GitWebhook相关
// @title 查询GitWebhook记录
// @description 查询GitWebhook记录
// @Summary 查询GitWebhook记录
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.SearchIdStringReq true "有ID输入ID, 否则填入name和页码，全空则全部返回"
// @Success 200 {object} api.Response "{"data":{}, "meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/git-gitWebhook/git-webhook [get]
func GetGitWebhook(c *gin.Context) {
	var param api.SearchIdStringReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	data, total, err := serviceGitWebhook.GitWebhook().GetGitWebhook(param)
	if err != nil {
		logger.Log().Error("GitWebhook", "查询GitWebhook记录", err)
		c.JSON(500, api.Err("查询GitWebhook记录失败", err))
		return
	}
	c.JSON(200, api.PageResult{
		Data: data,
		Meta: api.Meta{
			Msg: "Success",
		},
		Total:    total,
		Page:     param.PageInfo.Page,
		PageSize: param.PageInfo.PageSize,
	})

}

// UpdateGitWebhook
// @Tags GitWebhook相关
// @title 修改GitWebhook记录
// @description 修改GitWebhook记录
// @Summary 修改GitWebhook记录
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body gitWebhook.UpdateGitWebhookReq true "输入修改GitWebhook记录所需参数"
// @Success 200 {object} api.Response "{"data":{}, "meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/git-gitWebhook/git-webhook [put]
func UpdateGitWebhook(c *gin.Context) {
	var param gitWebhook.UpdateGitWebhookReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	data, err := serviceGitWebhook.GitWebhook().UpdateGitWebhook(param)
	if err != nil {
		logger.Log().Error("GitWebhook", "更改GitWebhook记录", err)
		c.JSON(500, api.Err("更改GitWebhook记录失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: data,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// DeleteGitWebhook
// @Tags GitWebhook相关
// @title 删除GitWebhook记录
// @description 删除GitWebhook记录
// @Summary 删除GitWebhook记录
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "输入需要删除的GitWebhook记录的IDs"
// @Success 200 {object} api.Response "{"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/git-gitWebhook/git-webhook [delete]
func DeleteGitWebhook(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	if err := serviceGitWebhook.GitWebhook().DeleteGitWebhook(param.Ids); err != nil {
		logger.Log().Error("GitWebhook", "删除GitWebhook记录", err)
		c.JSON(500, api.Err("删除GitWebhook记录失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}
