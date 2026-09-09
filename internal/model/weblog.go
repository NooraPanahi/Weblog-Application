package model

import "time"

type Weblog struct {
	ID         int64
	Title      string
	Content    string
	Image      *string
	AuthorID   int64
	AuthorName string
	Privacy    string
	CreatedAt  time.Time
}
