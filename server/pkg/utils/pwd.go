package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// StringCost 字符串加密难度
var StringCost = 12

//@function: GenerateFromPassword
//@description: 字符串加密
//@param: str string
//@return: string

func GenerateFromPassword(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), StringCost)
	return string(bytes), err
}

// CheckAdminPassword 校验密码
func CheckPassword(userPwd, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPwd), []byte(password))
	return err == nil
}
