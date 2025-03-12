package usecase

import (
	"article_service/internal/entity"
	"context"
)

type IArticleRepository interface {
	AddArticle(ctx context.Context, authorId int, name string, content string) (int, error)
	FindArticleByID(ctx context.Context, id int) (entity.Article, error)
	FindAllArticle(ctx context.Context, limit, offset uint64) ([]entity.Article, error)
	DeleteArticle(ctx context.Context, id int) error
}

type IArticleService interface {
	AddArticle(ctx context.Context, authorId int, title, content string) (int, error)
	FindArticle(ctx context.Context, id int) (entity.Article, error)
	DeleteArticle(ctx context.Context, id int) (bool, error)
	FindAllArticle(ctx context.Context, limit, offset int) ([]entity.Article, error)
}
