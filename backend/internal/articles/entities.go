package articles

import "time"

// Article entity

type Article struct {
	ID        int64
	AuthorID  int64
	Title     string
	Content   string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
