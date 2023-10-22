package middleware

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/utils/jwt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var casbinService = service.CasbinServiceApp()

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取请求的URI
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		cClaims, isExist := c.Get("claims")
		if !isExist {
			c.JSON(401, api.Err("未获取到token携带的claims", nil))
			c.Abort()
			return
		}
		claims, ok := cClaims.(*jwt.CustomClaims)
		if !ok {
			c.JSON(401, api.Err("token携带的claims不合法", nil))
			c.Abort()
			return
		}
		userGroup := claims.User.UserGroups
		sCount := 0
		var sub string
		// 遍历用户对应的所有组
		for _, group := range userGroup {
			if group.ID == 1 {
				// 管理员
				sub = "admin"
			} else {
				sub = strconv.FormatUint(uint64(group.ID), 10)
			}
			e := casbinService.Casbin()
			success, _ := e.Enforce(sub, obj, act)
			if success {
				sCount += 1
			}
		}

		// 判断策略中是否存在
		if sCount > 0 {
			c.Next()
		} else {
			c.JSON(403, api.Err("权限不足", nil))
			c.Abort()
			return
		}
	}
}
