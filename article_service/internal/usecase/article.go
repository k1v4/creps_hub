package usecase

import (
	"article_service/internal/entity"
	"context"
	"fmt"
	uploaderv1 "github.com/k1v4/protos/gen/file_uploader"
)

type ArticleUseCase struct {
	repo   IArticleRepository
	client uploaderv1.FileUploaderClient
}

func NewArticleUseCase(repo IArticleRepository, client uploaderv1.FileUploaderClient) *ArticleUseCase {
	return &ArticleUseCase{
		repo:   repo,
		client: client,
	}
}

func (a *ArticleUseCase) AddArticle(ctx context.Context, authorId int, title, content, imageName string, imageData []byte) (int, error) {
	const op = "ArticleUseCase.AddArticle"

	image, err := a.client.UploadFile(ctx, &uploaderv1.ImageUploadRequest{
		ImageData: imageData,
		FileName:  imageName,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	articleId, err := a.repo.AddArticle(ctx, authorId, title, content, image.GetUrl())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return articleId, nil
}

func (a *ArticleUseCase) FindArticle(ctx context.Context, id int) (entity.Article, error) {
	const op = "ArticleUseCase.FindArticle"

	articleByID, err := a.repo.FindArticleByID(ctx, id)
	if err != nil {
		return entity.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return articleByID, nil
}

func (a *ArticleUseCase) DeleteArticle(ctx context.Context, articleId, authorId int) (bool, error) {
	const op = "ArticleUseCase.DeleteArticle"

	err := a.repo.DeleteArticle(ctx, articleId, authorId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (a *ArticleUseCase) FindAllArticle(ctx context.Context, limit, offset int) ([]entity.ArticleUser, error) {
	const op = "ArticleUseCase.FindAllArticle"

	allArticle, err := a.repo.FindAllArticle(ctx, uint64(limit), uint64(offset))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return allArticle, nil
}

func (a *ArticleUseCase) FindAllArticleByUser(ctx context.Context, userId int) ([]entity.Article, error) {
	const op = "ArticleUseCase.FindAllArticle"

	allArticle, err := a.repo.FindAllArticlesByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return allArticle, nil
}
