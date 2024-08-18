package postgres

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"news-service/internal/entities"
)

type NewsPostgres struct {
	db *sqlx.DB
}

func NewNewsPostgres(db *sqlx.DB) *NewsPostgres {
	return &NewsPostgres{db: db}
}

func (p *NewsPostgres) AddNew(ctx *gin.Context, embed *entities.Embed) (*int, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
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

	var nid int
	newsQuery := `INSERT INTO news (title, description, url, timestamp, color, footer, image, thumbnail, video, author) 
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	err = p.db.GetContext(ctx.Request.Context(), &nid, newsQuery, embed.Title, embed.Description, embed.URL, embed.Timestamp, embed.Color, fid, imgid, tid, vid, aid)

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
