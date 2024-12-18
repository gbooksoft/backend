package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var routeLimits = map[string]redis_rate.Limit{}

func RegisterLimit(method, path string, perSecond int) {
	routeLimits[method+"~"+path] = redis_rate.Limit{
		Rate:   perSecond,
		Burst:  perSecond,
		Period: time.Second,
	}
}

func RateLimiterMiddleware(limiter *redis_rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем маршрут
		route := c.FullPath()
		reqMethod := c.Request.Method

		// Получаем лимит для текущего маршрута
		limit, exists := routeLimits[reqMethod+"~"+route]
		if !exists {
			// пропускаем дальше
			c.Next()
			return
		}

		// Генерация ключа для лимитера (на основе IP клиента)
		key := c.ClientIP()

		// Проверяем лимиты
		res, err := limiter.Allow(c, key, limit)
		if err != nil {
			log.Err(err).Msg("rate limiter internal err")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", limit.Burst))
		c.Header("X-Rate-Limit-Remaining", fmt.Sprintf("%d", res.Remaining))
		c.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", time.Now().Add(res.RetryAfter).Unix()))

		// Если лимит исчерпан, возвращаем 429
		if res.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		// Продолжаем выполнение запроса
		c.Next()
	}
}
