package middleware

import (
	// "embed"
	_ "embed"
	"net/http"
	opa "self_go_gin/util/open_policy_agent"
	"strings"

	"github.com/gin-gonic/gin"
)

// OpaMiddleware OPA 權限認證中間件
func OpaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		role := "guest"
		resource := "unknown"
		if strings.Contains(path, "users") {
			role = "users"
			resource = "user"
		} else if strings.Contains(path, "admins") {
			role = "admins"
			resource = "all"
		}

		var action string
		switch c.Request.Method {
		case http.MethodGet:
			action = "read"
		case http.MethodPost:
			action = "create"
		case http.MethodPut:
			action = "edit"
		case http.MethodDelete:
			action = "delete"
		default:
			action = "unknown"
		}

		result, err := opa.GetQueryResult(role, action, resource)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		// zap.S().Info("result:", result)

		// check if the user is allowed to access the resource
		if result[0].Expressions[0].Value == true {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{
			"message": "access forbidden",
		})
		c.Abort()
	}
}
