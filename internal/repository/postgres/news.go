package postgres

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"news-service/init/logger"
	"news-service/internal/entities"
	"news-service/pkg/constants"
	"time"
)

type NewsPostgres struct {
	db *sqlx.DB
}

func NewNewsPostgres(db *sqlx.DB) *NewsPostgres {
	return &NewsPostgres{db: db}
}

func (p *NewsPostgres) AddNew(ctx *gin.Context, discordId int, embed *entities.Embed) (*int, error) {
	tx, err := p.db.BeginTx(ctx.Request.Context(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}

	var fid *int
	if embed.Footer != nil {
		footerQuery := `INSERT INTO footers (text, icon_url) VALUES ($1, $2) RETURNING id`
		err = p.db.GetContext(ctx.Request.Context(), &fid, footerQuery, embed.Footer.Text, embed.Footer.IconURL)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var imgid *int
	if embed.Image != nil {
		imageQuery := `INSERT INTO resources (url) VALUES ($1) RETURNING id`
		err = p.db.GetContext(ctx.Request.Context(), &imgid, imageQuery, embed.Image.URL)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var tid *int
	if embed.Thumbnail != nil {
		thumbnailQuery := `INSERT INTO resources (url) VALUES ($1) RETURNING id`
		err = p.db.GetContext(ctx.Request.Context(), &tid, thumbnailQuery, embed.Thumbnail.URL)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var vid *int
	if embed.Video != nil {
		videoQuery := `INSERT INTO resources (url) VALUES ($1) RETURNING id`
		err = p.db.GetContext(ctx.Request.Context(), &vid, videoQuery, embed.Video.URL)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var aid *int
	if embed.Author != nil {
		authorQuery := `INSERT INTO authors (name, url, icon_url) VALUES ($1, $2, $3) RETURNING id`
		err = p.db.GetContext(ctx.Request.Context(), &aid, authorQuery, embed.Author.Name, embed.Author.URL, embed.Author.IconURL)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	fids := make([]int, 0, len(embed.Fields))
	if embed.Fields != nil {
		var fid int
		fieldQuery := `INSERT INTO fields (name, value, inline) VALUES ($1, $2, $3) RETURNING id`
		for _, field := range embed.Fields {
			err = p.db.GetContext(ctx.Request.Context(), &fid, fieldQuery, field.Name, field.Value, field.Inline)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			fids = append(fids, fid)
		}
	}

	if embed.Timestamp == nil {
		now := time.Now().UTC()
		embed.Timestamp = &now
	}

	var nid int
	newsQuery := `INSERT INTO news (discord_id, title, description, url, color, footer, image, thumbnail, video, author, created_at) 
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	err = p.db.GetContext(ctx.Request.Context(), &nid, newsQuery, discordId, embed.Title, embed.Description, embed.URL, embed.Color, fid, imgid, tid, vid, aid, embed.Timestamp)
	if err != nil {
		return nil, err
	}

	if len(fids) != 0 {
		newsFieldsQuery := `INSERT INTO news_fields (news_id, field_id) VALUES ($1, $2)`
		for _, fid := range fids {
			_, err = p.db.ExecContext(ctx.Request.Context(), newsFieldsQuery, nid, fid)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	return &nid, tx.Commit()
}

//query := `
//	SELECT n.title, n.description, n.url, n.timestamp, n.color,
//	    f.text, f.icon_url AS footer,
//	    ri.url AS image,
//	    rt.url AS thumbnail,
//	    rv.url AS video,
//	    a.name, a.url, a.icon_url AS author
//	FROM news n
//		LEFT JOIN footers f ON f.id = n.footer
//		LEFT JOIN resources ri ON ri.id = n.image
//		LEFT JOIN resources rt ON rt.id = n.thumbnail
//		LEFT JOIN resources rv ON rv.id = n.video
//		LEFT JOIN authors a ON a.id = n.author
//	WHERE n.description
//	LIKE '%' || $1 || '%'
//	OR n.title
//	LIKE '%' || $2 || '%'
//`

func (p *NewsPostgres) GetAllNews(ctx *gin.Context) ([]*entities.Embed, error) {
	var entityEmbeds []*entities.Embed
	fieldMap := make(map[int][]entities.EmbedField)

	query := `
		SELECT 
		    n.discord_id, n.title, n.description, n.url, n.created_at, n.color,
		    f.text AS footer_text, f.icon_url AS footer_icon_url, 
		    ri.url AS image_url,
		    rt.url AS thumbnail_url,
		    rv.url AS video_url,
		    a.name AS author_name, a.url AS author_url, a.icon_url AS author_icon_url,
		FROM news n
			LEFT JOIN footers f ON f.id = n.footer
			LEFT JOIN resources ri ON ri.id = n.image
			LEFT JOIN resources rt ON rt.id = n.thumbnail
			LEFT JOIN resources rv ON rv.id = n.video
			LEFT JOIN authors a ON a.id = n.author
			LEFT JOIN news_fields nf ON nf.news_id = n.id
	`

	rows, err := p.db.QueryContext(ctx.Request.Context(), query)
	if err != nil {
		logger.Error(err.Error(), constants.LoggerPostgres)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var embed entities.Embed
		var footer entities.EmbedFooter
		var image, thumbnail, video entities.EmbedResource
		var author entities.EmbedAuthor

		var fieldName, fieldValue string
		var fieldInline sql.NullBool

		err = rows.Scan(
			&embed.DiscordId, &embed.Title, &embed.Description, &embed.URL, &embed.Timestamp, &embed.Color,
			&footer.Text, &footer.IconURL,
			&image.URL, &thumbnail.URL, &video.URL,
			&author.Name, &author.URL, &author.IconURL,
			&fieldName, &fieldValue, &fieldInline,
		)
		if err != nil {
			logger.Error(err.Error(), constants.LoggerPostgres)
			return nil, err
		}

		embed.Footer = &footer
		embed.Image = &image
		embed.Thumbnail = &thumbnail
		embed.Video = &video
		embed.Author = &author

		if len(fieldName) > 0 {
			field := entities.EmbedField{
				Name:   fieldName,
				Value:  fieldValue,
				Inline: &fieldInline.Bool,
			}
			fieldMap[embed.DiscordId] = append(fieldMap[embed.DiscordId], field)
		}

		if rows.Err() != nil {
			logger.Error(rows.Err().Error(), constants.LoggerPostgres)
			return nil, err
		}
	}

	for _, embed := range entityEmbeds {
		embed.Fields = fieldMap[embed.DiscordId]
	}

	return entityEmbeds, nil
}

func (p *NewsPostgres) GetNewsById(ctx *gin.Context, discordId int) (*entities.Embed, error) {
	query := `
		SELECT 
		    n.discord_id, n.title, n.description, n.url, n.created_at, n.color,
		    f.text AS footer_text, f.icon_url AS footer_icon_url, 
		    ri.url AS image_url,
		    rt.url AS thumbnail_url,
		    rv.url AS video_url,
		    a.name AS author_name, a.url AS author_url, a.icon_url AS author_icon_url
		FROM news n
			LEFT JOIN footers f ON f.id = n.footer
			LEFT JOIN resources ri ON ri.id = n.image
			LEFT JOIN resources rt ON rt.id = n.thumbnail
			LEFT JOIN resources rv ON rv.id = n.video
			LEFT JOIN authors a ON a.id = n.author
			LEFT JOIN news_fields nf ON nf.news_id = n.id
		WHERE discord_id = $1
	`

	row := p.db.QueryRowxContext(ctx.Request.Context(), query, discordId)
	if row.Err() != nil {
		logger.Error(row.Err().Error(), constants.LoggerPostgres)
		return nil, row.Err()
	}

	var embed = new(entities.Embed)
	var footer entities.EmbedFooter
	var image, thumbnail, video entities.EmbedResource
	var author entities.EmbedAuthor

	err := row.Scan(
		&embed.DiscordId, &embed.Title, &embed.Description, &embed.URL, &embed.Timestamp, &embed.Color,
		&footer.Text, &footer.IconURL,
		&image.URL, &thumbnail.URL, &video.URL,
		&author.Name, &author.URL, &author.IconURL,
	)
	if err != nil {
		logger.Error(err.Error(), constants.LoggerPostgres)
		return nil, err
	}

	embed.Footer = &footer
	embed.Image = &image
	embed.Thumbnail = &thumbnail
	embed.Video = &video
	embed.Author = &author

	return embed, nil
}
