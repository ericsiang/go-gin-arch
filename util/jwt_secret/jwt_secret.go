package jwt_secret

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// LoginRole 登入角色
type LoginRole string

// login role
const (
	// 	LoginUser 用戶角色
	LoginUser  LoginRole = "user"
	// LoginAdmin 管理員角色
	LoginAdmin LoginRole = "admin"
)

// Claims 用戶聲明
type Claims struct {
	UserID  uint `json:"user_id"`
	AdminID uint `json:"admin_id"`
	jwt.MapClaims
}

// SetJwtSecret 設置JWT密鑰
func SetJwtSecret(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateToken 根據用戶的用戶id 生成JWT token
func GenerateToken(checkLoginRole LoginRole, loginID uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(1 * time.Hour)
	claims := Claims{}
	switch checkLoginRole {
	case LoginUser:
		claims.UserID = loginID
	case LoginAdmin:
		claims.AdminID = loginID
	default:
		return "", errors.New("LoginRole not allow")
	}
	claims.MapClaims = jwt.MapClaims{
		"exp": expireTime.Unix(),
		"iss": "gin-blog",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret) //簽名字串
	return token, err

}

// ParseToken 解析JWT token
func ParseToken(token string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
