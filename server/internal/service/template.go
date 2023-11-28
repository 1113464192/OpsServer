package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service/dbOper"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/util"
	"fqhWeb/pkg/util2"
)

type TaskService struct {
}

var (
	insTask = &TaskService{}
)

func Task() *TaskService {
	return insTask
}

// 修改或新增任务模板
func (s *TaskService) UpdateTaskTemplate(param *api.UpdateTaskTemplateReq) (projectInfo any, err error) {
	var task model.TaskTemplate
	var count int64
	if model.DB.Model(&task).Where("pid = ? AND type_name = ? AND task_name = ? AND id != ?", param.Pid, param.TypeName, param.TaskName, param.ID).Count(&count); count > 0 {
		return task, errors.New("该项目中的 任务类型的 任务名已被使用")
	}
	conditionJson, err := util.ConvertToJsonPair(param.Condition)
	if err != nil {
		return nil, err
	}
	portRuleJson, err := util.ConvertToJsonPair(param.PortRule)
	if err != nil {
		return nil, err
	}
	argsJson, err := util.ConvertToJsonPair(param.Args)
	if err != nil {
		return nil, err
	}

	if param.ID != 0 {
		// 修改
		if !util2.CheckIdExists(&task, param.ID) {
			return task, errors.New("不存在")
		}
		if err := model.DB.Model(&task).Where("id = ?", param.ID).First(&task).Error; err != nil {
			return task, errors.New("任务模板数据库查询失败: " + err.Error())
		}
		task.TypeName = param.TypeName
		task.TaskName = param.TaskName
		task.CmdTem = param.CmdTem
		task.ConfigTem = param.ConfigTem
		task.Condition = conditionJson
		task.Comment = param.Comment
		task.Pid = param.Pid
		task.PortRule = portRuleJson
		task.Args = argsJson
		if err = model.DB.Save(&task).Error; err != nil {
			return task, fmt.Errorf("数据保存失败: %v", err)
		}
	} else {
		task = model.TaskTemplate{
			TypeName:  param.TypeName,
			TaskName:  param.TaskName,
			CmdTem:    param.CmdTem,
			ConfigTem: param.ConfigTem,
			Condition: conditionJson,
			Comment:   param.Comment,
			Pid:       param.Pid,
			PortRule:  portRuleJson,
			Args:      argsJson,
		}
		if err = model.DB.Create(&task).Error; err != nil {
			return task, errors.New("创建项目任务模板失败")
		}
	}
	var result []api.TaskTemRes
	if result, err = s.GetTemplateResults(&task); err != nil {
		return nil, err
	}
	return result, err
}

// 获取任务模板
func (s *TaskService) GetProjectTask(param *api.GetProjectTaskReq) (projectObj any, total int64, err error) {
	var task []model.TaskTemplate
	db := model.DB.Model(&task)
	searchReq := &api.SearchReq{
		Condition: db,
		Table:     &task,
		PageInfo:  param.PageInfo,
	}
	// 如果传了模板ID
	// 直接返回模板的所有内容
	if param.ID != 0 {
		db = db.Where("id = ?", param.ID)
		searchReq.Condition = db
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
		var result []api.TaskTemRes
		if result, err = s.GetTemplateResults(&task); err != nil {
			return nil, total, err
		}
		return result, total, err
		// 如果传了类型名和项目ID
		// 返回模板名+模板ID切片
	} else if param.Pid != 0 && param.TypeName != "" {
		db = db.Where("pid = ? AND type_name = ?", param.Pid, param.TypeName)
		searchReq.Condition = db
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
		// 返回模板名切片
		var result []api.TaskInfo
		for _, record := range task {
			taskInfo := api.TaskInfo{
				TaskName: record.TaskName,
				ID:       record.ID,
			}
			result = append(result, taskInfo)
		}
		return result, total, err
		// 如果只传了项目ID
		// 只返回包含的类型名
	} else if param.Pid != 0 && param.TypeName == "" {
		db = db.Where("pid = ?", param.Pid)
		searchReq.Condition = db
		if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
			return nil, 0, err
		}
		var result []string
		for _, record := range task {
			if !util.IsSliceContain(result, record) {
				result = append(result, record.TypeName)
			}
		}
		return result, total, err
	}
	return nil, 0, errors.New("参数有误")
	// } else if param.Pid != 0 {
	// 	db = db.Where("pid = ?", param.Pid)
	// 	searchReq.Condition = db
	// 	if total, err = dbOper.DbOper().DbFind(searchReq); err != nil {
	// 		return nil, 0, err
	// 	}
	// 	// 返回 模板类型 > 模板名
	// 	var results = make(map[string][]api.TaskInfo)
	// 	for _, record := range task {
	// 		typeName := record.TypeName
	// 		taskInfo := api.TaskInfo{
	// 			TaskName: record.TaskName,
	// 			ID:       record.ID,
	// 		}
	// 		if _, ok := results[typeName]; !ok {
	// 			results[typeName] = []api.TaskInfo{taskInfo}
	// 		} else {
	// 			results[typeName] = append(results[typeName], taskInfo)
	// 		}
	// 	}
	// 	return results, total, err
	// }

}

