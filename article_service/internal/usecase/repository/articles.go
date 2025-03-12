package repository

import (
	"article_service/internal/entity"
	"article_service/pkg/DataBase/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type ArticleRepository struct {
	*postgres.Postgres
}

func NewArticleRepository(pg *postgres.Postgres) *ArticleRepository {
	return &ArticleRepository{pg}
}

func (a *ArticleRepository) AddArticle(ctx context.Context, authorId int, name string, content string) (int, error) {
	const op = "ArticleRepository.AddArticle"

	s, args, err := a.Builder.Insert("articles").
		Columns("author_id", "name", "text").
		Values(authorId, name, content).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = a.Pool.QueryRow(ctx, s, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *ArticleRepository) FindArticleByID(ctx context.Context, id int) (entity.Article, error) {
	const op = "ArticleRepository.FindArticleByID"

	s, args, err := a.Builder.Select("*").
		From("articles").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	var article entity.Article
	err = a.Pool.QueryRow(ctx, s, args...).
		Scan(&article.ID,
			&article.AuthorID,
			&article.PublicationDate,
			&article.Name,
			&article.Text)
	if err != nil {
		return entity.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

func (a *ArticleRepository) FindAllArticle(ctx context.Context, limit, offset uint64) ([]entity.Article, error) {
	const op = "ArticleRepository.FindAllArticle"

	s, args, err := a.Builder.Select("*").
		From("articles").
		OrderBy("publication_date DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := a.Pool.Query(ctx, s, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var articles []entity.Article
	for rows.Next() {
		var article entity.Article
		err = rows.Scan(&article.ID, &article.AuthorID, &article.PublicationDate, &article.Name, &article.Text)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (a *ArticleRepository) DeleteArticle(ctx context.Context, id int) error {
	const op = "ArticleRepository.DeleteArticle"

	s, args, err := a.Builder.Delete("articles").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = a.Pool.Exec(ctx, s, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// TODO надо ли разрешать редактирование
//func (a *ArticleRepository) UpdateArticle() {}
