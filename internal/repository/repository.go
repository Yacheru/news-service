package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"news-service/internal/entities"
	"news-service/internal/repository/postgres"
)

type News interface {
	AddNew(ctx *gin.Context, embed *entities.Embed) (*int, error)
}

type Repository struct {
	News
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		News: postgres.NewNewsPostgres(db),
	}
}
