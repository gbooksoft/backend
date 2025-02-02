package middleware

import (
	"fmt"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeylimiter"
	"net/http"
	"time"
)

type rateLimit struct {
	Rate   int
	Burst  int
	Period time.Duration
}

var routeLimits = map[string]rateLimit{}

func RegisterLimit(method, path string, perSecond int) {
	routeLimits[method+"~"+path] = rateLimit{
		Rate:   perSecond,
		Burst:  perSecond,
		Period: time.Second,
	}
}

func RateLimiterMiddleware() gin.HandlerFunc {
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

		limiter, err := valkeylimiter.NewRateLimiter(valkeylimiter.RateLimiterOption{
			ClientBuilder: func(option valkey.ClientOption) (valkey.Client, error) {
				return *database.GetValkeyClient(), nil
			},
			KeyPrefix: "limiter_" + reqMethod + "~" + route,
			Limit:     limit.Rate,
			Window:    limit.Period,
		})

		if err != nil {
			log.Err(err).Msg("rate limiter internal err")
			// Продолжаем выполнение запроса, мы не можем положить приложение если редис не работает
			c.Next()
			return
		}

		key := c.ClientIP()

		res, err := limiter.Allow(c, key)
		if err != nil {
			log.Err(err).Msg("rate limiter internal err")

			c.Next()
			return
		}

		c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", limit.Burst))
		c.Header("X-Rate-Limit-Remaining", fmt.Sprintf("%d", res.Remaining))

		// Если лимит исчерпан, возвращаем 429
		if !res.Allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		// Продолжаем выполнение запроса
		c.Next()
	}
}
