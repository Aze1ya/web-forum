package repo

import (
	"database/sql"
	"fmt"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (r Repository) CreateDisLikePost(dislike *entity.DisLike) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: create dislike post - %w", err)
	}
	defer tx.Rollback()

	tempLikeId := 0
	tempDislikeId := 0
	query := `select post_id, 
       			ifnull((select like_id from like where post_id=$1 and user_id=(select user_id from user where username=$2)), 0),
       			ifnull((select dislike_id from dislike where post_id=$1 and user_id=(select user_id from user where username=$2)), 0)
				from post
				where post.post_id = $1;`
	row := tx.QueryRow(query, dislike.PostID, dislike.Creator)
	err = row.Scan(&dislike.PostID, &tempLikeId, &tempDislikeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("repo sqlite: create dislike post - %w", utils.ErrInvalidPostId)
		} else {
			return fmt.Errorf("repo sqlite: create dislike post - %w", err)
		}
	}

	if tempDislikeId == 0 {
		query = `delete from like where like_id = $1;
				insert into dislike (user_id, post_id, comment_id)
				values ((select user_id from user where username = $2), $3, $4);`
		_, err = tx.Exec(query, tempLikeId, dislike.Creator, dislike.PostID, dislike.CommentID)
		if err != nil {
			return fmt.Errorf("repo sqlite: create dislike post - %w", err)
		}
	} else {
		query = `delete from dislike where dislike_id = $1;`
		_, err = tx.Exec(query, tempDislikeId)
		if err != nil {
			return fmt.Errorf("repo sqlite: create dislike post - %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: create dislike post - %w", err)
	}
	return nil
}

func (r Repository) CreateDisLikeComment(dislike *entity.DisLike) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: create dislike comment- %w", err)
	}
	defer tx.Rollback()

	tempLikeId := 0
	tempDislikeId := 0
	query := `select comment_id, 
       			ifnull((select like_id from like where comment_id=$1 and user_id=(select user_id from user where username=$2)), 0),
       			ifnull((select dislike_id from dislike where comment_id=$1 and user_id=(select user_id from user where username=$2)), 0)
				from comment
				where comment.comment_id = $1 and comment.post_id=$3;`
	row := tx.QueryRow(query, dislike.CommentID, dislike.Creator, dislike.PostID)
	err = row.Scan(&dislike.CommentID, &tempLikeId, &tempDislikeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("repo sqlite: create dislike comment - %w", utils.ErrInvalidCommentId)
		} else {
			return fmt.Errorf("repo sqlite: create dislike comment - %w", err)
		}
	}

	if tempDislikeId == 0 {
		query = `delete from like where like_id = $1;
				insert into dislike (user_id, post_id, comment_id)
				values ((select user_id from user where username = $2), $3, $4);`
		_, err = tx.Exec(query, tempLikeId, dislike.Creator, 0, dislike.CommentID)
		if err != nil {
			return fmt.Errorf("repo sqlite: create dislike comment - %w", err)
		}
	} else {
		query = `delete from dislike where dislike_id = $1;`
		_, err = tx.Exec(query, tempDislikeId)
		if err != nil {
			return fmt.Errorf("repo sqlite: create dislike comment - %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: create dislike comment - %w", err)
	}
	return nil
}
