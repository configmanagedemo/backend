package model

// CreateFile 创建文件
func CreateFile(file *File) (err error) {
	return db.Create(file).Error
}

// DeleteFile 删除文件
func DeleteFile(file *File) (err error) {
	return db.Delete(file).Error
}
