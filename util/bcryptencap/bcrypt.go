package bcryptencap

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateFromPassword 生成密碼哈希值
func GenerateFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CompareHashAndPassword 比較哈希值和密碼是否匹配
func CompareHashAndPassword(realPassword, checkPassword []byte) error {
	return bcrypt.CompareHashAndPassword(realPassword, checkPassword)
}
