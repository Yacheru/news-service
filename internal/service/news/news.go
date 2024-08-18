package news

import (
	"github.com/gin-gonic/gin"
	"news-service/internal/discord"

	"news-service/internal/entities"
	"news-service/internal/repository"
)

type ServiceNews struct {
	repo    repository.News
	webhook discord.Sender
}

func NewNewsService(repo repository.News, webhook discord.Sender) *ServiceNews {
	return &ServiceNews{repo: repo, webhook: webhook}
}

func (s *ServiceNews) AddNew(ctx *gin.Context, embed *entities.Embed) (*int, error) {
	if err := s.webhook.SendEmbed(embed); err != nil {
		return nil, err
	}

	return s.repo.AddNew(ctx, embed)
}
