package controller

import (
	"fmt"
	"main/internal/pkg/e"
	"main/internal/pkg/e/auths"
	"main/internal/pkg/logger"
	"main/internal/pkg/web/model"
	"main/internal/pkg/web/service"

	"github.com/gin-gonic/gin"
)

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthViewRole, "") {
		service.ResponseError(c, e.NoAuth, "no auth for get role list")
		return
	}
	data := service.GetRoleList()
	service.ResponseOK(c, &gin.H{
		"roles": data,
	})
}

type createRoleReq struct {
	RoleName string `json:"role_name" binding:"required"`
	Desc     string `json:"desc"`
	Auth     []struct {
		Flag     string   `json:"flag"`     // 标志
		Desc     string   `json:"desc"`     // 描述
		Resource []string `json:"resource"` // 资源
	} `json:"auths"`
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	req := createRoleReq{}
	if err := c.ShouldBind(&req); err != nil {
		service.ResponseError(c, e.LostParam, "create role lost params")
		return
	}
	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthCreateRole, "") {
		service.ResponseError(c, e.NoAuth, "no auth for create role")
		return
	}
	var authInfo model.AuthInfo
	authInfo.Auth = []model.Auth{}
	for _, auth := range req.Auth {
		authInfo.Auth = append(authInfo.Auth, model.Auth{
			Flag:     auth.Flag,
			Desc:     auth.Desc,
			Resource: auth.Resource,
		})
	}
	ok, roleID := service.CreateRole(req.RoleName, req.Desc, authInfo)
	if !ok {
		service.ResponseError(c, e.Error, "create role fail")
		return
	}
	service.ResponseOK(c, &gin.H{
		"role_id": roleID,
	})
}

type getRoleURI struct {
	RoleID uint `uri:"role_id" binding:"required"`
}

// GetRoleByID 获取角色信息
func GetRoleByID(c *gin.Context) {
	req := getRoleURI{}
	if err := c.ShouldBindUri(&req); err != nil {
		logger.ERROR("get role by id lost params")
		service.ResponseError(c, e.LostParam, "get role lost params")
		return
	}
	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthViewRole, "") {
		service.ResponseError(c, e.NoAuth, "no auth for get role")
		return
	}
	data := service.GetRoleByRoleID(req.RoleID)
	service.ResponseOK(c, &data)
}

type updateRoleURI struct {
	RoleID uint `uri:"role_id" binding:"required"`
}
type updateRoleReq struct {
	RoleName string `json:"role_name" binding:"required"`
	Desc     string `json:"desc"`
	Auth     []struct {
		Flag     string   `json:"flag"`     // 标志
		Desc     string   `json:"desc"`     // 描述
		Resource []string `json:"resource"` // 资源
	} `json:"auths"`
}

// UpdateRole 更新角色权限
func UpdateRole(c *gin.Context) {
	req := updateRoleReq{}
	uri := updateRoleURI{}
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.ERROR("update role lost params")
		service.ResponseError(c, e.LostParam, "lost params")
		return
	}
	if err := c.ShouldBind(&req); err != nil {
		logger.ERROR("update role req lost params")
		service.ResponseError(c, e.LostParam, "lost params")
		return
	}
	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthUpdateRole, "") {
		service.ResponseError(c, e.NoAuth, "no auth for update role")
		return
	}
	var authInfo model.AuthInfo
	authInfo.Auth = []model.Auth{}
	for _, auth := range req.Auth {
		authInfo.Auth = append(authInfo.Auth, model.Auth{
			Flag:     auth.Flag,
			Desc:     auth.Desc,
			Resource: auth.Resource,
		})
	}
	if ok := service.UpdateRole(uri.RoleID, req.RoleName, req.Desc, authInfo); !ok {
		service.ResponseError(c, e.Error, "update role fail")
		return
	}
	service.ResponseOK(c, nil)
}

type deleteRoleURI struct {
	RoleID uint `uri:"role_id" binding:"required"`
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	uri := deleteRoleURI{}
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.ERROR(fmt.Sprintf("delete role lost params, err:%s", err.Error()))
		service.ResponseError(c, e.LostParam, "lost params")
		return
	}
	uid, _ := service.GetUID(c)
	if !service.CheckAuth(uid, auths.AuthDeleteRole, "") {
		service.ResponseError(c, e.NoAuth, "no auth for delete role")
		return
	}
	if !service.DeleteRole(uri.RoleID) {
		service.ResponseError(c, e.Error, "delete role fail")
		return
	}
	service.ResponseOK(c, nil)
}
