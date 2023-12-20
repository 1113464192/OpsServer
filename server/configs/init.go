package configs

import (
	"fmt"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile(util.GetRootPath() + "/configs/config.yaml")
	// viper.SetConfigFile(util.GetRootPath() + "\\configs\\config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置信息失败: %s ", err))
	}
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("获取配置失败, err:%s ", err))
	}
	logger.BuildLogger(Conf.Logger.Level)
	// 监视配置文件的变化。当配置文件发生改变时，它会自动重新加载并更新 Viper 实例中的配置信息。这样可以避免在应用程序运行时手动重新加载配置文件
	viper.WatchConfig()
	// OnConfigChange，回调函数，用于在配置文件发生变化时进行处理。您可以将自定义的函数传递给 OnConfigChange()，在配置文件发生更改时，该函数将被调用
	viper.OnConfigChange(func(in fsnotify.Event) { // 传递配置文件变更事件的参数类型，以便在 OnConfigChange() 回调函数中获取有关配置文件变化的详细信息。
		logger.Log().Warning("config", "Conf", "配置文件触发修改重载"+in.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			logger.Log().Panic("config", "Conf", "配置文件写入结构体变量失败")
		}
	})

}
