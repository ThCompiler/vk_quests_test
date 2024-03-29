package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"vk_quests/internal/pkg/types"
	"vk_quests/pkg/logger"
)

const DataFormat = "2006/01/02 - 15:04:05"

const (
	RequestID logger.Field = "request_id"
	Method    logger.Field = "method"
	URL       logger.Field = "url"

	LoggerField types.ContextField = "logger"
)

func RequestLogger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		method := c.Request.Method
		requestID := uuid.New()

		if raw != "" {
			path = path + "?" + raw
		}

		lg := l.With(URL, path).With(RequestID, requestID).With(Method, method)
		c.Set(string(LoggerField), lg)

		clientIP := c.ClientIP()

		l.Info("[HTTP] Start - | %v | %s | %s  %v |",
			start.Format(DataFormat),
			clientIP,
			method,
			path,
		)

		// Process request
		c.Next()

		// Stop timer
		timeStamp := time.Now()
		latency := timeStamp.Sub(start)
		statusCode := c.Writer.Status()

		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		l.Info("[HTTP] End - %d | %v | %s | %s  %v | %v |",
			statusCode,
			timeStamp.Format(DataFormat),
			clientIP,
			method,
			path,
			latency,
		)
	}
}

func GetLogger(c *gin.Context) logger.Interface {
	if lg, ok := c.Get(string(LoggerField)); ok {
		return lg.(logger.Interface)
	}

	return logger.DefaultLogger
}
