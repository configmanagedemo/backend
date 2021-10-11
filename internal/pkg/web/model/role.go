package model

import "fmt"

// GetRoleList 获取角色列表
func GetRoleList() (roleList []Role, err error) {
	err = db.Model(&Role{}).Find(&roleList).Error
	return
}

// GetRole 获取角色
func GetRole(id uint) (role Role, err error) {
	err = db.Model(&role).First(&role, id).Error
	l := len(role.AuthInfo.Auth)
	fmt.Println(l)
	return
}

// UpdateRole 更新角色
func UpdateRole(role *Role) (err error) {
	err = db.Model(role).Save(role).Error
	return
}

// CreateRole 创建角色
func CreateRole(role *Role) (id uint, err error) {
	err = db.Model(role).Create(role).Error
	return role.ID, err
}

// DeleteRole 删除角色
func DeleteRole(role *Role) (err error) {
	err = db.Model(role).Delete(role).Error
	return
}
