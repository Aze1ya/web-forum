package entity

import "time"

type Comment struct {
	ID                int
	Creator           string
	PostID            int
	Body              string
	LikesCount        int
	DislikesCount     int
	CreationDate      time.Time
	CreationDateFront string
}
