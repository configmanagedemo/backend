package controller

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"main/internal/pkg/e"
	"main/internal/pkg/e/auths"
	"main/internal/pkg/logger"
	"main/internal/pkg/web/service"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type uploadFileReq struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Desc string                `form:"desc" binding:"required"`
}

func UploadBFile(c *gin.Context) {
	var err error
	req := uploadFileReq{}
	if err := c.ShouldBind(&req); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	var file multipart.File
	file, err = req.File.Open()
	if err != nil {
		service.ResponseError(c, e.Error, err.Error())
		return
	}

	fileData, _ := ioutil.ReadAll(file)
	uid, exists := service.GetUID(c)
	if !exists {
		service.ResponseError(c, e.NotLogin, "")
		return
	}
	// auth check
	if !service.CheckAuth(uid, auths.AuthUploadFile, req.File.Filename) {
		service.ResponseError(c, e.NoAuth, "no auth to upload "+req.File.Filename)
		return
	}

	var (
		bfileID uint
		ok      bool
	)
	fileHash := fmt.Sprintf("%x", md5.Sum(fileData)) // nolint:gosec

	if bfileID, ok = service.UploadBFile(fileData, req.File.Filename, uint(req.File.Size), fileHash, req.Desc, uid); !ok {
		service.ResponseError(c, e.Error, "upload fail")
		return
	}

	var taskID uint
	if taskID, ok = service.CreateTask(req.File.Filename, uid, fileHash, bfileID); !ok {
		service.ResponseError(c, e.Error, "create task fail")
		return
	}

	service.ResponseOK(c, &gin.H{
		"filename":  req.File.Filename,
		"size":      req.File.Size,
		"bfile_id":  bfileID,
		"task_id":   taskID,
		"file_hash": fileHash,
	})
}

func AddBFile(c *gin.Context) {

}

func DeleteBFile(c *gin.Context) {

}

type downloadBFileReq struct {
	ID uint `uri:"bfile_id" binding:"required"`
}

func DownloadBFile(c *gin.Context) {
	req := downloadBFileReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		logger.DEBUG(err.Error())
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	data, filename, _ := service.GetBFile(req.ID)

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(data)))
	if _, err := c.Writer.Write(data); err != nil {
		logger.ERROR("echo file error," + err.Error())
	}
}

func GetBFileNameList(c *gin.Context) {
	bfileNameList := service.GetBFileNameList()

	service.ResponseOK(c, &gin.H{
		"bfilename_list": bfileNameList,
	})
}

type getBFileListURI struct {
	Filename string `uri:"filename" binding:"required"`
}
type getBFileListQuery struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit" binding:"required"`
}

func GetBFileList(c *gin.Context) {
	reqURI := getBFileListURI{}
	reqQuery := getBFileListQuery{}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}
	if err := c.ShouldBind(&reqQuery); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}
	data := service.GetBFileListByFilename(reqURI.Filename, reqQuery.Offset, reqQuery.Limit)
	service.ResponseOK(c, &data)
}

type getBFileInfoQuery struct {
	FileID int `uri:"bfile_id" binding:"required"`
}

func GetBFileInfoByID(c *gin.Context) {
	reqQuery := getBFileInfoQuery{}
	if err := c.ShouldBindUri(&reqQuery); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}
	data := service.GetBFileInfo(uint(reqQuery.FileID))
	service.ResponseOK(c, &data)
}

type rollbackBFileURI struct {
	FileID uint `uri:"bfile_id" binding:"required"`
}

func RollbackBFile(c *gin.Context) {
	reqURI := rollbackBFileURI{}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}
	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthUploadFile, "") {
		service.ResponseError(c, e.NoAuth, "no auth for rollback")
		return
	}
	ok, filename, hash := service.RollbackBFile(reqURI.FileID)
	if !ok {
		service.ResponseError(c, e.Error, "rollback bfile fail")
		return
	}
	var taskID uint
	if taskID, ok = service.CreateTask(filename, uid, hash, reqURI.FileID); !ok {
		service.ResponseError(c, e.Error, "create task fail")
		return
	}
	service.ResponseOK(c, &gin.H{
		"task_id": taskID,
	})
}
