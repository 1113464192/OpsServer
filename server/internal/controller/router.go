package controller

import (
	"fqhWeb/internal/middleware"
	_ "fqhWeb/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoute() *gin.Engine {
	r := gin.Default()
	// if configs.Conf.System.Mode == "product" {
	// gin.SetMode(gin.ReleaseMode)
	// swagger.SwaggerInfo.Host = "127.0.0.1:9080"
	// }
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.Cors())
	// ---------API版本区分----------
	v1 := r.Group("/api/v1")
	v1.GET("ping", Test)
	// ---------登录----------
	v1.POST("login", UserLogin)
	// ---------Git-Webhook相关----------
	gitWebhookRouter := v1.Group("git-gitWebhook")
	{
		gitWebhookRouter.POST("github/:pid/:hid", HandleGithubWebhook)          // 接收github的webhook做处理
		gitWebhookRouter.POST("gitlab/:pid/:hid", HandleGitlabWebhook)          // 接收gitlab的webhook做处理
		gitWebhookRouter.PATCH("project-update-status", UpdateGitWebhookStatus) // 更改git-webhook记录的状态码
	}
	// ------------验证相关------------
	v1.Use(middleware.JWTAuthMiddleware()).Use(middleware.CasbinHandler()).Use(middleware.UserActionRecord())
	{
		// -------------接口权限测试--------------
		v1.GET("ping2", Test2)

		// ------------API Casbin RBAC相关----------------
		apiRouter := v1.Group("api")
		{
			apiRouter.GET("apis", GetApiList)    // 获取Api列表
			apiRouter.POST("fresh", FreshCasbin) // 刷新casbin缓存
			apiRouter.POST("api", UpdateApi)     // api添加/修改
			apiRouter.DELETE("api", DeleteApi)   // 删除api
		}
		// ------------用户相关------------
		userRoute := v1.Group("user")
		{
			userRoute.POST("user", UpdateUser)                       // 新增/修改用户
			userRoute.GET("users", GetUserList)                      // 查询用户切片
			userRoute.DELETE("users", DeleteUser)                    // 删除用户
			userRoute.PATCH("status", UpdateStatus)                  // 修改用户状态
			userRoute.PATCH("password", UpdatePasswd)                // 修改用户密码
			userRoute.PATCH("self-password", UpdateSelfPasswd)       // 修改自己的密码
			userRoute.GET("self-user", GetSelfInfo)                  // 查询自己的信息
			userRoute.GET("ass-group", GetAssGroup)                  // 根据用户ID查询组
			userRoute.GET("self-ass-group", GetSelfAssGroup)         // 获取自己关联的组
			userRoute.PUT("logout", UserLogout)                      // 登出
			userRoute.GET("action-log", GetRecordList)               // 查询用户所有的历史操作
			userRoute.GET("action-log/exist-date", GetRecordLogDate) // 查询有多少个月份表可供查询
			userRoute.POST("key-file", UpdateKeyFileContext)         // 添加私钥文件
			userRoute.POST("key-str", UpdateKeyContext)              // 添加私钥内容(与文件二选一)
		}
		// ------------用户组相关--------------
		groupRoute := v1.Group("group")
		{
			groupRoute.POST("group", UpdateGroup)            // 新增/修改组
			groupRoute.PUT("ass-user", UpdateGroupAssUser)   // 用户组关联用户
			groupRoute.DELETE("groups", DeleteUserGroup)     // 删除用户组
			groupRoute.GET("groups", GetGroup)               // 查询组切片
			groupRoute.GET("ass-user", GetAssUser)           // 查询组关联的用户
			groupRoute.GET("ass-project", GetAssProject)     // 查询组关联的项目
			groupRoute.PUT("ass-menus", UpdateGroupAssMenus) // 用户组关联菜单
			groupRoute.GET("ass-menus", GetGroupAssMenus)    // 查询组关联的菜单
			groupRoute.GET("apis", GetCasbinList)            // 获取Group的API权限列表
			groupRoute.PUT("apis", UpdateCasbin)             // 为用户组分配API权限
		}
		// ------------菜单相关--------------
		menuRoute := v1.Group("menu")
		{
			menuRoute.POST("menu", UpdateMenu)    // 新增/修改组
			menuRoute.GET("menus", GetMenuList)   // 获取菜单信息
			menuRoute.DELETE("menus", DeleteMenu) // 删除菜单
		}
		// -----------项目相关-------------
		projectRoute := v1.Group("project")
		{
			projectRoute.POST("project", UpdateProject)         // 新增/修改项目
			projectRoute.GET("project", GetProject)             // 查询项目
			projectRoute.GET("ass-host", GetHostAss)            // 查询项目关联的服务器
			projectRoute.PUT("ass-host", UpdateHostAss)         // 项目关联服务器
			projectRoute.DELETE("project", DeleteProject)       // 删除项目
			projectRoute.GET("cloud-type-set", GetCloudTypeSet) // 获取可填写的云商
		}
		// --------------云相关---------------
		cloudRoute := v1.Group("cloud")
		{
			cloudRoute.POST("vpc", CreateCloudVpc)                               // 创建项目的VPC
			cloudRoute.POST("vpc-subnet", CreateCloudVpcSubnet)                  // 创建项目VPC的subnet
			cloudRoute.POST("security-group", CreateCloudSecurityGroup)          // 创建项目的安全组
			cloudRoute.GET("host-in-subnet-sum", GetCloudHostInVpcSubnetSum)     // 项目机器沾指定subnet的总数
			cloudRoute.PUT("instance-config", UpdateCloudInstanceConfig)         // 创建/更新云项目的实例配置
			cloudRoute.GET("instance-config", GetCloudInstanceConfig)            // 获取云项目的实例配置
			cloudRoute.GET("instance", GetCloudInsInfo)                          // 获取云实例的详细信息
			cloudRoute.GET("vpc-id", GetCloudVpcId)                              // 根据VPCname获取VPC-ID
			cloudRoute.GET("vpc-subnet-id", GetCloudVpcSubnetId)                 // 根据subnetName和VPCID获取subnet-id
			cloudRoute.GET("security-group-id", GetCloudSecurityGroupId)         // 根据安全组Name获取安全组ID
			cloudRoute.GET("project-id", GetCloudProjectId)                      // 根据项目名获取云项目的ID(非数据库项目ID)
			cloudRoute.GET("GetCloudInstanceTypeList", GetCloudInstanceTypeList) // 获取对应的云实例类型列表
			cloudRoute.POST("create-instance", CreateCloudInstance)              // 创建云主机
			cloudRoute.POST("return-instance", ReturnCloudInstance)              // 退还云主机
		}

		// --------------Git-Webhook相关---------------------
		{
			gitWebhookRouter.PUT("git-webhook", UpdateGitWebhook)    // 修改git-webhook记录的数据
			gitWebhookRouter.GET("git-webhook", GetGitWebhook)       // 查询git-webhook记录
			gitWebhookRouter.DELETE("git-webhook", DeleteGitWebhook) // 删除git-webhook记录
		}
		// -----------主机相关-------------
		hostRoute := v1.Group("host")
		{
			hostRoute.POST("host", UpdateHost)       // 新增/修改服务器
			hostRoute.GET("password", GetHostPasswd) // 获取服务器密码
			hostRoute.DELETE("host", DeleteHost)     // 删除主机
			hostRoute.GET("host", GetHost)           // 获取主机信息
		}
		// -----------域名相关--------------
		domainRoute := v1.Group("domain")
		{
			domainRoute.POST("domain", UpdateDomain)      // 新增/修改的域名
			domainRoute.PUT("ass-host", UpdateDomainAss)  // 更新域名关联的服务器
			domainRoute.DELETE("domain", DeleteDomain)    // 删除域名
			domainRoute.GET("ass-host", GetDomainAssHost) // 获取域名关联服务器
		}
		// -----------任务模板相关-----------
		taskRoute := v1.Group("template")
		{
			taskRoute.POST("template", UpdateTemplate)       // 新增/修改模板
			taskRoute.PUT("ass-host", UpdateTemplateAssHost) // 模板关联服务器
			taskRoute.GET("template", GetProjectTemplate)    // 获取模板
			taskRoute.DELETE("template", DeleteTemplate)     // 删除模板
			taskRoute.GET("condition-set", GetConditionSet)  // 任务模板筛选主机的可多选条件集合
		}
		// -----------SSH操作相关-----------
		sshRoute := v1.Group("ssh")
		{
			sshRoute.POST("test-ssh", TestSSH) // 测试ssh
		}
		// -----------Webssh相关-------------
		websshRoute := v1.Group("webssh")
		{
			websshRoute.GET("webssh-conn", WebsshConn)                   // 连接webssh
			websshRoute.GET("webssh-exist-date", GetWebsshRecordLogDate) // 查询存在的日期
			websshRoute.GET("webssh-record", GetWebsshRecord)            // 查询Webssh操作记录
		}
		// -----------运维操作相关-----------
		opsRoute := v1.Group("ops")
		{
			// 工单相关
			opsRoute.POST("submit-task", SubmitTask)        // 装服是SSH模式，更倾向做http的C/S模式，更稳定且反馈也更详细
			opsRoute.GET("task", GetTask)                   // 查看任务工单
			opsRoute.GET("ssh-exec-param", GetSSHExecParam) // 提取SSH执行参数
			opsRoute.PUT("approve-task", ApproveTask)       // 用户审批工单
			opsRoute.PUT("status", UpdateTaskStatus)        // 更改工单状态
			opsRoute.DELETE("task", DeleteTask)             // 删除任务工单
			opsRoute.POST("exec-ssh-task", OpsExecSSHTask)  // 执行人执行工单操作
		}
		// -----------服务端操作相关-----------
		serverRoute := v1.Group("server")
		{
			serverRoute.PUT("record", UpdateServerRecord)    // 更改单服记录列表
			serverRoute.GET("record", GetServerRecord)       // 获取单服记录列表
			serverRoute.DELETE("record", DeleteServerRecord) // 删除单服记录
		}
	}
	return r
}
