//go:build windows
// +build windows

package main

import (
	"fqhWeb/configs"
	"fqhWeb/internal/controller"
	"fqhWeb/internal/crontab"
	"fqhWeb/internal/model"
)

func main() {
	configs.Init()
	model.Database()
	crontab.Cron()
	if configs.Conf.System.Mode != "product" {
		model.DB.AutoMigrate(
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
		)
	}

	r := controller.NewRoute()

	r.Run(":9081")
}
