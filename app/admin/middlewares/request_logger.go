package middlewares

import (
	"bytes"
	"encoding/json"
	"go-admin/internal/lib/logger"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JSON 简单脱敏
func sanitizeJSON(input []byte, sensitiveKeys []string) []byte {
	var data map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return input
	}

	for _, key := range sensitiveKeys {
		if _, exists := data[key]; exists {
			data[key] = "***"
		}
	}

	output, _ := json.Marshal(data)
	return output
}

// JSON 深层嵌套脱敏
func sanitizeNestedJSON(data map[string]interface{}, sensitiveKeys []string) {
	for key, value := range data {
		// 如果是敏感字段，替换为 "***"
		for _, sensitiveKey := range sensitiveKeys {
			if key == sensitiveKey {
				data[key] = "***"
			}
		}

		// 如果字段是嵌套的 map，递归处理
		if nestedMap, ok := value.(map[string]interface{}); ok {
			sanitizeNestedJSON(nestedMap, sensitiveKeys)
		}

		// 如果字段是嵌套的数组，递归处理数组中的每个元素
		if nestedArray, ok := value.([]interface{}); ok {
			for _, item := range nestedArray {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					sanitizeNestedJSON(nestedMap, sensitiveKeys)
				}
			}
		}
	}
}

func sanitizeJSONWithNesting(input []byte, sensitiveKeys []string) []byte {
	var data map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return input
	}

	sanitizeNestedJSON(data, sensitiveKeys)

	output, _ := json.Marshal(data)
	return output
}

// RequestLoggerMiddleware 是记录请求详细信息的中间件
func RequestLoggerMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		start := time.Now()

		// 记录请求的基本信息
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// 记录查询参数
		queryParams := c.Request.URL.Query()

		// 记录请求体（支持 JSON）
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// 读取后需要重新设置请求体，否则后续处理器无法读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 脱敏处理
		sensitiveKeys := []string{"password", "token"}
		// sanitizedBodyBytes := sanitizeJSON(bodyBytes, sensitiveKeys)
		sanitizedBodyBytes := sanitizeJSONWithNesting(bodyBytes, sensitiveKeys)

		// 调用后续处理器
		c.Next()

		// 请求结束时间
		latency := time.Since(start)

		// 记录响应状态码
		statusCode := c.Writer.Status()

		// 使用 zap 记录日志（开发环境）
		zapLogger.Debug("HTTP Request",
			zap.String("trace_id", c.Request.Context().Value(logger.TraceID).(string)),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Any("query_params", queryParams),
			zap.ByteString("body", sanitizedBodyBytes),
			zap.Int("status_code", statusCode),
			zap.Duration("latency", latency),
		)
	}
}
