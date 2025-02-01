package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func Add(e *gin.Engine) {
	e.Use(middleware.RateLimiterMiddleware(redis_rate.NewLimiter(database.GetRedisClient())))

	initializeDependencies()

	authRoutes(e)
	accountRoutes(e)
	projectRoutes(e)
}
