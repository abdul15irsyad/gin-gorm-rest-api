package middlewares

import (
	"gin-gorm-rest-api/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorMiddleware struct {
	logger *zap.Logger
}

func NewErrorMiddleware(libLogger *lib.LibLogger) *ErrorMiddleware {
	return &ErrorMiddleware{libLogger.Logger}
}

func (em *ErrorMiddleware) Handler(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			em.logger.Error("something went wrong",
				zap.Any("error", err),
			)

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
		}
	}()

	ctx.Next()
}
