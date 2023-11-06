package articles_test

import (
	"KnowledgeSharingPlatform/internal"
	"KnowledgeSharingPlatform/internal/articles"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"KnowledgeSharingPlatform/internal/bootstrap"
	"KnowledgeSharingPlatform/internal/test"
	"github.com/stretchr/testify/suite"
)

type ArticleTestSuite struct {
	suite.Suite
	Server            *httptest.Server
	config            internal.DbConfig
	existingArticleID int64
	token             string
}

func (suite *ArticleTestSuite) SetupTest() {
	db := bootstrap.SetupDatabase("test_db.sqlite")
	err := test.InitializeDatabase(db)
	if err != nil {
		panic("Failed to initialize the database")
	}

	suite.config = internal.DbConfig{
		DB:         db,
		DbFilePath: "test_db.sqlite",
	}
	test.LoadTestConfig()
	router := bootstrap.SetupRouter(bootstrap.SetupHandlers(bootstrap.SetupUseCases(bootstrap.SetupAdapters(db))))
	suite.Server = httptest.NewServer(router)

	_, suite.token = test.RegisterUserAndLogin(suite.T(), suite.Server.URL)

	createArticle := articles.CreateArticleDTO{
		Title:    "Initial Test Article",
		Content:  "Initial content for the test article",
		AuthorID: 1, // Replace with a valid author ID
		Tags:     []string{"initial"},
	}
	payload, err := json.Marshal(createArticle)
	require.NoError(suite.T(), err)

	createURL := fmt.Sprintf("%s/articles", suite.Server.URL)
	req, _ := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	res, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)

	defer res.Body.Close()
	data := test.GetBodyAsString(res)
	require.Equal(suite.T(), http.StatusCreated, res.StatusCode, fmt.Sprintf("response message: %s", data))

	var createdArticle articles.Article
	err = json.Unmarshal(data, &createdArticle)
	require.NoError(suite.T(), err)
	suite.existingArticleID = createdArticle.ID
}

func (suite *ArticleTestSuite) TearDownTest() {
	suite.Server.Close()
	test.DeleteTable(suite.config)
}

func (suite *ArticleTestSuite) TestCreateArticle() {
	article := articles.CreateArticleDTO{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
		Tags:     []string{"test"},
	}
	payload, err := json.Marshal(article)
	require.NoError(suite.T(), err)

	createURL := fmt.Sprintf("%s/articles", suite.Server.URL)
	req, _ := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+suite.token)

	res, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)

	defer res.Body.Close()
	//data := test.GetBodyAsString(res)
	require.Equal(suite.T(), http.StatusCreated, res.StatusCode)

	var createdArticle articles.Article
	err = json.NewDecoder(res.Body).Decode(&createdArticle)
	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), createdArticle.ID)
	require.Equal(suite.T(), "Test Article", createdArticle.Title)
	require.Equal(suite.T(), "Test Content", createdArticle.Content)
}

func (suite *ArticleTestSuite) TestUpdateArticle() {
	// Assume an article already exists, create an UpdateArticleRequest DTO with new data
	updateArticle := articles.UpdateArticleDTO{
		ID:      suite.existingArticleID,
		Title:   "Updated Title",
		Content: "Updated Content",
		Tags:    []string{"updated"},
	}
	payload, err := json.Marshal(updateArticle)
	require.NoError(suite.T(), err)

	updateURL := fmt.Sprintf("%s/articles/%d", suite.Server.URL, updateArticle.ID)
	req, _ := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+suite.token)
	res, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)

	defer res.Body.Close()
	data := test.GetBodyAsString(res)
	require.Equal(suite.T(), http.StatusOK, res.StatusCode, fmt.Sprintf("response message: %s", data))

	var updatedArticle articles.Article
	err = json.Unmarshal(data, &updatedArticle)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "Updated Title", updatedArticle.Title)
	require.Equal(suite.T(), "Updated Content", updatedArticle.Content)
}

func (suite *ArticleTestSuite) TestDeleteArticle() {
	articleID := suite.existingArticleID

	deleteURL := fmt.Sprintf("%s/articles/%d", suite.Server.URL, articleID)
	req, _ := http.NewRequest(http.MethodDelete, deleteURL, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	res, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)

	defer res.Body.Close()
	require.Equal(suite.T(), http.StatusNoContent, res.StatusCode)

	getURL := fmt.Sprintf("%s/articles/%d", suite.Server.URL, articleID)
	req, _ = http.NewRequest(http.MethodGet, getURL, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	res, err = http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)

	defer res.Body.Close()
	require.Equal(suite.T(), http.StatusNotFound, res.StatusCode)
}

func (suite *ArticleTestSuite) TestGetArticle() {
	// Assume an article already exists, get its ID
	articleID := suite.existingArticleID // Replace with an actual ID of an existing article

	getURL := fmt.Sprintf("%s/articles/%d", suite.Server.URL, articleID)
	req, _ := http.NewRequest(http.MethodGet, getURL, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	res, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)

	defer res.Body.Close()
	require.Equal(suite.T(), http.StatusOK, res.StatusCode)

	var article articles.Article
	err = json.NewDecoder(res.Body).Decode(&article)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), articleID, article.ID)
}

func (suite *ArticleTestSuite) TestListArticles() {
	listURL := fmt.Sprintf("%s/articles", suite.Server.URL)
	req, _ := http.NewRequest(http.MethodGet, listURL, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	res, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)

	defer res.Body.Close()
	require.Equal(suite.T(), http.StatusOK, res.StatusCode)

	var articlesList []articles.Article
	err = json.NewDecoder(res.Body).Decode(&articlesList)
	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), articlesList)
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
