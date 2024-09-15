package news

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"news-service/internal/discord"
	"news-service/pkg/constants"

	"news-service/internal/entities"
	"news-service/internal/repository"
)

type ServiceNews struct {
	NewsPostgres repository.NewsPostgres
	NewsRedis    repository.NewsRedis
	webhook      discord.Sender
}

func NewNewsService(NewsPostgres repository.NewsPostgres, NewsRedis repository.NewsRedis, webhook discord.Sender) *ServiceNews {
	return &ServiceNews{
		NewsPostgres: NewsPostgres,
		NewsRedis:    NewsRedis,
		webhook:      webhook,
	}
}

func (s *ServiceNews) AddNew(ctx *gin.Context, embed *entities.Embed) (*int, error) {
	discordId, err := s.webhook.SendEmbed(embed)
	if err != nil {
		return nil, err
	}

	embed.DiscordId = discordId

	id, err := s.NewsPostgres.AddNew(ctx, discordId, embed)
	if err != nil {
		return nil, err
	}

	if err := s.NewsRedis.AddNews(ctx, discordId, embed); err != nil {
		return nil, err
	}

	return id, nil
}

func (s *ServiceNews) GetAllNews(ctx *gin.Context, params *entities.Params) ([]*entities.Embed, error) {
	_, err := s.NewsRedis.GetAllNews(ctx)
	if err != nil {
		return nil, err
	}

	news, err := s.NewsPostgres.GetAllNews(ctx)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (s *ServiceNews) GetNewsById(ctx *gin.Context, discordId int) (*entities.Embed, error) {
	embed, err := s.NewsRedis.GetNewsById(ctx, discordId)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if embed == nil || errors.Is(err, redis.Nil) {
		embed, err = s.NewsPostgres.GetNewsById(ctx, discordId)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.NoNewsFoundError
		}

		if err != nil {
			return nil, err
		}
	}

	return embed, nil
}
