//go:build windows

package main

import (
	"fqhWeb/configs"
	"fqhWeb/internal/controller"
	"fqhWeb/internal/crontab"
	"fqhWeb/internal/model"
	"log"
)

func main() {
	configs.Init()
	model.Database()
	crontab.Cron()
	if configs.Conf.System.Mode != "product" {
		err := model.DB.AutoMigrate(
			&model.User{},
			&model.UserGroup{},
			&model.Menus{},
			&model.JwtBlacklist{},
			&model.Api{},
			&model.ActRecord{},
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

	r := controller.NewRoute()

	err := r.Run(":9081")
	if err != nil {
		log.Fatalf("启动报错: \n%v", err)
		return
	}
}
