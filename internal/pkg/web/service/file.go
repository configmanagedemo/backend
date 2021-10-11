package service

import (
	"fmt"
	"main/internal/pkg/cache"
	"main/internal/pkg/logger"
	"main/internal/pkg/web/model"

	"github.com/gin-gonic/gin"
)

// UploadBFile 上传.b文件
func UploadBFile(fileData []byte, filename string, fileSize uint, fileHash, desc, uploader string) (id uint, ok bool) {
	var bfile model.BFile
	bfile.Filename = filename
	bfile.FileSize = fileSize
	bfile.Desc = desc
	bfile.IsUse = true
	bfile.File.Data = fileData
	bfile.File.Hash = fileHash
	bfile.UploaderUID = uploader

	if err := model.UpdateFileSetUnUse(filename); err != nil {
		logger.ERROR(err.Error())
		return 0, false
	}

	if err := model.CreateBFile(&bfile); err != nil {
		logger.ERROR(err.Error())
		return 0, false
	}
	return bfile.ID, true
}

type tBFileData struct {
	fileData []byte
	filename string
}

// GetBFile 获取.b文件
func GetBFile(fileID uint) (fileData []byte, filename string, err error) {
	var bfile model.BFile
	// first query from lru
	if bFileData, ok := cache.LRUGet(fileID); ok {
		logger.DEBUG(fmt.Sprintf("fileid:%d find in cache.", fileID))
		fileData = bFileData.(tBFileData).fileData
		filename = bFileData.(tBFileData).filename
		return
	}
	// or find in db
	if err = model.GetBFile(fileID, &bfile); err != nil {
		logger.ERROR(err.Error())
		return
	}
	var data tBFileData
	data.fileData = bfile.File.Data
	data.filename = bfile.Filename
	logger.DEBUG(fmt.Sprintf("fileid:%d load from db.", fileID))
	cache.LRUPut(fileID, data)

	fileData = bfile.File.Data
	filename = bfile.Filename
	return
}

// GetBFileInfo 获取.b文件信息
func GetBFileInfo(fileID uint) (res gin.H) {
	var file model.BFile
	err := model.GetBFile(fileID, &file)
	if err != nil {
		logger.ERROR(fmt.Sprintf("GetBFileInfo fail, fileID[%d], err: %s", fileID, err.Error()))
	}
	res = gin.H{
		"file_id":    file.ID,
		"filename":   file.Filename,
		"filesize":   file.FileSize,
		"desc":       file.Desc,
		"is_use":     file.IsUse,
		"uploader":   file.User.Username,
		"created_at": file.CreatedAt,
	}
	return
}

// GetBFileNameList 获取.b文件名列表
func GetBFileNameList() (filenameList []string) {
	filenameList, _ = model.GetBFileNameList()
	return
}

// GetBFileListByFilename 根据文件名获取文件列表
func GetBFileListByFilename(filename string, offset, limit int) (res gin.H) {
	bfiles, total, _ := model.GetBFileListByFilename(filename, offset, limit)
	var data []gin.H
	for index := range bfiles {
		data = append(data, gin.H{
			"file_id":    bfiles[index].ID,
			"filename":   bfiles[index].Filename,
			"filesize":   bfiles[index].FileSize,
			"desc":       bfiles[index].Desc,
			"is_use":     bfiles[index].IsUse,
			"uploader":   bfiles[index].User.Username,
			"created_at": bfiles[index].CreatedAt,
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

// RollbackBFile 回滚.b文件
func RollbackBFile(fileID uint) (ok bool, filename, hash string) {
	var bfile model.BFile
	if err := model.GetBFile(fileID, &bfile); err != nil {
		logger.ERROR(err.Error())
		return false, "", ""
	}
	if err := model.UpdateFileSetUnUse(bfile.Filename); err != nil {
		logger.ERROR(err.Error())
		return false, "", ""
	}
	if err := model.UpdateFileSetUsed(fileID); err != nil {
		logger.ERROR(err.Error())
		return false, "", ""
	}
	return true, bfile.Filename, bfile.File.Hash
}
