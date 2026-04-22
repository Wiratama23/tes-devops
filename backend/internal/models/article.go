package models

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ArticlesID  int       `json:"articles_id"`
	UID         uuid.UUID `json:"uid"`
	Title       string    `json:"title"`
	ArticleText string    `json:"article_text"`
	DateCreated time.Time `json:"date_created"`
	UpdatedAt   time.Time `json:"updated_at"`
}
