package repository

import (
	"context"
	"fmt"
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/errors"
	"time"
)

type BlogRepository interface {
	// Create and Update
	CreateArticle(ctx context.Context, data *domain.BlogArtikel) error
	UpdateArticle(ctx context.Context, data *domain.BlogArtikel) error
	UpdateArticleStatus(ctx context.Context, id string, status constants.ArticleStatus, publishAt *time.Time) error
	IncrementViews(ctx context.Context, id string) error
	SlugExists(ctx context.Context, slug string) (bool, error)

	// Read operations
	GetArticle(ctx context.Context, id string) (domain.BlogArtikel, error)
	GetArticleBySlug(ctx context.Context, slug string) (domain.BlogArtikel, error)
	ListArticles(ctx context.Context, req requests.ListArtikel) ([]domain.BlogArtikel, int, error)
	GetArticleStats(ctx context.Context) (dto.BlogStats, error)

	// Public endpoints
	ListPublicArticles(ctx context.Context, req requests.ListArtikel) ([]domain.BlogArtikel, int, error)
	GetPublicArticleWithRelated(ctx context.Context, slug string) (article domain.BlogArtikel, related []domain.BlogArtikel, err error)

	// Tag related operations
	AddArticleTags(ctx context.Context, articleID string, tagIDs []string) error
	RemoveArticleTags(ctx context.Context, articleID string, tagIDs []string) error
	ClearArticleTags(ctx context.Context, articleID string) error

	// Delete operations
	DeleteArticle(ctx context.Context, id string) error
	HardDeleteArticle(ctx context.Context, id string) error
}

type blogRepository struct {
	db *database.Database
}

func NewBlogRepository(db *database.Database) BlogRepository {
	return &blogRepository{
		db: db,
	}
}

func (r *blogRepository) CreateArticle(ctx context.Context, data *domain.BlogArtikel) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	if err != nil {
		return errors.CheckUniqueViolation(err)
	}
	return err
}

func (r *blogRepository) UpdateArticle(ctx context.Context, data *domain.BlogArtikel) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(data).
		Where("id = ?", data.ID).
		ExcludeColumn("id", "created_at", "views").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *blogRepository) UpdateArticleStatus(ctx context.Context, id string, status constants.ArticleStatus, publishAt *time.Time) error {
	query := r.db.InitQuery(ctx).NewUpdate().
		Table("blog_artikels").
		Set("status = ?", status).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id)

	if publishAt != nil {
		query.Set("published_at = ?", publishAt)
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *blogRepository) IncrementViews(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Table("blog_artikels").
		Set("views = views + 1").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *blogRepository) GetArticle(ctx context.Context, id string) (res domain.BlogArtikel, err error) {
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Relation("Author").
		Relation("Tags").
		Where(`"ba"."id" = ?`, id).
		Scan(ctx)
	return res, err
}

func (r *blogRepository) GetArticleBySlug(ctx context.Context, slug string) (res domain.BlogArtikel, err error) {
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Relation("Author").
		Relation("Tags").
		Where(`"ba"."slug" = ?`, slug).
		Scan(ctx)
	return res, err
}

func (r *blogRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	return r.db.InitQuery(ctx).
		NewSelect().
		Model((*domain.BlogArtikel)(nil)).
		Where("slug = ?", slug).
		Exists(ctx)
}

