package entity

import "time"

type Post struct {
	ID                int
	Creator           string
	Title             string
	Body              string
	LikesCount        int
	DislikesCount     int
	CommentsCount     int
	CreationDate      time.Time
	CreationDateFront string
	Categories        []string
	Comments          []Comment
}
