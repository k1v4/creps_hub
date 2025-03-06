package articleUseCase

import (
	"article_service/internal/domain/entity"
	"context"
	"time"
)

type ArticleService interface {
	AddNewArticle(ctx context.Context, name, text string, tags []string, createdAt time.Duration) (int64, error)
	GetFullArticle(ctx context.Context, id int64) (*entity.Article, error)
}

type articleUseCase struct {
	articleService ArticleService
}

func (b *articleUseCase) GetAllForList() {
	// получить инфу о всех книгах
}

func (b *articleUseCase) PostArticle(ctx context.Context, name, text string, tags []string, createdAt time.Duration) (int64, error) {
	return b.articleService.AddNewArticle(ctx, name, text, tags, createdAt)
}

func (b *articleUseCase) GetFullArticle(ctx context.Context, id int64) (*entity.Article, error) {
	return b.articleService.GetFullArticle(ctx, id)
}
