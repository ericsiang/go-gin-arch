package middleware

import (
	"net/http"
	"self_go_gin/infra/cache/redis"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

// RateLimit 限流中間件
func RateLimit(redisLimitKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := redis_rate.NewLimiter(redis.GetRedisClient())
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
