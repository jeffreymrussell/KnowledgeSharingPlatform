package articles

import (
	"KnowledgeSharingPlatform/internal"
	"KnowledgeSharingPlatform/internal/middlewares"
	"encoding/json"
	"net/http"
	"strconv"
)

// ArticleHandler struct holds the use case that will be used by the handlers
type ArticleHandler struct {
	Usecase *ArticleUsecase
}

// CreateArticle handles the creation of a new article
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var createArticleDTO CreateArticleDTO
	if err := json.NewDecoder(r.Body).Decode(&createArticleDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article, err := h.Usecase.CreateArticle(createArticleDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

// UpdateArticle handles the updating of an existing article
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/articles/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var updateArticleDTO UpdateArticleDTO
	if err := json.NewDecoder(r.Body).Decode(&updateArticleDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updateArticleDTO.ID = id // Ensure the DTO has the correct ID
	userIDStr := r.Context().Value(middlewares.UserContextKey).(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	updateArticleDTO.AuthorID = userID

	article, err := h.Usecase.UpdateArticle(updateArticleDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

// DeleteArticle handles the deletion of an article
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/articles/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	err = h.Usecase.DeleteArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetArticle handles retrieving an article
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/articles/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := h.Usecase.GetArticle(id)
	if err != nil {
		if err == internal.ErrNotFound {
			http.Error(w, "Could not find article with that ID", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

// ListArticles handles listing articles based on filters
// ListArticles handles listing articles based on filters
func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()

	// Initialize filters
	filters := make(map[string]interface{})

	// Date range filter
	startDate := queryParams.Get("startDate")
	endDate := queryParams.Get("endDate")
	if startDate != "" && endDate != "" {
		filters["startDate"] = startDate
		filters["endDate"] = endDate
	}

	// Tags filter
	tags := queryParams["tags"]
	if len(tags) > 0 {
		filters["tags"] = tags
	}

	// Author filter
	if authorID := queryParams.Get("authorId"); authorID != "" {
		filters["authorId"], _ = strconv.ParseInt(authorID, 10, 64)
	}

	// Search keyword filter
	if search := queryParams.Get("search"); search != "" {
		filters["search"] = search
	}

	// Sorting
	if sortBy := queryParams.Get("sortBy"); sortBy != "" {
		filters["sortBy"] = sortBy
	}
	if order := queryParams.Get("order"); order != "" {
		filters["order"] = order
	}

	// Call the use case with the constructed filters
	articles, err := h.Usecase.List(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}
