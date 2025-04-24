package middleware

import (
	"github.com/gin-gonic/gin"
	"go-dyndns/pkg/logger"
	"time"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Info(
			"HTTP",
			"Request",
			logger.Field{Key: "method", Value: c.Request.Method},
			logger.Field{Key: "path", Value: c.Request.URL.Path},
			logger.Field{Key: "status", Value: c.Writer.Status()},
			logger.Field{Key: "client_ip", Value: c.ClientIP()},
			logger.Field{Key: "duration", Value: time.Since(start)},
			logger.Field{Key: "request_id", Value: c.GetString("RequestID")},
		)
	}
}
