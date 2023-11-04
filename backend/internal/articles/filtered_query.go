package articles

import (
	"database/sql"
	"fmt"
	"strings"
)

// ArticleQuery holds the state of a query being built.
type ArticleQuery struct {
	DB           *sql.DB
	whereClauses []string
	args         []interface{}
	orderBy      string
}

// NewArticleQuery creates a new instance of ArticleQuery.
func NewArticleQuery(db *sql.DB) *ArticleQuery {
	return &ArticleQuery{
		DB: db,
	}
}

// WhereDateRange adds a date range condition to the query.
func (q *ArticleQuery) WhereDateRange(field string, start, end interface{}) *ArticleQuery {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("%s BETWEEN ? AND ?", field))
	q.args = append(q.args, start, end)
	return q
}

// WhereTagsIn adds a condition for matching tags.
func (q *ArticleQuery) WhereTagsIn(tags []string) *ArticleQuery {
	placeholders := make([]string, len(tags))
	for i := range tags {
		placeholders[i] = "?"
	}
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("id IN (SELECT article_id FROM article_tags WHERE tag_id IN (SELECT id FROM tags WHERE name IN (%s)))", strings.Join(placeholders, ",")))
	q.args = append(q.args, toStringInterface(tags)...)
	return q
}

// WhereEquals adds an equality condition to the query.
func (q *ArticleQuery) WhereEquals(field string, value interface{}) *ArticleQuery {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("%s = ?", field))
	q.args = append(q.args, value)
	return q
}

// WhereKeywordSearch adds a keyword search condition to the query.
func (q *ArticleQuery) WhereKeywordSearch(keyword string) *ArticleQuery {
	q.whereClauses = append(q.whereClauses, "(title LIKE ? OR content LIKE ?)")
	q.args = append(q.args, "%"+keyword+"%", "%"+keyword+"%")
	return q
}

// OrderBy adds an order by clause to the query.
func (q *ArticleQuery) OrderBy(field string, order string) *ArticleQuery {
	q.orderBy = fmt.Sprintf("ORDER BY %s %s", field, order)
	return q
}

func (q *ArticleQuery) Build() string {
	query := "SELECT id, author_id, title, content, created_at, updated_at FROM articles"
	if len(q.whereClauses) > 0 {
		query += " WHERE " + strings.Join(q.whereClauses, " AND ")
	}
	if q.orderBy != "" {
		query += " " + q.orderBy
	}

	return query
}

// Helper function to convert a slice of strings to a slice of empty interfaces.
func toStringInterface(strs []string) []interface{} {
	interfaces := make([]interface{}, len(strs))
	for i, s := range strs {
		interfaces[i] = s
	}
	return interfaces
}
