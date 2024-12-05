package server

import (
	"context"
	"fmt"
	"go-admin/app/admin/routers"
	"go-admin/internal/lib/config"
	"go-admin/internal/lib/logger"
	"net/http"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	// setup logger
	logger.Setup()

	ctx := context.Background()

	// logger.Info(context.Background(), "###Server started###\n")
	logger.Info(ctx, "###Server started###\n")

	// start web server
	endless.DefaultReadTimeOut = config.Settings.AdminServer.Server.ReadTimeout
	endless.DefaultWriteTimeOut = config.Settings.AdminServer.Server.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", config.Settings.AdminServer.Server.HttpPort)
	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		logger.Info(ctx, fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(ctx, fmt.Sprintf("Server error: %v\n", err))
	}

}
