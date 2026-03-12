package middleware

import (
	"errors"
	"net/http"
	"self_go_gin/util/jwt_secret"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JwtAuthMiddleware JWT 認證中間件
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var isPass = true
		var message string

		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			message = "token not found"
			isPass = false
		} else {
			if find := strings.Contains(bearerToken, "Bearer "); !find {
				message = "Bearer token fail"
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": message,
				})

				c.Abort()
				return
			}
			jwtToken := strings.Split(bearerToken, " ")
			token := jwtToken[1]

			claims, err := jwt_secret.ParseToken(token)
			if claims != nil {
				if claims.UserID == 0 && claims.AdminID == 0 {
					message = "jwt data fail"
					isPass = false
				}
				if claims.UserID != 0 {
					c.Set("usersID", claims.UserID)
				} else if claims.AdminID != 0 {
					c.Set("adminID", claims.AdminID)
				}
			}

			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					message = "token expired"
					isPass = false
				default:
					message = "token fail"
					isPass = false
				}
			}
		}

		if !isPass {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": message,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
