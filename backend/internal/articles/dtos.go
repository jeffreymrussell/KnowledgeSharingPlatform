package articles

// CreateArticleDTO carries the data required to create a new article
type CreateArticleDTO struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	AuthorID int64    `json:"author_id"`
	Tags     []string `json:"tags"`
}

// Validate performs basic validation on CreateArticleDTO fields
func (dto *CreateArticleDTO) Validate() error {
	// Add validation logic here (e.g., check for empty fields, etc.)
	return nil
}

// UpdateArticleDTO carries the data required to update an existing article
type UpdateArticleDTO struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	AuthorID int64    `json:"author_id"` // This can be used to verify if the user making the request is the author of the article
	Tags     []string `json:"tags"`
}

// Validate performs basic validation on UpdateArticleDTO fields
func (dto *UpdateArticleDTO) Validate() error {
	// Add validation logic here
	return nil
}

// ArticleResponse is used to send article data back to the client
type ArticleResponse struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	AuthorID  int64    `json:"author_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// ConvertToResponse converts an article entity to an ArticleResponse DTO
func ConvertToResponse(article Article) ArticleResponse {
	return ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		AuthorID:  article.AuthorID,
		Tags:      article.Tags,
		CreatedAt: article.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: article.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ListArticlesFilters defines the filters that can be applied when listing articles
type ListArticlesFilters struct {
	// Define any filters you want to apply when listing articles (e.g., by tag, author, date range, etc.)
}
