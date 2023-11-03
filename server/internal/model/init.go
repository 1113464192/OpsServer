package model

import (
	"fmt"
	"fqhWeb/configs"
	m_logger "fqhWeb/pkg/logger"
	"fqhWeb/pkg/utils"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var LogFile *os.File

func Database() {
	if utils.IsDir(utils.GetRootPath() + "/logs/mysql/") {
		os.Mkdir(utils.GetRootPath()+"/logs/mysql/", 0777)
	}

	logFileName := utils.GetRootPath() + "/logs/mysql/mysql.log"

	var err error

	LogFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("无法打开日志文件:", err)
		return
	}

	newLogger := logger.New(
		log.New(LogFile, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	dsn := configs.Conf.Mysql.Conf
	if dsn == "" {
		m_logger.Log().Error("mysql", "")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "fqh_",
			SingularTable: true,
		},
		Logger:          newLogger,
		CreateBatchSize: configs.Conf.Mysql.CreateBatchSize,
		NowFunc: func() time.Time {
			tmp := time.Now().Local().Format("2006-01-02 15:04:05")
			now, _ := time.ParseInLocation("2006-01-02 15:04:05", tmp, time.Local)
			return now
		},
	})
	if err != nil {
		m_logger.Log().Error("mysql", "service_log", err)
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		m_logger.Log().Error("mysql", "service_log", err)
		panic(err)
	}
	// 连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	// 设置连接的最大生命周期为 0，这意味着连接在连接池中没有最大生命周期的限制，它可以一直保持打开状态
	sqlDB.SetConnMaxLifetime(0)
	DB = db
}
