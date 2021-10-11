package service

import (
	"fmt"
	"main/config"
	"main/internal/pkg/logger"
	"main/internal/pkg/web/model"

	"github.com/gin-gonic/gin"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// CheckLogin 登录检查
func CheckLogin(username, password string) (uid string, ok bool) {
	var user model.User
	err := model.FindUserByUsername(username, &user)
	if err != nil {
		return user.UID, false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+config.Conf.Svr.AppKey)); err != nil {
		return user.UID, false
	}
	return user.UID, true
}

// CreateUser 创建用户
func CreateUser(username, password string, roleID uint) (ok bool) {
	uid := uuid.NewV1().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(password+config.Conf.Svr.AppKey), bcrypt.DefaultCost)
	if err != nil {
		logger.ERROR("hash generate fail")
		return false
	}
	err = model.CreateUser(&model.User{
		UID:      uid,
		Username: username,
		Password: string(hash),
		RoleID:   roleID,
	})
	if err != nil {
		logger.ERROR("create user fail")
		return false
	}
	return true
}

// ChangeUserRole 更改用户角色
func ChangeUserRole(uid string, roleID uint) (ok bool) {
	var user model.User
	if err := model.FindUserByUID(uid, &user); err != nil {
		logger.ERROR("getUserByUID fail, uid:" + uid)
		return false
	}
	user.RoleID = roleID
	if err := model.UpdateUser(&user); err != nil {
		logger.ERROR("update user in db fail, uid:" + uid)
		return false
	}

	return true
}

// ChangePassword 更改用户密码
func ChangeUserPassword(uid, oldPwd, newPwd string) (ok bool) {
	var user model.User
	if err := model.FindUserByUID(uid, &user); err != nil {
		logger.ERROR("getUserByUID fail, uid:" + uid)
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd+config.Conf.Svr.AppKey)); err != nil {
		logger.ERROR("old password not correct")
		return false
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd+config.Conf.Svr.AppKey), bcrypt.DefaultCost)
	if err != nil {
		logger.ERROR("generate password hash fail")
		return false
	}
	user.Password = string(hash)

	if err := model.UpdateUser(&user); err != nil {
		logger.ERROR("update user in db fail, uid:" + uid)
		return false
	}

	return true
}

// FindUserByUID 通过uid获取用户
func FindUserByUID(uid string, user *model.User) (ok bool) {
	if err := model.FindUserByUID(uid, user); err != nil {
		logger.ERROR(fmt.Sprintf("findUserByUID fail, uid: %s", uid))
		return false
	}

	return true
}

// GetUserList 获取用户列表
func GetUserList(offset, limit int) (res gin.H) {
	total, users, err := model.GetUserList(offset, limit)
	if err != nil {
		return
	}
	var data []gin.H
	for index := range users {
		data = append(data, gin.H{
			"id":         users[index].ID,
			"uid":        users[index].UID,
			"username":   users[index].Username,
			"role":       users[index].Role.Name,
			"created_at": users[index].CreatedAt,
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

// UpdateUser修改用户信息
func UpdateUser(uid, username string, roleID uint) (ok bool) {
	var user model.User
	if err := model.FindUserByUID(uid, &user); err != nil {
		return false
	}

	user.Username = username
	user.RoleID = roleID

	if err := model.UpdateUser(&user); err != nil {
		return false
	}

	return true
}

// DeleteUser 删除用户
func DeleteUser(uid string) (ok bool) {
	var user model.User
	user.UID = uid
	if err := model.DeleteUserByUID(uid); err != nil {
		logger.ERROR(fmt.Sprintf("uid:[%s] delete fail, msg:[%s]", uid, err.Error()))
		return false
	}
	return true
}
