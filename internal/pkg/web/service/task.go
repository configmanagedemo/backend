package service

import (
	"fmt"
	"main/internal/pkg/web/model"

	"github.com/gin-gonic/gin"
)

// CreateTask 创建任务
func CreateTask(filename, uid, fileHash string, fileID uint) (id uint, ok bool) {
	task := model.Task{
		Key:         filename,
		Type:        fmt.Sprintf("%d", fileID),
		Status:      0,
		OperatorUID: uid,
	}
	if err := model.CreateTask(&task); err != nil {
		return 0, false
	}

	if err := Notify(&gin.H{
		"filename": filename,
		"file_id":  fileID,
		"task_id":  task.ID,
		"file_md5": fileHash,
	}); err != nil {
		task.Status = 2
		if err := model.UpdateTask(&task); err != nil {
			return task.ID, false
		}
	}

	return task.ID, true
}

// FinishTask 完成任务
func FinishTask(id uint) bool {
	task := model.Task{}
	if err := model.FindTaskByID(&task, id); err != nil {
		return false
	}
	task.Status = 1
	if err := model.UpdateTask(&task); err != nil {
		return false
	}
	return true
}

// GetTaskListByKey 通过key获取任务列表
func GetTaskListByKey(key string, offset, limit int) (res gin.H) {
	tasks, total, _ := model.GetTaskListByKey(key, offset, limit)
	var data []gin.H

	for index := range tasks {
		data = append(data, gin.H{
			"id":         tasks[index].ID,
			"key":        tasks[index].Key,
			"type":       tasks[index].Type,
			"status":     tasks[index].Status,
			"operator":   tasks[index].User.Username,
			"created_at": tasks[index].CreatedAt,
			"updated_at": tasks[index].UpdatedAt,
		})
	}
	res = gin.H{
		"total":  total,
		"offset": offset,
		"limit":  limit,
		"data":   data,
	}
	return
}
