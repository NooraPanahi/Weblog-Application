package model

import "time"

type Comment struct {
	ID        int64
	WeblogID  int64
	UserID    int64
	Content   string
	CreatedAt time.Time
}
