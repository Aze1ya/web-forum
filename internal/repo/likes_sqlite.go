package repo

import (
	"database/sql"
	"fmt"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (r Repository) CreateLikePost(like *entity.Like) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: create like post - %w", err)
	}
	defer tx.Rollback()

	tempLikeId := 0
	tempDislikeId := 0
	query := `select post_id, 
       			ifnull((select like_id from like where post_id=$1 and user_id=(select user_id from user where username=$2)), 0),
       			ifnull((select dislike_id from dislike where post_id=$1 and user_id=(select user_id from user where username=$2)), 0)
				from post
				where post.post_id = $1;`
	row := tx.QueryRow(query, like.PostID, like.Creator)
	err = row.Scan(&like.PostID, &tempLikeId, &tempDislikeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("repo sqlite: create like post - %w", utils.ErrInvalidPostId)
		} else {
			return fmt.Errorf("repo sqlite: create like post - %w", err)
		}
	}

	if tempLikeId == 0 {
		query = `delete from dislike where dislike_id = $1;
				insert into like (user_id, post_id, comment_id)
				values ((select user_id from user where username = $2), $3, $4);`
		_, err = tx.Exec(query, tempDislikeId, like.Creator, like.PostID, like.CommentID)
		if err != nil {
			return fmt.Errorf("repo sqlite: create like post - %w", err)
		}
	} else {
		query = `delete from like where like_id = $1;`
		_, err = tx.Exec(query, tempLikeId)
		if err != nil {
			return fmt.Errorf("repo sqlite: create like post - %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: create like post - %w", err)
	}
	return nil
}

func (r Repository) CreateLikeComment(like *entity.Like) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: create like comment- %w", err)
	}
	defer tx.Rollback()
	tempLikeId := 0
	tempDislikeId := 0
	query := `select comment_id, 
       			ifnull((select like_id from like where comment_id=$1 and user_id=(select user_id from user where username=$2)), 0),
       			ifnull((select dislike_id from dislike where comment_id=$1 and user_id=(select user_id from user where username=$2)), 0)
				from comment
				where comment.comment_id = $1 and comment.post_id=$3;`
	row := tx.QueryRow(query, like.CommentID, like.Creator, like.PostID)
	err = row.Scan(&like.CommentID, &tempLikeId, &tempDislikeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("repo sqlite: create like comment - %w", utils.ErrInvalidCommentId)
		} else {
			return fmt.Errorf("repo sqlite: create like comment - %w", err)
		}
	}

	if tempLikeId == 0 {
		query = `delete from dislike where dislike_id = $1;
				insert into like (user_id, post_id, comment_id)
				values ((select user_id from user where username = $2), $3, $4);`
		_, err = tx.Exec(query, tempDislikeId, like.Creator, 0, like.CommentID)
		if err != nil {
			return fmt.Errorf("repo sqlite: create like comment - %w", err)
		}
	} else {
		query = `delete from like where like_id = $1;`
		_, err = tx.Exec(query, tempLikeId)
		if err != nil {
			return fmt.Errorf("repo sqlite: create like comment - %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: create like comment - %w", err)
	}
	return nil
}
