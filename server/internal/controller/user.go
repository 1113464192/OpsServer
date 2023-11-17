package controller

import (
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils"
	"fqhWeb/pkg/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserLogin
// @Tags 用户相关
// @title 用户登录
// @description 用户名长度不少于3位，密码不少于6位
// @Summary 用户登录
// @Produce  application/json
// @Param data formData api.AuthLoginReq true "用户名, 密码"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/login [post]
func UserLogin(c *gin.Context) {
	var loginReq api.AuthLoginReq
	if err := c.ShouldBind(&loginReq); err == nil {
		u := &model.User{Username: loginReq.Username, Password: loginReq.Password}
		if userInfo, err := service.User().Login(u); err != nil {
			logger.Log().Error("User", "账号或密码错误", err)
			c.JSON(403, api.Err("账号或密码错误", err))
			return
		} else {
			c.JSON(200, api.Response{
				Data: userInfo,
				Meta: api.Meta{
					Msg: "登录成功",
				},
			})
			return
		}
	} else {
		c.JSON(200, api.ErrorResponse(err))
	}
}

// UserLogin
// @Tags 用户相关
// @title 用户登出
// @description 登出 - 把JWT拉入黑名单
// @Summary 登出
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/logout [put]
func UserLogout(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	token := parts[1]
	jwt := &model.JwtBlacklist{Jwt: token}
	if err := service.Jwt().JwtAddBlacklist(jwt); err != nil {
		logger.Log().Error("User", "用户登出失败，jwt没有拉入黑名单", err)
		c.JSON(500, api.Err("用户登出失败，jwt没有拉入黑名单", err))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "登出成功",
			},
		})
		return
	}
}

