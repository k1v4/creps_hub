package mongo

import (
	"article_service/internal/domain/entity"
	"article_service/pkg/DataBase/postgres"
	"context"
	"time"
)

type articleStorage struct {
	db *postgres.DB
}

func NewArticleStorage(db *postgres.DB) *articleStorage {
	return &articleStorage{
		db: db,
	}
}

func (a *articleStorage) PostArticle(ctx context.Context, name, text string, tags []string, createdAt time.Duration) (int64, error) {
	return 0, nil
}

func (a *articleStorage) GetArticle(ctx context.Context, id int64) (*entity.Article, error) {
	return nil, nil
}
