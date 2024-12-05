package handlers

import (
	"go-admin/internal/lib/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

const (
	// Common codes
	StatusCodeSuccess = 0
	StatusCodeFail    = 400

	// Business error codes
	// 1xxxx User error
	// 2xxxx Other error
	StatusCodeRegisterError = 10001
	StatusCodeLoginError    = 10002
)

type resp struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	Response(c, StatusCodeSuccess, "", data)
}

func Fail(c *gin.Context, msg string) {
	Response(c, StatusCodeFail, msg, nil)
}

func Response(c *gin.Context, code int64, message string, data interface{}) {
	ctx := c.Request.Context()

	// set header trace_id
	if ctx.Value(logger.TraceID) != nil {
		c.Header("X-Trace-ID", ctx.Value(logger.TraceID).(string))
	}

	// create a done channel to tell the request it's done
	doneChan := make(chan resp)

	// here you put the actual work needed for the request
	// and then send the doneChan with the status and body
	// to finish the request by writing the response
	go func() {
		result := resp{
			Code:    code,
			Message: message,
			Data:    data,
		}
		if code == StatusCodeSuccess && result.Message == "" {
			result.Message = "Success"
		}
		doneChan <- result
	}()

	// non-blocking select on two channels see if the request
	// times out or finishes
	select {
	// if the context is done it timed out or was cancelled
	// so don't return anything
	case <-ctx.Done():
		c.Render(StatusCodeSuccess, &render.JSON{
			Data: resp{
				Code:    StatusCodeFail,
				Message: "Fail",
				Data:    data,
			},
		})
		// if the request finished then finish the request by
		// writing the response
	case result := <-doneChan:
		c.Render(StatusCodeSuccess, &render.JSON{
			Data: result,
		})
	}

}
