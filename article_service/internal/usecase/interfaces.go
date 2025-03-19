package usecase

import (
	"article_service/internal/entity"
	"context"
)

type IArticleRepository interface {
	AddArticle(ctx context.Context, authorId int, name string, content string) (int, error)
	FindArticleByID(ctx context.Context, id int) (entity.Article, error)
	FindAllArticle(ctx context.Context, limit, offset uint64) ([]entity.ArticleUser, error)
	DeleteArticle(ctx context.Context, articleId, authorId int) error
	FindAllArticlesByUser(ctx context.Context, userId int) ([]entity.Article, error)
}

type IArticleService interface {
	AddArticle(ctx context.Context, authorId int, title, content string) (int, error)
	FindArticle(ctx context.Context, id int) (entity.Article, error)
	DeleteArticle(ctx context.Context, articleId, authorId int) (bool, error)
	FindAllArticle(ctx context.Context, limit, offset int) ([]entity.ArticleUser, error)
	FindAllArticleByUser(ctx context.Context, userId int) ([]entity.Article, error)
}
