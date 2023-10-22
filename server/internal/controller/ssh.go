package controller

// TestSSH
// @Tags SSH相关
// @title SSH Forward测试
// @description
// @Summary SSH Forward测试
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.clientConfig true "创建成功，data返回密码"
// @Success 200 {} string "{"data":{},"meta":{msg":"Success"}}"
// @Failure 500 {string} string "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/ssh/testSSH [post]
// func TestSSH(c *gin.Context) {
// 	var params *ssh.ClientConfigService
// 	err := c.ShouldBind()
// 	// ①先从项目中获取操作的机器和命令模板
// 	// 不定长参数接收参数

// 	fmt.Println("\n ①先从项目中获取操作的机机器 \n")
// 	clientConfig := &sshService.ClientConfigService{
// 		Host:      ,
// 		Port:      22,
// 		Username:  "root",
// 		Password:  clientPassword,
// 		Key:       clientKey,
// 		KeyPasswd: clientKeyPasswd,
// 	}

// 	// ②走工单审批
// 	fmt.Println("\n ②走工单审批 \n")

// 	// ③执行操作

// }
