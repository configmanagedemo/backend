package model

// GetAuthListByUid 根据uid获取权限列表
func GetAuthListByUID(uid string) (authList []Auth, err error) {
	var user User
	result := db.Model(&User{}).Preload("Role").Where(&User{UID: uid}).First(&user)
	return user.Role.AuthInfo.Auth, result.Error
}

// GetAuthList 获取权限列表
func GetAuthList() (authList []Auth, err error) {
	var authInfo AuthInfo
	err = db.Model(&Role{}).Find(&authInfo).Error
	authList = authInfo.Auth
	return
}
