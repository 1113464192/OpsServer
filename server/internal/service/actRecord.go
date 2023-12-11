package service

import (
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"strings"
	"time"
)

type ActRecord struct {
}

var (
	insLog = ActRecord{}
)

func Record() *ActRecord {
	return &insLog
}

func (s *ActRecord) checkRecordTableExists() (tableName string, exist bool) {
	currYearMon := time.Now().Local()
	nowTime := currYearMon.Format("2006_01")
	tableName = fmt.Sprintf("%sact_record_%s", configs.Conf.Mysql.TablePrefix, nowTime)
	// 等待表的创建，最多等待20s
	for i := 0; i < 5; i++ {
		if model.DB.Migrator().HasTable(tableName) {
			exist = true
			break
		}
		_ = model.DB.Table(tableName).AutoMigrate(&model.ActRecord{})
		time.Sleep(time.Second)
	}
	return tableName, exist
}

// RecordCreate 插入日志
func (s *ActRecord) RecordCreate(log *model.ActRecord) (err error) {
	tableName, exist := s.checkRecordTableExists()
	if !exist {
		return errors.New("当月表尚未创建，请联系运维查看")
	}
	if err = model.DB.Table(tableName).Create(&log).Error; err != nil {
		return fmt.Errorf("在当月记录表中创建记录失败: %v", err)
	}
	return err
}

// 查询有多少个月份表可供查询
func (s *ActRecord) GetRecordLogDate() (dates []string, err error) {
	// 构建原生 SQL 查询语句
	sql := fmt.Sprintf(`SHOW TABLES LIKE '%sact_record_%%'`, configs.Conf.Mysql.TablePrefix)
	if err = model.DB.Raw(sql).Scan(&dates).Error; err != nil {
		return nil, fmt.Errorf("获取所有记录表的表名失败: %v", err)
	}
	for i := 0; i < len(dates); i++ {
		dates[i] = strings.Replace(dates[i], configs.Conf.Mysql.TablePrefix+"act_record_", "", -1)
	}
	return dates, err
}

// 查询月份记录
func (s *ActRecord) GetRecordList(param api.GetRecordListReq) (list any, total int64, err error) {
	tableName := fmt.Sprintf("%sact_record_%s", configs.Conf.Mysql.TablePrefix, param.Date)
	if !model.DB.Migrator().HasTable(tableName) {
		return nil, 0, errors.New("没有这个日期的行为记录表存在: " + tableName)
	}

	var logs []model.ActRecord
	db := model.DB.Table(tableName).Order("id desc")
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &logs,
		PageInfo:  param.PageInfo,
	}
	if param.Id != 0 {
		db = db.Where("user_id = ?", param.Id)
		searchReq.Condition = db
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
	} else {
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
	}
	return &logs, total, err
}
