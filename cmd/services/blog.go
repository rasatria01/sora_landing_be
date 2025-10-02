package services

import (
	"context"
	"database/sql"
	"errors"

	"net/http"
	"sora_landing_be/cmd/constants"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/database"
	internal_err "sora_landing_be/pkg/errors"
	"sora_landing_be/pkg/logger"
	"sora_landing_be/pkg/utils"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type BlogService interface {
	// Create and Update operations
	CreateArticle(ctx context.Context, userID string, payload requests.BlogArtikel) error
	CreateArticleFromURL(ctx context.Context, userID string, payload requests.FromURL) error
	// CreateByLink(ctx context.Context, payl)
	UpdateArticle(ctx context.Context, id string, payload requests.UpdateArtikel) error
	UpdateArticleStatus(ctx context.Context, id string, payload requests.UpdateArticleStatus) error

	// Read operations
	GetArticle(ctx context.Context, id string) (response.BlogArticle, error)
	GetArticleBySlug(ctx context.Context, slug string) (response.BlogArticle, error)
	ListArticles(ctx context.Context, params requests.ListArtikel) (dto.PaginationResponse[response.BlogArticleList], error)
	GetArticleStats(ctx context.Context) (dto.BlogStats, error)

	// Public endpoints
	ListPublicArticles(ctx context.Context, params requests.ListArtikel) (dto.PaginationResponse[response.PublicArticleList], error)
	GetPublicArticleBySlug(ctx context.Context, slug string) (response.PublicArticleDetail, error)

	// Tag operations
	UpdateArticleTags(ctx context.Context, articleID string, tagIDs []string) error

	// Delete operations
	DeleteArticle(ctx context.Context, id string) error
	HardDeleteArticle(ctx context.Context, id string) error
}

type blogService struct {
	blogRepo repository.BlogRepository
	tagRepo  repository.TagRepository
	catRepo  repository.CategoryRepository
}

func NewBlogService(blogRepo repository.BlogRepository, tagRepo repository.TagRepository, catRepo repository.CategoryRepository) BlogService {
	return &blogService{
		blogRepo: blogRepo,
		tagRepo:  tagRepo,
		catRepo:  catRepo,
	}
}

func (s *blogService) CreateArticle(ctx context.Context, userID string, payload requests.BlogArtikel) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// Generate unique slug
		uniqueSlug, err := utils.GenerateUniqueSlug(ctx, s.blogRepo, payload.Title)
		if err != nil {
			return err
		}

		// Convert to domain model
		article := payload.ToDomain(userID, uniqueSlug)

		// Create the article
		err = s.blogRepo.CreateArticle(ctx, article)
		if err != nil {
			return err
		}

		// Add tags if provided
		if len(payload.TagIDs) > 0 {
			err = s.blogRepo.AddArticleTags(ctx, article.ID, payload.TagIDs)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *blogService) CreateArticleFromURL(ctx context.Context, userID string, payload requests.FromURL) error {
	return database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		// --- 1. Ensure "external" category exists ---
		cat, err := s.catRepo.GetCategoryByName(ctx, "external")
		if err != nil {
			logger.Log.Error("Error fetching category", zap.Error(err))
			return err
		}

		if cat == nil {
			slug, err := utils.GenerateUniqueSlug(ctx, s.catRepo, "external")
			if err != nil {
				logger.Log.Error("Error generating category slug", zap.Error(err))
				return err
			}

			newCat := &domain.Category{
				Name:        "external",
				Slug:        slug,
				CreatedByID: userID,
			}

			catID, err := s.catRepo.CreateCategoryReturnID(ctx, newCat)
			if err != nil {
				logger.Log.Error("Error creating category", zap.Error(err))
				return err
			}

			cat = &domain.Category{
				BaseEntity: domain.BaseEntity{ID: catID},
				Name:       "external",
				Slug:       slug,
			}
		}

		// --- 2. Ensure "external" tag exists ---
		tag, err := s.tagRepo.GetTagByName(ctx, "external")
		if err != nil {
			logger.Log.Error("Error fetching tag", zap.Error(err))
			return err
		}

		if tag == nil {
			slug, err := utils.GenerateUniqueSlug(ctx, s.tagRepo, "external")
			if err != nil {
				logger.Log.Error("Error generating tag slug", zap.Error(err))
				return err
			}

			newTag := &domain.Tag{
				Name:        "external",
				Slug:        slug,
				CreatedByID: userID,
			}

			tagID, err := s.tagRepo.CreateTagReturnID(ctx, newTag)
			if err != nil {
				logger.Log.Error("Error creating tag", zap.Error(err))
				return err
			}

			tag = &domain.Tag{
				BaseEntity: domain.BaseEntity{ID: tagID},
				Name:       "external",
				Slug:       slug,
			}
		}

		// --- 3. Fetch article from URL ---
		extractedArticle, err := readability.FromURL(payload.URL, 30*time.Second)
		if err != nil {
			logger.Log.Error("Readability error", zap.Error(err), zap.String("url", payload.URL))
			return err
		}

		// --- 4. Generate unique slug for the article ---
		uniqueSlug, err := utils.GenerateUniqueSlug(ctx, s.blogRepo, extractedArticle.Title)
		if err != nil {
			logger.Log.Error("Error generating article slug", zap.Error(err))
			return err
		}

		// --- 5. Convert to domain model ---
		articleDomain := &domain.BlogArtikel{
			Title:       extractedArticle.Title,
			Slug:        uniqueSlug,
			Content:     extractedArticle.Content,
			Excerpt:     extractedArticle.Excerpt,
			CategoryID:  cat.ID,
			ImageURL:    extractedArticle.Image,
			AuthorID:    userID,
			Status:      constants.StatusPublished,
			PublishedAt: time.Now(),
			Tags:        []*domain.Tag{}, // start empty
		}

		// --- 6. Save article ---
		if err := s.blogRepo.CreateArticlefromURL(ctx, articleDomain); err != nil {
			logger.Log.Error("Error creating article", zap.Error(err))
			return err
		}

		// --- 7. Add "external" tag to article ---
		if tag != nil && tag.ID != "" {
			if err := s.blogRepo.AddArticleTags(ctx, articleDomain.ID, []string{tag.ID}); err != nil {
				logger.Log.Error("Error adding tag to article", zap.Error(err), zap.String("article_id", articleDomain.ID))
				return err
			}
		}

		logger.Log.Info("Article successfully created",
			zap.String("article_id", articleDomain.ID),
			zap.String("category", cat.Name),
			zap.String("tag", tag.Name),
			zap.String("url", payload.URL),
		)

		return nil
	})
}

