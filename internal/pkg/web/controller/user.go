package controller

import (
	"main/internal/pkg/cache"
	"main/internal/pkg/e"
	"main/internal/pkg/e/auths"
	"main/internal/pkg/logger"
	"main/internal/pkg/web/model"
	"main/internal/pkg/web/service"

	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func Login(c *gin.Context) {
	req := loginReq{}
	if err := c.ShouldBind(&req); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	var (
		uid string
		ok  bool
	)
	if uid, ok = service.CheckLogin(req.Username, req.Password); !ok {
		service.ResponseError(c, e.LoginFail, "login fail, wrong username or password")
		return
	}

	key := service.GetUUID()
	if err := cache.Set(key, uid, 3600); err != nil {
		logger.ERROR("cache set fail")
		service.ResponseError(c, e.LoginFail, "login fail")
		return
	}

	host := c.GetHeader("Host")
	c.SetCookie("token", key, 3600, "/", host, false, true)
	// refresh auth
	service.RefreshAuthCache(uid)

	service.ResponseOK(c, &gin.H{
		"ok": true,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	key := ""
	if token := c.GetHeader("x-token"); token != "" {
		_ = cache.Del(token)
		key = token
	}
	if token, err := c.Cookie("token"); err == nil {
		_ = cache.Del(token)
		key = token
	}
	host := c.GetHeader("Host")
	c.SetCookie("token", key, 0, "/", host, false, true)
	service.ResponseOK(c, &gin.H{})
}

type createUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	req := createUserReq{}
	if err := c.ShouldBind(&req); err != nil {
		logger.ERROR(err.Error())
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	uid, exists := service.GetUID(c)
	if !exists {
		service.ResponseError(c, e.NotLogin, "")
		return
	}

	if !service.CheckAuth(uid, auths.AuthCreateUser, "") {
		service.ResponseError(c, e.NoAuth, "no auth for create user")
		return
	}

	if ok := service.CreateUser(req.Username, req.Password, req.RoleID); !ok {
		service.ResponseError(c, e.Error, "create user fail")
		return
	}

	service.ResponseOK(c, &gin.H{
		"ok": true,
	})
}

type changePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var req changePasswordReq
	if err := c.ShouldBind(&req); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	uid, exists := service.GetUID(c)
	if !exists {
		service.ResponseError(c, e.NotLogin, "")
		return
	}

	if ok := service.ChangeUserPassword(uid, req.OldPassword, req.NewPassword); !ok {
		service.ResponseError(c, e.Error, "change password fail, wrong password")
		return
	}
	service.ResponseOK(c, &gin.H{})
}

// GetMyInfo 获取用户个人信息
func GetMyInfo(c *gin.Context) {
	uid, exists := service.GetUID(c)
	if !exists {
		service.ResponseError(c, e.NotLogin, "")
		return
	}
	var user model.User
	if ok := service.FindUserByUID(uid, &user); !ok {
		service.ResponseError(c, e.Error, "can not find user")
		return
	}

	service.ResponseOK(c, &gin.H{
		"username": user.Username,
		"role_id":  user.RoleID,
	})
}

type getUserListReq struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit" binding:"required"`
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	req := getUserListReq{}
	if err := (c.ShouldBind(&req)); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}
	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthViewUser, "") {
		service.ResponseError(c, e.NoAuth, "no auth for get user list")
		return
	}

	data := service.GetUserList(req.Offset, req.Limit)
	service.ResponseOK(c, &data)
}

type getUserReq struct {
	UID string `uri:"uid" binding:"required"`
}

// GetUserByUID 获取用户
func GetUserByUID(c *gin.Context) {
	req := getUserReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthViewUser, "") {
		service.ResponseError(c, e.NoAuth, "no auth for get user by id")
		return
	}

	var user model.User
	if ok := service.FindUserByUID(req.UID, &user); !ok {
		service.ResponseError(c, e.Error, "can not find user")
		return
	}

	service.ResponseOK(c, &gin.H{
		"id":       user.ID,
		"uid":      user.UID,
		"username": user.Username,
		"role_id":  user.RoleID,
	})
}

type updateUserURI struct {
	UID string `uri:"uid" binding:"required"`
}
type updateUserForm struct {
	Username string `json:"username" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// UpdateUserByUID 更新用户信息
func UpdateUserByUID(c *gin.Context) {
	reqURI := updateUserURI{}
	reqForm := updateUserForm{}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	if err := c.ShouldBind(&reqForm); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthUpdateUser, "") {
		service.ResponseError(c, e.NoAuth, "no auth for update user")
		return
	}

	if ok := service.UpdateUser(reqURI.UID, reqForm.Username, reqForm.RoleID); !ok {
		service.ResponseError(c, e.Error, "update fail")
		return
	}
	service.RefreshAuthCache(reqURI.UID)
	service.ResponseOK(c, &gin.H{})
}

type deleteUserURI struct {
	UID string `uri:"uid" binding:"required"`
}

// DeleteUserByUID 删除用户
func DeleteUserByUID(c *gin.Context) {
	uid, _ := service.GetUID(c)
	reqURI := deleteUserURI{}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		service.ResponseError(c, e.LostParam, err.Error())
		return
	}

	if uid == reqURI.UID {
		service.ResponseError(c, e.Error, "can not delete yourself")
		return
	}

	if !service.CheckAuth(uid, auths.AuthDeleteUser, "") {
		service.ResponseError(c, e.NoAuth, "no auth for delete user")
		return
	}

	if ok := service.DeleteUser(reqURI.UID); !ok {
		service.ResponseError(c, e.Error, "delete user fail")
		return
	}

	service.ResponseOK(c, &gin.H{})
}
