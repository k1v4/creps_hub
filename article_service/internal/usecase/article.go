package usecase

import (
	"article_service/internal/entity"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	uploaderv1 "github.com/k1v4/protos/gen/file_uploader"
	"github.com/redis/go-redis/v9"
	"regexp"
	"strings"
	"time"
)

type ArticleUseCase struct {
	repo   IArticleRepository
	client uploaderv1.FileUploaderClient
	cache  *redis.Client
}

func NewArticleUseCase(repo IArticleRepository, client uploaderv1.FileUploaderClient, cache *redis.Client) *ArticleUseCase {
	return &ArticleUseCase{
		repo:   repo,
		client: client,
		cache:  cache,
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

	processContent, err := a.processImages(ctx, content)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	articleId, err := a.repo.AddArticle(ctx, authorId, title, processContent, image.GetUrl())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return articleId, nil
}

func (a *ArticleUseCase) FindArticle(ctx context.Context, id int) (entity.Article, error) {
	const op = "ArticleUseCase.FindArticle"
	var res entity.Article

	err := a.cache.Get(ctx, fmt.Sprintf("%d", id)).Scan(&res)
	if err == nil {
		return res, nil
	} else if err != nil && !errors.Is(err, redis.Nil) {
		return entity.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	articleByID, err := a.repo.FindArticleByID(ctx, id)
	if err != nil {
		return entity.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	set := a.cache.Set(ctx, fmt.Sprintf("%d", id), &articleByID, 15*time.Minute)
	if set.Err() != nil {
		return res, fmt.Errorf("%s: %w", op, set.Err())
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

func (a *ArticleUseCase) processImages(ctx context.Context, content string) (string, error) {
	const op = "usecase.processImages"

	re := regexp.MustCompile(`data:image/(.*?);base64,(.*?)"`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		// Извлекаем тип изображения и base64-данные
		imageType := match[1]
		base64Data := match[2]

		// Генерируем уникальное имя файла
		uniqueFileName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), imageType)

		imageData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return "", fmt.Errorf("ошибка декодирования base64: %v", err)
		}

		image, err := a.client.UploadFile(ctx, &uploaderv1.ImageUploadRequest{
			ImageData: imageData,
			FileName:  uniqueFileName,
		})
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		// Заменяем base64 на URL изображения
		content = strings.Replace(content, match[0], fmt.Sprintf(`%s"/`, image.GetUrl()), 1)

	}

	return content, nil
}
