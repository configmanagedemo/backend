package service

import (
	"encoding/json"
	"fmt"
	"main/internal/pkg/cache"
	"main/internal/pkg/e/auths"
	"main/internal/pkg/logger"
	"main/internal/pkg/web/model"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// getAuthInfo 根据uid从缓存和DB获取权限信息
func getAuthInfo(uid string) (authList []model.Auth) {
	var (
		err        error
		authString string
	)
	authString, err = cache.Get(uid + "Auth")
	if err != nil {
		if err != cache.Nil {
			logger.ERROR(fmt.Sprintf("redis get auth info fail, uid:%s", uid))
		}
		RefreshAuthCache(uid)
		authString, err = cache.Get(uid + "Auth")
		if err != nil {
			logger.ERROR(fmt.Sprintf("redis set auth cache fail, uid:%s", uid))
		}
	}
	err = json.Unmarshal([]byte(authString), &authList)
	if err != nil {
		logger.ERROR("getAuthInfo json unmarshal fail")
		return
	}
	return authList
}

// CheckAuth 根据用户 flag检测权限
func CheckAuth(uid, flag, resource string) (ok bool) {
	if flag == "" {
		return true
	}

	authList := getAuthInfo(uid)
	for _, auth := range authList {
		if auth.Flag == auths.AuthAll {
			return true
		}
		if auth.Flag != flag {
			continue
		}
		if resource == "" {
			return true
		}
		for _, res := range auth.Resource {
			if res == auths.AuthAll || res == resource {
				return true
			}
		}
	}
	return false
}

// RefreshAuthCache 刷新权限
func RefreshAuthCache(uid string) {
	authList, err := model.GetAuthListByUID(uid)
	if err != nil {
		logger.WARN(fmt.Sprintf("db get auth info fail, uid:%s", uid))
	}
	s, err := json.Marshal(authList)
	if err != nil {
		logger.ERROR("getAuthInfo json marshal fail")
		return
	}
	fmt.Printf("uid:%s set auth:%s\n", uid, s)
	err = cache.Set(uid+"Auth", string(s), time.Second*5)
	if err != nil {
		logger.ERROR("getAuthInfo set auth info to cache fail")
		return
	}
}

// GetRoleList 获取权限列表
func GetRoleList() (data []gin.H) {
	var (
		roleList []model.Role
		err      error
	)
	if roleList, err = model.GetRoleList(); err != nil {
		logger.WARN("GetRoleList fail")
		return
	}

	for index := range roleList {
		data = append(data, gin.H{
			"role_id":    roleList[index].ID,
			"role_name":  roleList[index].Name,
			"role_desc":  roleList[index].Desc,
			"created_at": roleList[index].CreatedAt,
		})
	}
	return
}

// CreateRole 创建角色
func CreateRole(roleName, desc string, authInfo model.AuthInfo) (ok bool, roleID uint) {
	var role model.Role
	role.Name = roleName
	role.Desc = desc
	role.AuthInfo = authInfo
	roleID, err := model.CreateRole(&role)
	if err != nil {
		logger.ERROR(fmt.Sprintf("db create role fail, err:%s", err.Error()))
		return false, 0
	}
	return true, roleID
}

// GetRoleByRoleID 角色ID获取角色信息
func GetRoleByRoleID(roleID uint) (data gin.H) {
	role, err := model.GetRole(roleID)
	fmt.Println(role)
	if err != nil {
		logger.ERROR(fmt.Sprintf("db get role fail, roleID:%d, err:%s", roleID, err.Error()))
		return
	}
	// authInfo, err := json.Marshal(role.AuthInfo)
	// if err != nil {
	// 	logger.ERROR("json marshal fail")
	// }
	data = gin.H{
		"role_id":   role.ID,
		"role_name": role.Name,
		"desc":      role.Desc,
		"auth":      role.AuthInfo.Auth,
	}
	return data
}

// UpdateRole 更新角色
func UpdateRole(roleID uint, roleName, desc string, authInfo model.AuthInfo) (ok bool) {
	role, err := model.GetRole(roleID)
	if err != nil {
		logger.ERROR(fmt.Sprintf("db get role fail, roleID:%d, err:^%s", roleID, err.Error()))
		return false
	}
	role.Name = roleName
	role.Desc = desc
	role.AuthInfo = authInfo
	err = model.UpdateRole(&role)
	if err != nil {
		logger.ERROR(fmt.Sprintf("db update role fail, roleID:%d, err:%s", roleID, err.Error()))
		return false
	}
	return true
}

// DeleteRole 删除角色
func DeleteRole(roleID uint) (ok bool) {
	var role model.Role
	role.ID = roleID
	if err := model.DeleteRole(&role); err != nil {
		logger.ERROR(fmt.Sprintf("db delete role fail, roleID:%d, err:%s", roleID, err.Error()))
		return false
	}
	return true
}

// CreateToken 创建一个token
func CreateToken(expireAt time.Time, uid string) (ok bool) {
	tokenStr := uuid.NewV1().String()
	token := model.Token{
		Token:    tokenStr,
		ExpireAt: expireAt,
		Enable:   true,
		UID:      uid,
	}
	if err := model.CreateToken(&token); err != nil {
		return false
	}
	return true
}

// ChangeTokenEnable 更改token信息
func ChangeTokenEnable(id uint, enable bool) (ok bool) {
	if err := model.UpdateTokenByMap(id, map[string]interface{}{
		"enable": enable,
	}); err != nil {
		return false
	}
	return true
}

// ChangeTokenExpire 更改token过期信息
func ChangeTokenExpire(id uint, expireAt time.Time) (ok bool) {
	if err := model.UpdateTokenByMap(id, map[string]interface{}{
		"expire_at": expireAt,
	}); err != nil {
		return false
	}
	return true
}

// GetTokenList 获取token列表
func GetTokenList(offset, limit int) (res gin.H) {
	total, tokens, _ := model.GetTokenList(offset, limit)
	data := []gin.H{}
	for idx := range tokens {
		data = append(data, gin.H{
			"id":        tokens[idx].ID,
			"token":     tokens[idx].Token,
			"expire_at": tokens[idx].ExpireAt,
			"enable":    tokens[idx].Enable,
			"uid":       tokens[idx].UID,
			"username":  tokens[idx].User.Username,
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
