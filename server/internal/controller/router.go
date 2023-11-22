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
	// ------------验证相关------------
	v1.Use(middleware.JWTAuthMiddleware()).Use(middleware.CasbinHandler()).Use(middleware.UserActionRecord())
	{
		// -------------接口权限测试--------------
		v1.GET("ping2", Test2)

		// ------------API Casbin RBAC相关----------------
		apiRouter := v1.Group("api")
		apiRouter.GET("getApiList", GetApiList)       // api列表
		apiRouter.POST("fresh", FreshCasbin)          // 刷新casbin缓存
		apiRouter.GET("getCasbinList", GetCasbinList) // 获取用户已有的API权限列表
		apiRouter.POST("updateApi", UpdateApi)        // api添加/修改
		apiRouter.DELETE("delApi", DeleteApi)         // 删除api
		apiRouter.POST("updateCasbin", UpdateCasbin)  // 为用户组分配API权限
		// ------------用户相关------------
		userRoute := v1.Group("user")
		{
			userRoute.POST("update", UpdateUser)              // 新增/修改用户
			userRoute.GET("search", GetUserList)              // 查询用户切片
			userRoute.DELETE("delete", DeleteUser)            // 删除用户
			userRoute.PATCH("status", UpdateStatus)           // 修改用户状态
			userRoute.PATCH("password", UpdatePasswd)         // 修改用户密码
			userRoute.PATCH("selfPassword", UpdateSelfPasswd) // 修改自己的密码
			userRoute.GET("getSelfInfo", GetSelfInfo)         // 查询自己的信息
			userRoute.GET("getAssGroup", GetAssGroup)         // 根据用户ID查询组
			userRoute.GET("getSelfAssGroup", GetSelfAssGroup) // 获取自己关联的组
			userRoute.PUT("logout", UserLogout)               // 登出
			userRoute.GET("getActLog", GetRecordList)         // 查询用户所有的历史操作
			userRoute.POST("keyFile", UpdateKeyFileContext)   // 添加私钥文件
			userRoute.POST("keyStr", UpdateKeyContext)        // 添加私钥内容(与文件二选一)
		}
		// ------------用户组相关--------------
		groupRoute := v1.Group("group")
		{
			groupRoute.POST("update", UpdateGroup)         // 新增/修改组
			groupRoute.PUT("association", UpdateUserAss)   // 用户组关联用户
			groupRoute.DELETE("delete", DeleteUserGroup)   // 删除用户组
			groupRoute.GET("getGroups", GetGroup)          // 查询组切片
			groupRoute.GET("getUserAss", GetAssUser)       // 查询组关联的用户
			groupRoute.GET("getProjectAss", GetAssProject) // 查询组关联的项目
		}
		// ------------菜单相关--------------
		menuRoute := v1.Group("menu")
		{
			menuRoute.POST("update", UpdateMenu)        // 新增/修改组
			menuRoute.PUT("association", UpdateMenuAss) // 菜单关联组
			menuRoute.GET("getMenus", GetMenuList)      // 获取用户组ID对应菜单
			menuRoute.DELETE("delete", DeleteMenu)      // 删除菜单
		}
		// -----------项目相关-------------
		projectRoute := v1.Group("project")
		{
			projectRoute.POST("update", UpdateProject)             // 新增/修改项目
			projectRoute.GET("getProject", GetProject)             // 查询项目
			projectRoute.GET("getSelfProject", GetSelfProjectList) // 获取自身所有项目
			projectRoute.GET("getAssHost", GetHostAss)             // 查询项目关联的服务器
			projectRoute.PUT("association", UpdateHostAss)         // 项目关联服务器
			projectRoute.DELETE("delete", DeleteProject)           // 删除项目
		}
		// -----------主机相关-------------
		hostRoute := v1.Group("host")
		{
			hostRoute.POST("update", UpdateHost)             // 新增/修改服务器
			hostRoute.GET("getPasswd", GetHostPasswd)        // 获取服务器密码
			hostRoute.POST("updateDomain", UpdateDomain)     // 新增/修改的域名
			hostRoute.PUT("assDomain", UpdateDomainAss)      // 更新域名关联的服务器
			hostRoute.DELETE("delete", DeleteHost)           // 删除主机
			hostRoute.DELETE("deleteDomain", DeleteDomain)   // 删除域名
			hostRoute.GET("Host", GetHost)                   // 获取主机当前的状态
			hostRoute.GET("domainAssHost", GetDomainAssHost) // 获取域名关联服务器
		}
		// -----------任务模板相关-----------
		taskRoute := v1.Group("task")
		{
			taskRoute.POST("template", UpdateTaskTemplate)
			taskRoute.PUT("association", UpdateTaskAssHost)
			taskRoute.GET("getTemplate", GetProjectTask)
			taskRoute.DELETE("deleteTemplate", DeleteTaskTemplate)
			taskRoute.GET("conditionSet", GetConditionSet)
		}
		// -----------SSH操作相关-----------
		sshRoute := v1.Group("ssh")
		{
			sshRoute.POST("testSSH", TestSSH)
		}
		// -----------运维操作相关-----------
		opsRoute := v1.Group("ops")
		{
			// 工单相关
			opsRoute.POST("submitTask", SubmitTask)
			opsRoute.GET("getTask", GetTask)
			opsRoute.GET("getExecParam", GetExecParam)
			opsRoute.PUT("approveTask", ApproveTask)
			opsRoute.DELETE("delete", DeleteTask)
			opsRoute.POST("execTask", OpsExecTask)
		}
	}
	return r
}
