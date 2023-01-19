package entity

type Like struct {
	ID        int
	Creator   string
	PostID    int
	CommentID int
}

type DisLike struct {
	ID        int
	Creator   string
	PostID    int
	CommentID int
}
