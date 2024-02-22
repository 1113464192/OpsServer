package middleware

import (
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	jwtService = service.Jwt()
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(403, api.Err("用户访问令牌缺失", nil))
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(403, api.Err("用户访问令牌格式有误", nil))
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		// 判断是否在黑名单
		if jwtService.IsBlacklist(parts[1]) {
			c.JSON(403, api.Err("您的账户异地登陆或令牌失效", nil))
			c.Abort()
			return
		}
		// 判断token是否已过期
		claims, err := auth.ParseToken(parts[1])
		// mc 里面包含对应登录账号，签发人（Issuer） 过期时间（ExpiresAt）
		if err != nil {
			//已过期，把token拉到黑名单
			err2 := jwtService.JwtAddBlacklist(&model.JwtBlacklist{Jwt: parts[1]})
			if err2 != nil {
				logger.Log().Error("JWT", "拉取token到黑名单失败", err2)
				c.JSON(500, api.Err("拉取token到黑名单失败", err2))
				c.Abort()
				return
			}
			logger.Log().Info("JWT", "授权已过期", err)
			c.JSON(403, api.Err("授权已过期", err))
			c.Abort()
			return
		}
		// 将当前请求的claims信息保存到请求的上下文c上
		c.Set("claims", claims)
		c.Next() // 后续的处理函数可以用过c.Get("claims")来获取当前请求的用户信息
	}
}
