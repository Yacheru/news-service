package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"news-service/init/config"
	"news-service/init/logger"
	"news-service/internal/repository/postgres"
	"news-service/internal/server/http/routes"
	"news-service/pkg/constants"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

func NewServer(ctx context.Context, cfg *config.Config) (*HTTPServer, error) {
	db, err := postgres.NewPostgresConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	engine := SetupGin(cfg)
	api := engine.Group(cfg.APIEntry)
	routes.NewComponentsAndRoutes(api, db).Routes()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.APIPort),
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &HTTPServer{
		server: server,
	}, nil
}

func (s *HTTPServer) Run(cfg *config.Config) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err.Error(), constants.LoggerServer)
		}
	}()

	logger.InfoF("success to listen and serve on :%d port", constants.LoggerServer, cfg.APIPort)

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func SetupGin(cfg *config.Config) *gin.Engine {
	var mode = gin.ReleaseMode
	if cfg.APIDebug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithFormatter(logger.HTTPLogger))

	return engine
}