// UpdateUser
// @Tags 用户相关
// @title 新增/修改用户信息
// @description 新增不用传用户ID，修改才传用户ID
// @Summary 新增/修改用户信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateUserReq true "创建成功，data返回密码"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/update [post]
func UpdateUser(c *gin.Context) {
	var userReq api.UpdateUserReq
	if err := c.ShouldBind(&userReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	user, passwd, err := service.User().UpdateUser(&userReq)
	if err != nil {
		logger.Log().Error("User", "添加/修改用户", err)
		if err.Error() == "用户密码bcrypt加密失败" {
			c.JSON(500, api.Err("用户密码bcrypt加密失败", nil))
			return
		}
		c.JSON(500, api.Err("添加/修改用户", err))
		return
	}

	data := map[string]any{
		"string": user,
	}
	if passwd != "" {
		data["passwd"] = passwd
	}

	c.JSON(200, api.Response{
		Data: data,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetUserList
// @Tags 用户相关
// @title 用户列表
// @description 获取用户列表(IDs直接取用户无需其他参数，否则需要name和pageinfo)
// @Summary 获取用户列表
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetUserListReq true "所需参数,输入了ids则不再需要输入其他参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/search [get]
func GetUserList(c *gin.Context) {
	var param api.GetUserListReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}

	user, total, err := service.User().GetUserList(param)
	if err != nil {
		logger.Log().Error("User", "获取用户列表", err)
		c.JSON(500, api.Err("获取失败", err))
		return
	} else {
		c.JSON(200, api.PageResult{
			Meta: api.Meta{
				Msg: "Success",
			},
			Data:     user,
			Page:     param.Page,
			PageSize: param.PageSize,
			Total:    total,
		})
	}
}

// DeleteUser
// @Tags 用户相关
// @title 删除用户
// @description 删除用户
// @Summary 删除用户
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "用户IDs"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/delete [delete]
func DeleteUser(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	if err := service.User().DeleteUser(param.Ids); err != nil {
		logger.Log().Error("User", "删除用户", err)
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

// UpdatePasswd
// @Tags 用户相关
// @title 修改用户密码
// @description 修改用户密码
// @Summary 修改用户密码
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.PasswordReq true "创建成功，data返回密码"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/password [patch]
func UpdatePasswd(c *gin.Context) {
	var passwd api.PasswordReq
	var err error
	if err := c.ShouldBind(&passwd); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	passwd.Password, err = utils.GenerateFromPassword(passwd.Password)
	if err != nil {
		logger.Log().Error("User", "用户密码加密", err)
		c.JSON(500, api.Err("密码加密失败", err))
		return
	}

	if err := service.User().UpdatePasswd(&passwd); err != nil {
		logger.Log().Error("User", "更新密码", err)
		c.JSON(500, api.Err("更新密码失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}

// UpdateSelfPasswd
// @Tags 用户相关
// @title 修改用户自己的密码
// @description 修改用户自己的密码
// @Summary 修改用户自己的密码
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param password formData string true "需要修改的密码"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/selfPassword [patch]
func UpdateSelfPasswd(c *gin.Context) {
	var passwd api.PasswordReq
	var err error
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
		return
	}

	passwd.ID = claims.User.ID
	passwd.Password = c.PostForm("password")
	if err := c.ShouldBind(&passwd); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}

	passwd.Password, err = utils.GenerateFromPassword(passwd.Password)
	if err != nil {
		logger.Log().Error("User", "用户密码加密", err)
		c.JSON(500, api.Err("密码加密失败", err))
		return
	}

	if err := service.User().UpdatePasswd(&passwd); err != nil {
		logger.Log().Error("User", "更新密码", err)
		c.JSON(500, api.Err("更新密码失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}

// UpdateStatus
// @Tags 用户相关
// @title 修改用户状态
// @description status参数：恢复用户传1，禁用用户传2
// @Summary 修改用户状态
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.StatusReq true "用户状态(恢复传1，禁用传2)"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/status [patch]
func UpdateStatus(c *gin.Context) {
	var params api.StatusReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	if err := service.User().UpdateStatus(&params); err != nil {
		logger.Log().Error("User", "更改用户状态", err)
		c.JSON(500, api.Err("更改用户状态失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}

}

// GetSelfInfo
// @Tags 用户相关
// @title 用户个人信息
// @description 获取用户个人信息
// @Summary 用户个人信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/getSelfInfo [get]
func GetSelfInfo(c *gin.Context) {
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
		return
	}
	userInfo, err := service.User().GetSelfInfo(&claims.User.ID)
	if err != nil {
		logger.Log().Error("User", "获取用户个人信息", err)
		c.JSON(500, api.Err("获取用户个人信息失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Data: userInfo,
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}

// GetAssGroup
// @Tags 用户相关
// @title 获取关联组
// @description 返回用户关联的组
// @Summary 获取关联组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.IdReq true "用户的ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/getAssGroup [get]
func GetAssGroup(c *gin.Context) {
	var id api.IdReq
	if err := c.ShouldBind(&id); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	groupInfo, err := service.User().GetAssGroup(id.Id)
	if err != nil {
		logger.Log().Error("User", "获取用户关联组信息失败", err)
		c.JSON(500, api.Err("获取用户个人信息失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Data: groupInfo,
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}

// GetSelfAssGroup
// @Tags 用户相关
// @title 获取用户本身关联的组
// @description 返回用户本身关联的组
// @Summary 获取用户本身关联的组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/getSelfAssGroup [get]
func GetSelfAssGroup(c *gin.Context) {
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
		return
	}
	id := claims.User.ID
	groupInfo, err := service.User().GetAssGroup(id)
	if err != nil {
		logger.Log().Error("User", "获取用户关联组信息失败", err)
		c.JSON(500, api.Err("获取用户个人信息失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Data: groupInfo,
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}

// GetRecordList
// @Tags 用户相关
// @title 获取用户操作记录
// @description 获取用户操作记录，不包含get
// @Summary 获取用户操作记录
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetPagingByIdReq true "用户username"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/getActLog [get]
func GetRecordList(c *gin.Context) {
	var param api.GetPagingByIdReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	logs, total, err := service.Record().GetRecordList(&param)
	if err != nil {
		logger.Log().Error("User", "获取用户操作记录失败", err)
		c.JSON(500, api.Err("获取用户操作记录失败", err))
		return
	} else {
		c.JSON(200, api.PageResult{
			Meta: api.Meta{
				Msg: "Success",
			},
			Data:     logs,
			Total:    total,
			Page:     param.PageInfo.Page,
			PageSize: param.PageInfo.PageSize,
		})
	}
}

// UpdateKeyFileContext
// @Tags 用户相关
// @title 提交自身私钥文件
// @description 是私钥不要提交公钥！私钥如: id_rsa
// @Summary 提交自身私钥文件
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param keyFile formData file true "私钥文件上传"
// @Param Passphrase formData string true "私钥通行证密码上传"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/keyFile [post]
func UpdateKeyFileContext(c *gin.Context) {
	file, err := c.FormFile("keyFile")
	if err != nil {
		c.JSON(500, api.Err("上传失败", err))
		return
	}
	passphrase := c.PostForm("Passphrase")
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
		return
	}
	err = service.User().UpdateKeyFileContext(file, passphrase, claims.User.ID)
	if err != nil {
		logger.Log().Error("User", "上传文件写入个人密钥失败", err)
		c.JSON(500, api.Err("上传文件写入个人密钥失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}

// UpdateKeyContext
// @Tags 用户相关
// @title 提交自身私钥字符串
// @description 是私钥字符串不要提交公钥文件！私钥如: id_rsa的内容
// @Summary 提交自身私钥字符串
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param keyStr formData string true "私钥文本内容上传"
// @Param Passphrase formData string true "私钥通行证密码上传"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/user/keyStr [post]
func UpdateKeyContext(c *gin.Context) {
	keyStr := c.PostForm("keyStr")
	passphrase := c.PostForm("Passphrase")
	cClaims, _ := c.Get("claims")
	claims, ok := cClaims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(401, api.Err("token携带的claims不合法", nil))
		c.Abort()
		return
	}
	err := service.User().UpdateKeyContext(keyStr, passphrase, claims.User.ID)
	if err != nil {
		logger.Log().Error("User", "私钥字符串写入个人密钥失败", err)
		c.JSON(500, api.Err("私钥字符串写入个人密钥失败", err))
		return
	} else {
		c.JSON(200, api.Response{
			Meta: api.Meta{
				Msg: "Success",
			},
		})
	}
}
