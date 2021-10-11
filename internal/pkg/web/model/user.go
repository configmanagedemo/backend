package model

// CreateUser 创建用户
func CreateUser(user *User) (err error) {
	result := db.Create(user)
	return result.Error
}

// UpdateUser 更新用户
func UpdateUser(user *User) (err error) {
	result := db.Save(user)
	return result.Error
}

// DeleteUser 删除用户
func DeleteUser(user *User) (err error) {
	return db.Delete(user).Error
}

// DeleteUserByUID 通过uid删除用户
func DeleteUserByUID(uid string) (err error) {
	return db.Where("uid = ?", uid).Delete(User{}).Error
}

// FindUserByUsername 根据用户名获取用户
func FindUserByUsername(username string, user *User) (err error) {
	result := db.Where(&User{Username: username}).First(&user)
	return result.Error
}

// FindUserByUID 根据uid获取用户
func FindUserByUID(uid string, user *User) (err error) {
	result := db.Where(&User{UID: uid}).First(&user)
	return result.Error
}

// GetUserList 获取用户列表
func GetUserList(offset, limit int) (total int64, users []User, err error) {
	err = db.Model(&User{}).Count(&total).Error
	if err != nil {
		return
	}
	err = db.Model(&User{}).Joins("Role").Offset(offset).Limit(limit).Find(&users).Error
	return
}
