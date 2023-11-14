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
		apiRouter.POST("getApiList", GetApiList)     // api列表
		apiRouter.GET("fresh", FreshCasbin)          // 刷新casbin缓存
		apiRouter.GET("casbinList", GetCasbinList)   // 获取用户已有的API权限列表
		apiRouter.POST("updateApi", UpdateApi)       // api添加/修改
		apiRouter.DELETE("delApi", DeleteApi)        // 删除api
		apiRouter.POST("updateCasbin", UpdateCasbin) // 为用户分配API权限
		// ------------用户相关------------
		userRoute := v1.Group("user")
		{
			userRoute.POST("update", UpdateUser)
			userRoute.GET("search", GetUserList)
			userRoute.DELETE("delete", DeleteUser)
			userRoute.PATCH("status", UpdateStatus)
			userRoute.PATCH("password", UpdatePasswd)
			userRoute.PATCH("selfPassword", UpdateSelfPasswd)
			userRoute.GET("getSelfInfo", GetSelfInfo)
			userRoute.GET("getAssGroup", GetAssGroup)
			userRoute.GET("getSelfAssGroup", GetSelfAssGroup)
			userRoute.PUT("logout", UserLogout)
			userRoute.GET("getActLog", GetRecordList)
			userRoute.POST("keyFile", UpdateKeyFileContext)
			userRoute.POST("keyStr", UpdateKeyContext)
		}
		// ------------用户组相关--------------
		// 暂未实现单用户多用户组的权限绑定，只能一个用户对一个工作室，以后功能齐全再写
		groupRoute := v1.Group("group")
		{
			groupRoute.POST("update", UpdateGroup)
			groupRoute.PUT("association", UpdateUserAss)
			groupRoute.DELETE("delete", DeleteUserGroup)
			groupRoute.GET("getGroups", GetGroup)
			groupRoute.GET("getUserAss", GetAssUser)
			groupRoute.GET("getProjectAss", GetAssProject)
		}
		// ------------菜单相关--------------
		menuRoute := v1.Group("menu")
		{
			menuRoute.POST("update", UpdateMenu)
			menuRoute.PUT("association", UpdateMenuAss)
			menuRoute.GET("getMenus", GetMenuList)
			menuRoute.DELETE("delete", DeleteMenu)
		}
		// -----------项目相关-------------
		projectRoute := v1.Group("project")
		{
			projectRoute.POST("update", UpdateProject)
			projectRoute.GET("getProject", GetProject)
			projectRoute.GET("getSelfProject", GetSelfProjectList)
			projectRoute.GET("getHost", GetHostAss)
			projectRoute.PUT("association", UpdateHostAss)
			projectRoute.DELETE("delete", DeleteProject)
		}
		// -----------主机相关-------------
		hostRoute := v1.Group("host")
		{
			hostRoute.POST("update", UpdateHost)
			hostRoute.POST("updateDomain", UpdateDomain)
			hostRoute.PUT("association", UpdateProjectAss)
			hostRoute.PUT("assDomain", UpdateDomainAss)
			hostRoute.DELETE("delete", DeleteHost)
			hostRoute.DELETE("deleteDomain", DeleteDomain)
			hostRoute.GET("Host", GetHost)
			hostRoute.GET("Project", GetProjectAss)
			hostRoute.GET("domainAssHost", GetDomainAssHost)
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
		}
	}
	return r
}
