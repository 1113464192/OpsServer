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
			&model.Project{},
			&model.JwtBlacklist{},
			&model.Api{},
			&model.ActRecord{},
			&model.Project{},
			&model.Host{},
		)
	}

	r := controller.NewRoute()
	r.Run(":9081")
}