func (s *TaskService) DeleteTaskTemplate(ids []uint) (err error) {
	if err = util2.CheckIdsExists(model.TaskTemplate{}, ids); err != nil {
		return err
	}
	var task []model.TaskTemplate
	tx := model.DB.Begin()
	if err = tx.Find(&task, ids).Error; err != nil {
		return errors.New("查询任务信息失败")
	}
	if err = tx.Model(&task).Association("Hosts").Clear(); err != nil {
		tx.Rollback()
		return errors.New("清除表信息 任务与服务器关联 失败")
	}
	if err = tx.Where("id in (?)", ids).Delete(&model.TaskTemplate{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除任务失败")
	}
	tx.Commit()
	return err
}

// 任务关联主机
func (s *TaskService) UpdateHostAss(param api.UpdateTemplateAssHostReq) (err error) {
	var host []model.Host
	var task model.TaskTemplate
	// 判断所有项目是否都存在
	if err = util2.CheckIdsExists(model.Host{}, param.Hids); err != nil {
		return err
	}

	if !util2.CheckIdExists(&task, param.Tid) {
		return errors.New("任务模板ID不存在")
	}

	if err = model.DB.Find(&host, param.Hids).Error; err != nil {
		return errors.New("主机数据库查询操作失败")
	}
	if err = model.DB.First(&task, param.Tid).Error; err != nil {
		return errors.New("任务模板数据库查询操作失败")
	}
	if err = model.DB.Model(&task).Association("Hosts").Replace(&host); err != nil {
		return errors.New("任务模板与服务器数据库关联操作失败")
	}
	if err != nil {
		return err
	}
	return err

}

// 返回模板JSON结果
func (s *TaskService) GetTemplateResults(taskInfo any) (result []api.TaskTemRes, err error) {
	var res api.TaskTemRes
	if task, ok := taskInfo.(*[]model.TaskTemplate); ok {
		for _, task := range *task {
			res = api.TaskTemRes{
				ID:        task.ID,
				TypeName:  task.TypeName,
				TaskName:  task.TaskName,
				CmdTem:    task.CmdTem,
				ConfigTem: task.ConfigTem,
				Comment:   task.Comment,
				Pid:       task.Pid,
				Condition: task.Condition,
				PortRule:  task.PortRule,
				Args:      task.Args,
			}
			result = append(result, res)
		}
		return result, err
	}
	if task, ok := taskInfo.(*model.TaskTemplate); ok {
		res = api.TaskTemRes{
			ID:        task.ID,
			TypeName:  task.TypeName,
			TaskName:  task.TaskName,
			CmdTem:    task.CmdTem,
			ConfigTem: task.ConfigTem,
			Comment:   task.Comment,
			Pid:       task.Pid,
			Condition: task.Condition,
			PortRule:  task.PortRule,
			Args:      task.Args,
		}
		result = append(result, res)
		return result, err
	}
	return result, errors.New("转换任务模板结果失败")
}
