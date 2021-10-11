package model

// CreateTask 创建任务
func CreateTask(task *Task) (err error) {
	return db.Model(&task).Create(task).Error
}

// UpdateTask 更新任务
func UpdateTask(task *Task) (err error) {
	return db.Model(&task).Save(task).Error
}

// DeleteTask 删除任务
func DeleteTask(task *Task) (err error) {
	return db.Model(&task).Delete(task).Error
}

// FindTaskByID 获取任务
func FindTaskByID(task *Task, id uint) (err error) {
	return db.Model(&task).First(&task, id).Error
}

// GetTaskListByKey 获取任务列表
func GetTaskListByKey(key string, offset, limit int) (tasks []Task, total int64, err error) {
	err = db.Model(&Task{}).Where("`key` = ?", key).Count(&total).Error
	if err != nil {
		return
	}
	result := db.Model(&Task{}).Joins("User").Where(&Task{Key: key}).
		Order("created_at desc").Offset(offset).Limit(limit).Find(&tasks)
	err = result.Error
	return
}
