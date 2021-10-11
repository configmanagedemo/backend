package router

import (
	"net/http"

	"main/internal/pkg/logger"
	"main/internal/pkg/web/controller"
	"main/internal/pkg/web/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter for load router
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(logger.LogMiddleware())
	r.Use(gin.Recovery())
	gin.SetMode(gin.DebugMode)

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/test/middleware", middleware.AuthRequired(), func(c *gin.Context) {
			uid, _ := c.Get("uid")
			c.JSON(http.StatusOK, gin.H{
				"uid": uid,
				"msg": "ok",
			})
		})

		apiv1.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",
			})
		})
		// 上传.b文件
		apiv1.POST("/bfile", middleware.AuthRequired(), controller.UploadBFile)
		// 获取.b文件信息
		apiv1.GET("/bfile/:bfile_id", middleware.AuthRequired(), controller.GetBFileInfoByID)
		// 下载.b文件
		apiv1.GET("/bfile/:bfile_id/download", middleware.AuthRequired(), controller.DownloadBFile)
		// 获取所有.b文件名列表
		apiv1.GET("/bfiles", middleware.AuthRequired(), controller.GetBFileNameList)
		// 获取某个文件名的.b文件列表
		apiv1.GET("/bfiles/:filename", middleware.AuthRequired(), controller.GetBFileList)
		// 上报状态
		apiv1.POST("/task/:task_id", middleware.AuthRequired(), controller.ReportFinishTask)
		// 根据key获取任务列表
		apiv1.GET("/tasks", middleware.AuthRequired(), controller.GetTaskList)
		// 回滚.b文件
		apiv1.POST("/rollback/:bfile_id", middleware.AuthRequired(), controller.RollbackBFile)

		// 登录
		apiv1.POST("/login", controller.Login)
		// 创建用户
		apiv1.POST("/user", middleware.AuthRequired(), controller.CreateUser)
		// 获取用户信息
		apiv1.GET("/user/:uid", middleware.AuthRequired(), controller.GetUserByUID)
		// 修改用户信息
		apiv1.PUT("/user/:uid", middleware.AuthRequired(), controller.UpdateUserByUID)
		// 删除用户
		apiv1.DELETE("/user/:uid", middleware.AuthRequired(), controller.DeleteUserByUID)
		// 获取用户信息
		apiv1.GET("/my", middleware.AuthRequired(), controller.GetMyInfo)
		// 修改密码
		apiv1.PUT("/password", middleware.AuthRequired(), controller.ChangePassword)
		// 登出
		apiv1.GET("/logout", middleware.AuthRequired(), controller.Logout)
		// 获取用户列表
		apiv1.GET("/users", middleware.AuthRequired(), controller.GetUserList)

		// 获取角色列表
		apiv1.GET("/roles", middleware.AuthRequired(), controller.GetRoleList)
		// 获取指定角色
		apiv1.GET("/role/:role_id", middleware.AuthRequired(), controller.GetRoleByID)
		// 创建角色
		apiv1.POST("/role", middleware.AuthRequired(), controller.CreateRole)
		// 修改角色
		apiv1.PUT("/role/:role_id", middleware.AuthRequired(), controller.UpdateRole)
		// 删除角色
		apiv1.DELETE("/role/:role_id", middleware.AuthRequired(), controller.DeleteRole)
	}
	// webhook test
	r.POST("/webhook/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success.")
	})
	return r
}
