package model

import (
	"fmt"
	"main/internal/pkg/logger"

	"gorm.io/gorm"
)

// CreateBFile 创建.b文件
func CreateBFile(bfile *BFile) (err error) {
	result := db.Create(bfile)
	return result.Error
}

// UpdateBFile 更新.b文件
func UpdateBFile(bfile *BFile) (err error) {
	result := db.Save(bfile)
	return result.Error
}

// DeleteBFile 删除.b文件
func DeleteBFile(bfile *BFile) (err error) {
	result := db.Delete(bfile)
	return result.Error
}

// GetBFile 获取.b文件数据
func GetBFile(fileID uint, bfile *BFile) (err error) {
	db.Where(gorm.Model{ID: fileID}).First(bfile)
	return db.Model(&bfile).Association("File").Find(&bfile.File)
}

// UpdateFileSetUnUse 设置文件未使用
func UpdateFileSetUnUse(filename string) (err error) {
	result := db.Model(&BFile{}).Where("filename = ?", filename).Update("is_use", false)
	logger.INFO(fmt.Sprintf("affect row %d", result.RowsAffected))
	return result.Error
}

// UpdateFileSetUsed 设置文件使用
func UpdateFileSetUsed(fileID uint) (err error) {
	err = db.Model(&BFile{}).Where("id = ?", fileID).Update("is_use", true).Error
	return err
}

// GetBFileNameList 获取.b文件名列表
func GetBFileNameList() (filenameList []string, err error) {
	err = db.Model(&BFile{}).Distinct().Pluck("filename", &filenameList).Error
	return
}

// GetBFileListByFilename 根据文件名获取.b文件列表
func GetBFileListByFilename(filename string, offset, limit int) (bfiles []BFile, total int64, err error) {
	err = db.Model(&BFile{}).Where("filename = ?", filename).Count(&total).Error
	if err != nil {
		return
	}
	err = db.Model(&BFile{}).Joins("User").Where(&BFile{Filename: filename}).
		Order("created_at desc").Offset(offset).Limit(limit).Find(&bfiles).Error
	return
}
