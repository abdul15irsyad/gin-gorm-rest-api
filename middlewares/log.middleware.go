package middlewares

import (
	"gin-gorm-rest-api/services"

	"github.com/gin-gonic/gin"
)

type LogMiddleware struct {
	logService *services.LogService
}

func NewLogMiddleware(logService *services.LogService) *LogMiddleware {
	return &LogMiddleware{logService}
}

func (lm *LogMiddleware) Log(c *gin.Context) {
	c.Next()
	lm.logService.WithFields(map[string]interface{}{
		"endpoint": c.Request.URL.String(),
		"method":   c.Request.Method,
		"status":   c.Writer.Status(),
		"ip":       c.ClientIP(),
	}, "LogMiddleware")
}
