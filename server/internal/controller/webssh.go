package controller

import (
	"fqhWeb/internal/service/globalFunc"
	"fqhWeb/internal/service/webssh"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util/jwt"

	"github.com/gin-gonic/gin"
)

// WebsshConn
// @Tags Webssh相关
// @title 连接Webssh
// @description 连接Webssh,自动获取当前用户，防止冒用其它user
// @Summary 连接Webssh
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.WebsshConnReq true "传HostID、屏幕高宽"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/webssh/webssh-conn [get]
func WebsshConn(c *gin.Context) {
	var (
		param api.WebsshConnReq
		err   error
	)
	if err = globalFunc.IncreaseWebSSHConn(); err != nil {
		c.JSON(500, api.Err("已达到最大webssh数量", err))
		return
	}

	if err = c.ShouldBindQuery(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}

	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
		return
	}

	wsRes, err := webssh.WebSsh().WebSshHandle(c, &claims.User, param)
	if err != nil {
		logger.Log().Error("Webssh", wsRes+"连接Webssh失败", err)
		c.JSON(500, api.Err(wsRes+"连接Webssh失败", err))
		return
	}
}
