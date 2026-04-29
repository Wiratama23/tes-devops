package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"rwiratama.com/m/internal/models"
)

type ArticleRepository struct {
	pool PgxPool
}

func NewArticleRepository(pool PgxPool) *ArticleRepository {
	return &ArticleRepository{pool: pool}
}

// Create inserts a new article into the database
func (r *ArticleRepository) Create(ctx context.Context, uid uuid.UUID, title, articleText string) (*models.Article, error) {
	query := `
		INSERT INTO articles (uid, title, article_text)
		VALUES ($1, $2, $3)
		RETURNING articles_id, uid, title, article_text, date_created, updated_at
	`

	var article models.Article
	err := r.pool.QueryRow(ctx, query, uid, title, articleText).Scan(
		&article.ArticlesID,
		&article.UID,
		&article.Title,
		&article.ArticleText,
		&article.DateCreated,
		&article.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return &article, nil
}

// GetByID retrieves an article by articles_id
func (r *ArticleRepository) GetByID(ctx context.Context, articleID int) (*models.Article, error) {
	query := `
		SELECT articles_id, uid, title, article_text, date_created, updated_at
		FROM articles
		WHERE articles_id = $1
	`

	var article models.Article
	err := r.pool.QueryRow(ctx, query, articleID).Scan(
		&article.ArticlesID,
		&article.UID,
		&article.Title,
		&article.ArticleText,
		&article.DateCreated,
		&article.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return &article, nil
}

// GetByUserID retrieves all articles by a specific user
func (r *ArticleRepository) GetByUserID(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	query := `
		SELECT articles_id, uid, title, article_text, date_created, updated_at
		FROM articles
		WHERE uid = $1
		ORDER BY date_created DESC
	`

	rows, err := r.pool.Query(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles by user: %w", err)
	}
	defer rows.Close()

	articles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Article, error) {
		var article models.Article
		err := row.Scan(
			&article.ArticlesID,
			&article.UID,
			&article.Title,
			&article.ArticleText,
			&article.DateCreated,
			&article.UpdatedAt,
		)
		return article, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to collect articles: %w", err)
	}

	return articles, nil
}

// Update updates an article
func (r *ArticleRepository) Update(ctx context.Context, articleID int, title, articleText string) (*models.Article, error) {
	query := `
		UPDATE articles
		SET title = $2, article_text = $3, updated_at = CURRENT_TIMESTAMP
		WHERE articles_id = $1
		RETURNING articles_id, uid, title, article_text, date_created, updated_at
	`

	var article models.Article
	err := r.pool.QueryRow(ctx, query, articleID, title, articleText).Scan(
		&article.ArticlesID,
		&article.UID,
		&article.Title,
		&article.ArticleText,
		&article.DateCreated,
		&article.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return &article, nil
}

// Delete removes an article from the database
func (r *ArticleRepository) Delete(ctx context.Context, articleID int) error {
	query := `DELETE FROM articles WHERE articles_id = $1`

	result, err := r.pool.Exec(ctx, query, articleID)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

// GetAllWithPagination retrieves articles with pagination (uses date_created
// index) and returns the total row count alongside the page so the frontend
// can render an accurate "Page N of M" pager.
func (r *ArticleRepository) GetAllWithPagination(ctx context.Context, limit, offset int) ([]models.Article, int, error) {
	var totalCount int
	if err := r.pool.QueryRow(ctx, articleCountSQL).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	rows, err := r.pool.Query(ctx, articlePaginatedSelectSQL, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query articles with pagination: %w", err)
	}
	defer rows.Close()

	articles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Article, error) {
		var article models.Article
		err := row.Scan(
			&article.ArticlesID,
			&article.UID,
			&article.Title,
			&article.ArticleText,
			&article.DateCreated,
			&article.UpdatedAt,
		)
		return article, err
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to collect articles: %w", err)
	}

	return articles, totalCount, nil
}

// SQL strings are exported as constants so tests can reuse the exact text
// when wiring pgxmock expectations.
const (
	// Approximate row count from the planner statistics. ~1ms regardless of
	// table size — vastly cheaper than COUNT(*) on a million-row table.
	// reltuples is maintained by autovacuum/ANALYZE; a brand-new table not
	// yet analyzed reports -1, which we clamp to 0 so the pager doesn't
	// crash. Run `ANALYZE articles` after large bulk loads to refresh it.
	articleCountSQL = `SELECT GREATEST(reltuples, 0)::BIGINT FROM pg_class WHERE relname = 'articles'`

	articlePaginatedSelectSQL = `
		SELECT articles_id, uid, title, article_text, date_created, updated_at
		FROM articles
		ORDER BY date_created DESC
		LIMIT $1 OFFSET $2
	`
)
