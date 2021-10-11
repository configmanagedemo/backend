package model

// GetTokenList 获取Token列表
func GetTokenList(offset, limit int) (total int64, tokenList []Token, err error) {
	err = db.Model(&Token{}).Count(&total).Error
	if err != nil {
		return
	}
	err = db.Model(&Token{}).Joins("User").Offset(offset).Limit(limit).Find(&tokenList).Error
	return
}

// UpdateToken 更新token
func UpdateToken(token *Token) (err error) {
	err = db.Model(token).Save(token).Error
	return
}

// UpdateTokenByMap 更新token
func UpdateTokenByMap(id uint, data map[string]interface{}) (err error) {
	err = db.Model(&Token{}).Where("id = ?", id).Updates(data).Error
	return
}

// CreateToken 创建token
func CreateToken(token *Token) (err error) {
	err = db.Model(token).Create(token).Error
	return
}

// DeleteToken 删除token
func DeleteToken(token *Token) (err error) {
	err = db.Model(token).Delete(token).Error
	return
}
