package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (r *Repository) CreateComment(comment *entity.Comment) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: create comment - %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO comment(post_id, user_id, body, creation_time) 
				VALUES ($1, (SELECT user_id FROM user WHERE username = $2), $3, $4)`
	_, err = tx.Exec(query, comment.PostID, comment.Creator, comment.Body, comment.CreationDate)
	if err != nil {
		return fmt.Errorf("repo sqlite: create comment - %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: create comment - %w", err)
	}

	return nil
}

func (r Repository) GetCommentById(id int) (*entity.Comment, error) {
	query := `SELECT comment_id, post_id, username, body, creation_time, COUNT(like_id), COUNT(dislike_id)
				FROM comment 
				JOIN user USING (user_id)
				JOIN like USING (comment_id)
				JOIN dislike USING (comment_id)
				WHERE comment_id=$1`
	row := r.db.QueryRow(query, id)
	comment := new(entity.Comment)
	err := row.Scan(&comment.ID, &comment.PostID, &comment.Creator, &comment.Body, &comment.CreationDate, &comment.LikesCount, &comment.DislikesCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrSqlNotFound
		}
		return nil, fmt.Errorf("repo sqlite: get comment by id - %w", err)
	}

	return comment, nil
}
