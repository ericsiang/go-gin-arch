package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"net/http"
	"self_go_gin/container"
)

// RateLimit 限流中間件
func RateLimit(redisLimitKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		redisClient := container.GetContainer().GetRedisClient()
		limiter := redis_rate.NewLimiter(redisClient)
		// 限制每秒 5 個 request
		res, err := limiter.Allow(c, redisLimitKey, redis_rate.PerSecond(5))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if res.Allowed == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"msg": "Too many requests",
			})
			return
		}

		c.Next()
	}

}
