package routers

import (
	"fmt"
	"go-admin/app/admin/controllers"
	"go-admin/app/admin/middlewares"
	"go-admin/internal/lib/config"
	"go-admin/internal/lib/env"
	"go-admin/internal/lib/logger"
	"io"
	"os"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	docs "go-admin/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter init router
func InitRouter() *gin.Engine {
	// set gin
	gin.DisableConsoleColor()
	gin.SetMode(config.Settings.AdminServer.Server.GinMode)
	// set gin log to file
	gin.DefaultWriter = io.MultiWriter(logger.WriteSyncer, os.Stdout)

	// gin init
	router := gin.New()
	router.Use(gin.Logger())
	// router.Use(gin.LoggerWithWriter(logger.WriteSyncer))
	router.Use(gin.Recovery())

	// pprof
	pprof.Register(router)

	// health check
	router.GET("/health-check", controllers.HealthCheck)

	// router
	routerGroupBasePath := "/api/admin/v1"
	routerGroup := router.Group(routerGroupBasePath)
	routerGroup.Use(middlewares.TraceID())
	routerGroup.Use(middlewares.RequestLoggerMiddleware(logger.Logger))

	routerGroup.POST("/user/login", controllers.UserLogin)
	routerGroup.Use(middlewares.JWT())
	{
		routerGroup.GET("/user", controllers.UserList)
	}

	// Swagger docs
	if env.IsDev() || env.IsTest() {
		swaggerHost := fmt.Sprintf("localhost:%d", config.Settings.AdminServer.Server.HttpPort)
		if env.IsTest() {
			swaggerHost = "xxx"
		}
		docs.SwaggerInfo.Host = swaggerHost
		docs.SwaggerInfo.BasePath = routerGroupBasePath
		docs.SwaggerInfo.Title = "go-admin接口文档"
		docs.SwaggerInfo.Description = ``

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return router
}
