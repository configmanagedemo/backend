package controller

import (
	"fmt"
	"main/internal/pkg/e"
	"main/internal/pkg/logger"
	"main/internal/pkg/web/service"

	"github.com/gin-gonic/gin"
)

type reportFinishTaskReq struct {
	TaskID uint `uri:"task_id" binding:"required"`
}

func ReportFinishTask(c *gin.Context) {
	req := reportFinishTaskReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	if ok := service.FinishTask(req.TaskID); !ok {
		service.ResponseError(c, e.Error, "finish task fail")
		return
	}

	service.ResponseOK(c, &gin.H{
		"task_id": req.TaskID,
	})
}

type getTaskListReq struct {
	Key    string `form:"key" binding:"required"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit" binding:"required"`
}

func GetTaskList(c *gin.Context) {
	req := getTaskListReq{}
	if err := c.ShouldBind(&req); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}
	logger.INFO(fmt.Sprintf("req.key: %s", req.Key))
	data := service.GetTaskListByKey(req.Key, req.Offset, req.Limit)

	service.ResponseOK(c, &data)
}
