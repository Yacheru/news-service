package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"news-service/init/config"
	"news-service/internal/discord"
	"news-service/internal/repository"
	"news-service/internal/server/http/handlers"
	"news-service/internal/service"
)

type Router struct {
	router  *gin.RouterGroup
	handler *handlers.Handlers
}

func NewComponentsAndRoutes(router *gin.RouterGroup, db *sqlx.DB, cfg *config.Config) *Router {
	repo := repository.NewRepository(db)
	webhook := discord.NewWebhookClient(cfg)
	services := service.NewService(repo, webhook)
	handler := handlers.NewHandlers(services)

	return &Router{router: router, handler: handler}
}

func (r *Router) Routes() {
	{
		r.router.POST("/add", r.handler.AddNew)
		r.router.GET("/", r.handler.SearchNews)
	}
}
