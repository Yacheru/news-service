package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Router struct {
	router *gin.RouterGroup
}

func NewComponentsAndRoutes(router *gin.RouterGroup, db *sqlx.DB) *Router {

	return &Router{router: router}
}

func (r *Router) Routes() {
	{
		r.router.POST("/add")
		r.router.GET("/:id")
	}
}
