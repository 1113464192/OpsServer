package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Test
// @Tags 测试相关
// @title 测试Gin能否正常访问
// @description 无设置权限，返回"Hello world!~~(无权限版)"
// @Summary 测试Gin能否正常访问
// @Produce  application/json
// @Success 200 {} string "{"message": "Hello world!~~(无权限版)"}"
// @Router /api/v1/ping [get]
func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello world!~~(无权限版)",
	})
}

// Test2
// @Tags 测试相关
// @title 测试Gin通过认证能否正常访问
// @description 有设置权限，返回"Hello world!~~(有权限版)"
// @Summary 测试Gin通过认证能否正常访问
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Success 200 {} string "{"message": "Hello world!~~(有权限版)"}"
// @Router /api/v1/ping2 [get]
func Test2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello world!~~(有权限版)",
	})
}