func (r *blogRepository) ListArticles(ctx context.Context, req requests.ListArtikel) ([]domain.BlogArtikel, int, error) {
	var res []domain.BlogArtikel

	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Relation("Author").
		Relation("Tags")

	// Apply filters
	if req.CategoryID != "" {
		q.Where("category_id = ?", req.CategoryID)
	}
	if req.TagID != "" {
		q.Join("JOIN article_tags at ON at.blog_article_id = blog_artikels.id").
			Where("at.tag_id = ?", req.TagID)
	}
	if req.Status != "" {
		q.Where("ba.status = ?", req.Status)
	}
	if req.Search != "" {
		q.Where("ba.title ILIKE ? OR ba.content ILIKE ?",
			fmt.Sprintf("%%%s%%", req.Search),
			fmt.Sprintf("%%%s%%", req.Search))
	}
	if req.StartDate != nil {
		q.Where("ba.created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		q.Where("ba.created_at <= ?", req.EndDate)
	}

	// Apply sorting
	order := "DESC"
	orderBy := "created_at"
	if req.SortBy != "" {
		orderBy = req.SortBy
		if req.SortOrder == "asc" {
			order = "ASC"
		}
	}
	q.Order(fmt.Sprintf("%s %s", orderBy, order))

	// Apply pagination
	q.Limit(req.PageSize).
		Offset(req.CalculateOffset())

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}

func (r *blogRepository) GetArticleStats(ctx context.Context) (res dto.BlogStats, err error) {
	err = r.db.InitQuery(ctx).
		NewSelect().
		ColumnExpr("COUNT(*) as total_articles").
		ColumnExpr("COUNT(CASE WHEN status = 'published' THEN 1 END) as published_articles").
		ColumnExpr("COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft_articles").
		ColumnExpr("SUM(views) as total_views").
		Table("blog_artikels").
		Where("deleted_at IS NULL").
		Scan(ctx, &res)
	return res, err
}

func (r *blogRepository) AddArticleTags(ctx context.Context, articleID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	values := make([]interface{}, 0, len(tagIDs)*2)
	placeholders := make([]string, 0, len(tagIDs))

	for _, tagID := range tagIDs {
		values = append(values, articleID, tagID)
		placeholders = append(placeholders, "(?, ?)")
	}

	_, err := r.db.InitQuery(ctx).NewRaw(`
		INSERT INTO article_tags (blog_article_id, tag_id)
		VALUES `+placeholders[0], values...).
		Exec(ctx)

	return err
}

func (r *blogRepository) RemoveArticleTags(ctx context.Context, articleID string, tagIDs []string) error {
	_, err := r.db.InitQuery(ctx).NewDelete().
		Table("article_tags").
		Where("blog_article_id = ?", articleID).
		Where("tag_id IN (?)", tagIDs).
		Exec(ctx)
	return err
}

func (r *blogRepository) ClearArticleTags(ctx context.Context, articleID string) error {
	_, err := r.db.InitQuery(ctx).NewDelete().
		Table("article_tags").
		Where("blog_article_id = ?", articleID).
		Exec(ctx)
	return err
}

func (r *blogRepository) DeleteArticle(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Table("blog_artikels").
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *blogRepository) HardDeleteArticle(ctx context.Context, id string) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete tags associations first
	_, err = r.db.InitQuery(ctx).NewDelete().
		Table("article_tags").
		Where("blog_article_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Delete the article
	_, err = r.db.InitQuery(ctx).NewDelete().
		Table("blog_artikels").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *blogRepository) ListPublicArticles(ctx context.Context, req requests.ListArtikel) ([]domain.BlogArtikel, int, error) {
	var res []domain.BlogArtikel

	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Relation("Author").
		Relation("Tags")

	// Apply filters
	if req.CategoryID != "" {
		q.Where("category_id = ?", req.CategoryID)
	}
	if req.TagID != "" {
		q.Join("JOIN article_tags at ON at.blog_article_id = blog_artikels.id").
			Where("at.tag_id = ?", req.TagID)
	}
	if req.Status != "" {
		q.Where("status = ?", req.Status)
	}
	if req.Search != "" {
		q.Where("title ILIKE ? OR content ILIKE ?",
			fmt.Sprintf("%%%s%%", req.Search),
			fmt.Sprintf("%%%s%%", req.Search))
	}
	if req.StartDate != nil {
		q.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		q.Where("created_at <= ?", req.EndDate)
	}

	// Apply sorting
	order := "DESC"
	orderBy := "created_at"
	if req.SortBy != "" {
		orderBy = req.SortBy
		if req.SortOrder == "asc" {
			order = "ASC"
		}
	}
	q.Order(fmt.Sprintf("%s %s", orderBy, order))

	// Apply pagination
	q.Limit(req.PageSize).
		Offset(req.CalculateOffset())

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}

func (r *blogRepository) GetPublicArticleWithRelated(ctx context.Context, slug string) (article domain.BlogArtikel, related []domain.BlogArtikel, err error) {
	// Get the main article
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&article).
		Relation("Category").
		Relation("Author").
		Relation("Tags").
		Where(`"ba"."slug" = ?`, slug).
		Where("ba.status = ?", constants.StatusPublished).
		Scan(ctx)
	if err != nil {
		return article, nil, err
	}

	// Get related articles based on category and tags
	tagConditions := []string{"category_id = ?"}
	params := []interface{}{article.CategoryID}

	// Add tag conditions if article has tags
	if len(article.Tags) > 0 {
		tagIDs := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tagIDs[i] = tag.ID
		}
		tagConditions = append(tagConditions, "EXISTS (SELECT 1 FROM article_tags at WHERE at.blog_article_id = ba.id AND at.tag_id IN (?))")
		params = append(params, tagIDs)
	}

	// Build the query
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&related).
		Relation("Category").
		Relation("Author").
		Relation("Tags").
		Where("ba.id != ?", article.ID).
		Where("ba.status = ?", constants.StatusPublished).
		Where("ba.deleted_at IS NULL").
		Where("ba.published_at <= ?", time.Now()).
		WhereOr(tagConditions[0], params[0])

	// Add tag conditions if present
	if len(tagConditions) > 1 {
		q.WhereOr(tagConditions[1], params[1])
	}

	// Order by publish date and limit to 2 articles
	q.Order("published_at DESC").
		Limit(2)

	err = q.Scan(ctx)
	if err != nil {
		return article, nil, err
	}

	return article, related, nil
}
