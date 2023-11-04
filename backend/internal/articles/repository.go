package articles

import (
	"KnowledgeSharingPlatform/internal"
	"database/sql"
	"errors"
)

type SQLiteArticleAdapter struct {
	DB *sql.DB
}

func (adapter *SQLiteArticleAdapter) NewQuery() *ArticleQuery {
	return NewArticleQuery(adapter.DB)
}

func (adapter *SQLiteArticleAdapter) Save(article Article) (Article, error) {
	result, err := adapter.DB.Exec(`INSERT INTO articles (author_id, title, content, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		article.AuthorID, article.Title, article.Content)
	if err != nil {
		return Article{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Article{}, err
	}

	article.ID = id
	return article, nil
}

func (adapter *SQLiteArticleAdapter) Update(article Article) (Article, error) {
	_, err := adapter.DB.Exec(`UPDATE articles SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		article.Title, article.Content, article.ID)
	if err != nil {
		return Article{}, err
	}

	return article, nil
}

func (adapter *SQLiteArticleAdapter) Delete(id int64) error {
	_, err := adapter.DB.Exec(`DELETE FROM articles WHERE id = ?`, id)
	return err
}

func (adapter *SQLiteArticleAdapter) FindByID(id int64) (Article, error) {
	var article Article
	row := adapter.DB.QueryRow(`SELECT id, author_id, title, content, created_at, updated_at FROM articles WHERE id = ?`, id)
	err := row.Scan(&article.ID, &article.AuthorID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Article{}, internal.ErrNotFound
		}
		return Article{}, err
	}
	return article, nil
}

func (adapter *SQLiteArticleAdapter) FindAll(articleQuery *ArticleQuery) ([]Article, error) {
	// This is a simplified example. You'll need to build the SQL query based on the filters provided.
	query := articleQuery.Build()
	rows, err := adapter.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.ID, &article.AuthorID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
