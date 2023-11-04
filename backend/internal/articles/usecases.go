package articles

import (
	"KnowledgeSharingPlatform/internal/users"
	"errors"
	"strconv"
	"time"
)

type ArticleRepository interface {
	Save(article Article) (Article, error)
	Update(article Article) (Article, error)
	Delete(id int64) error
	FindByID(id int64) (Article, error)
	FindAll(query *ArticleQuery) ([]Article, error)
	NewQuery() *ArticleQuery
}

type ArticleUsecase struct {
	UserUsecase *users.UserUsecase
	ArticleRepo ArticleRepository
}

// CreateArticle creates a new article
func (usecase *ArticleUsecase) CreateArticle(createArticleDTO CreateArticleDTO) (Article, error) {
	// Validate the article data
	if err := createArticleDTO.Validate(); err != nil {
		return Article{}, err
	}

	// Create the article entity
	article := Article{
		AuthorID:  createArticleDTO.AuthorID,
		Title:     createArticleDTO.Title,
		Content:   createArticleDTO.Content,
		Tags:      createArticleDTO.Tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the article
	return usecase.ArticleRepo.Save(article)
}

// UpdateArticle updates an existing article
func (usecase *ArticleUsecase) UpdateArticle(updateArticleDTO UpdateArticleDTO) (Article, error) {
	// Retrieve the existing article
	article, err := usecase.ArticleRepo.FindByID(updateArticleDTO.ID)
	if err != nil {
		return Article{}, err
	}

	// Check if the user is the author of the article
	if article.AuthorID != updateArticleDTO.AuthorID {
		return Article{}, errors.New("unauthorized: user is not the author of the article" + strconv.FormatInt(article.AuthorID, 10) + "!=" + strconv.FormatInt(updateArticleDTO.AuthorID, 10))
	}

	// Update the article's data
	article.Title = updateArticleDTO.Title
	article.Content = updateArticleDTO.Content
	article.Tags = updateArticleDTO.Tags
	article.UpdatedAt = time.Now()

	// Save the updated article
	return usecase.ArticleRepo.Update(article)
}

// DeleteArticle deletes an article
func (usecase *ArticleUsecase) DeleteArticle(id int64) error {
	return usecase.ArticleRepo.Delete(id)
}

// GetArticle retrieves an article
func (usecase *ArticleUsecase) GetArticle(id int64) (Article, error) {
	return usecase.ArticleRepo.FindByID(id)
}

func (usecase *ArticleUsecase) List(filters map[string]interface{}) ([]Article, error) {
	query := usecase.ArticleRepo.NewQuery()

	if startDate, ok := filters["startDate"]; ok && filters["endDate"] != nil {
		query = query.WhereDateRange("created_at", startDate, filters["endDate"])
	}

	if tags, ok := filters["tags"].([]string); ok && len(tags) > 0 {
		query = query.WhereTagsIn(tags)
	}

	if authorId, ok := filters["authorId"].(int64); ok {
		query = query.WhereEquals("author_id", authorId)
	}

	if search, ok := filters["search"].(string); ok {
		query = query.WhereKeywordSearch(search)
	}

	if sortBy, ok := filters["sortBy"].(string); ok {
		order := "ASC"
		if sortOrder, ok := filters["order"].(string); ok {
			order = sortOrder
		}
		query = query.OrderBy(sortBy, order)
	}

	return usecase.ArticleRepo.FindAll(query)
}
