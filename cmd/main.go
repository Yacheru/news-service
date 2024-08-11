package main

import (
	"context"
	"os/signal"
	"syscall"

	"news-service/init/config"
	"news-service/init/logger"
	"news-service/internal/server"
	"news-service/pkg/constants"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cfg := &config.ServerConfig

	logger.InitLogger(cfg.APIDebug)

	if err := config.InitConfig(); err != nil {
		logger.Error(err.Error(), constants.LoggerCMD)

		cancel()
	}
	logger.Info("configuration loaded", constants.LoggerCMD)

	app, err := server.NewServer(ctx, cfg)
	if err != nil {
		logger.Error(err.Error(), constants.LoggerCMD)

		cancel()
	}
	logger.Info("server configured", constants.LoggerCMD)

	if err := app.Run(cfg); err != nil {
		logger.Error(err.Error(), constants.LoggerCMD)

		cancel()
	}
	logger.Info("server is running", constants.LoggerCMD)

	<-ctx.Done()

	if app != nil {
		if err := app.Shutdown(ctx); err != nil {
			logger.Error(err.Error(), constants.LoggerCMD)
		}
		
		logger.Info("http-server shutdown", constants.LoggerCMD)
	}

	logger.Info("service shutdown", constants.LoggerCMD)
}
