//go:build windows

package main

import (
	"fqhWeb/configs"
	"fqhWeb/internal/controller"
	"fqhWeb/internal/crontab"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/globalFunc"
	"log"
)

func main() {
	configs.Init()
	model.Database()
	crontab.Cron()
	if configs.Conf.System.Mode != "product" {
		// 数据库这块更推荐业务层实现关联，而不使用主外键，这里为了开发速度使用了主外键
		err := model.DB.AutoMigrate(
			&model.User{},
			&model.UserGroup{},
			&model.Menus{},
			&model.JwtBlacklist{},
			&model.Api{},
			&model.TaskTemplate{},
			&model.Domain{},
			&model.Project{},
			&model.Host{},
			&model.TaskRecord{},
			&model.ServerRecord{},
			&model.GitWebhookRecord{},
		)
		if err != nil {
			log.Fatalf("自动迁移报错: \n%v", err)
			return
		}
	}
	globalFunc.DeclareGlobalVar()

	r := controller.NewRoute()

	err := r.Run(":9080")
	if err != nil {
		log.Fatalf("启动报错: \n%v", err)
		return
	}
}