func (s *blogService) UpdateArticle(ctx context.Context, id string, payload requests.UpdateArtikel) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// Get existing article
		existing, err := s.blogRepo.GetArticle(ctx, id)
		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				return internal_err.NewDefaultError(http.StatusNotFound, internal_err.DataNotFound)
			}
			return err
		}

		// Generate new slug if title changed
		var uniqueSlug string
		if payload.Title != "" && existing.Title != payload.Title {
			uniqueSlug, err = utils.GenerateUniqueSlug(ctx, s.blogRepo, payload.Title)
			if err != nil {
				return err
			}
		} else {
			uniqueSlug = existing.Slug
		}

		// Update article
		article := payload.ToDomain(existing.AuthorID, uniqueSlug)
		article.ID = id
		article.Views = existing.Views
		article.PublishedAt = existing.PublishedAt

		err = s.blogRepo.UpdateArticle(ctx, article)
		if err != nil {
			return err
		}

		// Update tags if provided
		if payload.TagIDs != nil {
			err = s.blogRepo.ClearArticleTags(ctx, id)
			if err != nil {
				return err
			}

			if len(payload.TagIDs) > 0 {
				err = s.blogRepo.AddArticleTags(ctx, id, payload.TagIDs)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (s *blogService) UpdateArticleStatus(ctx context.Context, id string, payload requests.UpdateArticleStatus) error {
	// Get existing article
	_, err := s.blogRepo.GetArticle(ctx, id)
	if err != nil {
		return err
	}

	var publishAt *time.Time
	switch payload.Status {
	case constants.StatusPublished:
		now := time.Now()
		publishAt = &now
	case constants.StatusScheduled:
		publishAt = payload.PublishAt
	}

	return s.blogRepo.UpdateArticleStatus(ctx, id, payload.Status, publishAt)
}

func (s *blogService) GetArticle(ctx context.Context, id string) (response.BlogArticle, error) {
	var res response.BlogArticle

	article, err := s.blogRepo.GetArticle(ctx, id)
	if err != nil {
		return res, err
	}

	// Increment views asynchronously

	res.FromDomain(&article)
	return res, nil
}

func (s *blogService) GetArticleBySlug(ctx context.Context, slug string) (response.BlogArticle, error) {
	var res response.BlogArticle

	article, err := s.blogRepo.GetArticleBySlug(ctx, slug)
	if err != nil {
		return res, err
	}

	// Increment views asynchronously

	res.FromDomain(&article)
	return res, nil
}

func (s *blogService) ListArticles(ctx context.Context, params requests.ListArtikel) (dto.PaginationResponse[response.BlogArticleList], error) {
	var paginateRes dto.PaginationResponse[response.BlogArticleList]

	articles, count, err := s.blogRepo.ListArticles(ctx, params)
	if err != nil {
		return paginateRes, err
	}

	// Convert domain models to response DTOs
	list := make([]response.BlogArticleList, len(articles))
	for i, article := range articles {
		var item response.BlogArticleList
		item.FromDomain(&article)
		list[i] = item
	}

	paginateRes = dto.NewPaginationResponse(params.PaginationRequest, count, list)
	return paginateRes, nil
}

func (s *blogService) GetArticleStats(ctx context.Context) (dto.BlogStats, error) {
	stats, err := s.blogRepo.GetArticleStats(ctx)
	if err != nil {
		return dto.BlogStats{}, err
	}
	return stats, nil
}

func (s *blogService) UpdateArticleTags(ctx context.Context, articleID string, tagIDs []string) error {
	// Get existing article
	_, err := s.blogRepo.GetArticle(ctx, articleID)
	if err != nil {
		return err
	}

	// Verify all tags exist
	for _, tagID := range tagIDs {
		_, err := s.tagRepo.GetTag(ctx, tagID)
		if err != nil {
			return internal_err.NewDefaultError(http.StatusBadRequest, "Invalid tag ID")
		}
	}

	err = database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// Clear existing tags
		err := s.blogRepo.ClearArticleTags(ctx, articleID)
		if err != nil {
			return err
		}

		// Add new tags
		if len(tagIDs) > 0 {
			err = s.blogRepo.AddArticleTags(ctx, articleID, tagIDs)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *blogService) DeleteArticle(ctx context.Context, id string) error {
	// Verify article exists
	_, err := s.blogRepo.GetArticle(ctx, id)
	if err != nil {
		return err
	}

	return s.blogRepo.DeleteArticle(ctx, id)
}

func (s *blogService) HardDeleteArticle(ctx context.Context, id string) error {
	// Verify article exists
	_, err := s.blogRepo.GetArticle(ctx, id)
	if err != nil {
		return err
	}

	return s.blogRepo.HardDeleteArticle(ctx, id)
}

func (s *blogService) ListPublicArticles(ctx context.Context, params requests.ListArtikel) (dto.PaginationResponse[response.PublicArticleList], error) {
	var paginateRes dto.PaginationResponse[response.PublicArticleList]

	articles, count, err := s.blogRepo.ListPublicArticles(ctx, params)
	if err != nil {
		return paginateRes, err
	}

	// Convert domain models to response DTOs
	list := make([]response.PublicArticleList, len(articles))
	for i, article := range articles {
		var item response.PublicArticleList
		item.FromDomain(&article)
		list[i] = item
	}

	paginateRes = dto.NewPaginationResponse(params.PaginationRequest, count, list)
	return paginateRes, nil
}

func (s *blogService) GetPublicArticleBySlug(ctx context.Context, slug string) (response.PublicArticleDetail, error) {
	var res response.PublicArticleDetail

	article, related, err := s.blogRepo.GetPublicArticleWithRelated(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, internal_err.NewDefaultError(http.StatusNotFound, internal_err.DataNotFound)
		}
		return res, err
	}

	// Increment views asynchronously
	go func() {
		bgCtx := context.Background()
		_ = s.blogRepo.IncrementViews(bgCtx, article.ID)
	}()

	// Convert to response DTO with related articles
	res.FromDomain(&article, related)
	return res, nil
}
