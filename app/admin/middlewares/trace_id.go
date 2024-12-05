package middlewares

import (
	"context"
	"go-admin/internal/lib/logger"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := shortuuid.New()

		// set uuid to http header
		c.Writer.Header().Set("X-Trace-ID", uuid)

		// set uuid to gin context
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), logger.TraceID, uuid))

		// get trace id
		// fmt.Printf("trace_id: %v\n", c.Request.Context().Value(logger.TraceID))

		c.Next()
	}
}
