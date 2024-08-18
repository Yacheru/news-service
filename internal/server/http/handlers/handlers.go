package handlers

import (
	"net/http"
	"news-service/internal/service"

	"github.com/gin-gonic/gin"

	"news-service/init/logger"
	"news-service/internal/entities"
	"news-service/pkg/constants"
)

type Handlers struct {
	service service.News
}

func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) AddNew(ctx *gin.Context) {
	var embedEntity = new(entities.Embed)
	if err := ctx.ShouldBindJSON(embedEntity); err != nil {
		logger.Error(err.Error(), constants.LoggerHandlers)

		NewErrorResponse(ctx, http.StatusBadRequest, constants.BodyIsInvalid)
		return
	}

	id, err := h.service.AddNew(ctx, embedEntity)
	if err != nil {
		logger.Error(err.Error(), constants.LoggerHandlers)

		NewErrorResponse(ctx, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "news added successfully", id)
}

func (h *Handlers) SearchNews(ctx *gin.Context) {
	// TODO:
	/*
		- Search by title
		- Search by description
		- Search by author
		- Search by footer
	*/
}
