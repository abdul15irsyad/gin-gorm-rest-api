package middlewares

import (
	"fmt"
	"gin-gorm-rest-api/lib"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogMiddleware struct {
	logger *zap.Logger
}

func NewLogMiddleware(libLogger *lib.LibLogger) *LogMiddleware {
	return &LogMiddleware{logger: libLogger.Logger}
}

func (lm *LogMiddleware) Handler(c *gin.Context) {
	c.Next()

	method := c.Request.Method
	path := c.Request.URL.Path
	status := c.Writer.Status()
	ip := c.ClientIP()

	lm.logger.Info(fmt.Sprintf("%s %s %d", method, path, status),
		zap.Any("method", method),
		zap.Any("path", path),
		zap.Any("status", status),
		zap.Any("ip", ip),
	)
}
