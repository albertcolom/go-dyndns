package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		fullURL := c.Request.URL.RequestURI()
		clientIP := c.ClientIP()
		requestID := c.GetString("RequestID")

		log.Printf("[HTTP] %s %s | %d | %s | %s | %s",
			method, fullURL, statusCode, duration, clientIP, requestID,
		)
	}
}
