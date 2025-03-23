package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	// bcrypt.DefaultCost 工作因子 迭代次数
	// 工作因子越大，密码越复杂，安全性越高，但是加密时间越长
	// 得到加密后的密码 哈希值
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("encryptPassword error = %v\n", err)
		return "", err
	}
	return string(hash), nil
}
