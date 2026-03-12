package middleware

import (
	"net/http"
	opa "self_go_gin/util/open_policy_agent"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OpaMiddleware OPA 權限認證中間件
func OpaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := opa.GetQueryResult(c)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		zap.S().Info("result:", result)

		// check if the user is allowed to access the resource
		if result[0].Expressions[0].Value == true {
			c.Next()
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "access forbidden",
			})
			c.Abort()
			return
		}
	}
}
