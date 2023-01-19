package repo

import (
	"fmt"
	"strings"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (r Repository) CreatePost(post *entity.Post) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("repo sqlite: create post - %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO post(user_id, title, body, creation_time) 
				VALUES ((SELECT user_id FROM user WHERE username = $1), $2, $3, $4);`
	res, err := tx.Exec(query, post.Creator, post.Title, post.Body, post.CreationDate)
	if err != nil {
		return 0, fmt.Errorf("repo sqlite: create post - %w", err)
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("repo sqlite: create post - %w", err)
	}

	query = `INSERT INTO category (category_name) VALUES ($2);`
	for _, category := range post.Categories {
		_, err := tx.Exec(query, category)
		if err != nil && !strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, fmt.Errorf("repo sqlite: create post - %w", err)
		}
	}

	query = `INSERT INTO post_category (post_id, category_id) 
				VALUES ($1, (SELECT category_id FROM category WHERE category_name = $2));`
	for _, category := range post.Categories {
		_, err := tx.Exec(query, postID, category)
		if err != nil {
			return 0, fmt.Errorf("repo sqlite: create post - %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("repo sqlite: create post - %w", err)
	}

	return postID, nil
}

func (r Repository) GetPostById(post_id int) (*entity.Post, error) {
	query := `SELECT post_id, username, title, post.body, post.creation_time, category_name, 
				COUNT(distinct(comment.comment_id)), COUNT(distinct(like.like_id)), 
				COUNT(distinct(dislike.dislike_id)) 
				FROM post 
				JOIN user USING (user_id) 
				JOIN post_category USING (post_id) 
				JOIN category USING (category_id) 
				LEFT JOIN comment USING (post_id) 
				LEFT JOIN like USING (post_id) 
				LEFT JOIN dislike USING (post_id) 
				WHERE post.post_id=$1
				GROUP BY post_id, category_id;`

	post := new(entity.Post)
	rows, err := r.db.Query(query, post_id)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get post by id - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var tempCat string
		if err := rows.Scan(&post.ID, &post.Creator, &post.Title, &post.Body, &post.CreationDate, &tempCat, &post.CommentsCount, &post.LikesCount, &post.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get post by id - %w", err)
		}
		post.Categories = append(post.Categories, tempCat)
	}
	if post.ID == 0 {
		return nil, utils.ErrSqlNotFound
	}

	query = `SELECT comment_id, comment.post_id, username, body, comment.creation_time,
			 COUNT(like_id), COUNT(dislike_id) 
			 FROM comment 
			 JOIN user USING (user_id)
			 LEFT JOIN like USING (comment_id)
			 LEFT JOIN dislike USING (comment_id)
			 WHERE comment.post_id=$1
			 GROUP BY comment_id;`

	rows, err = r.db.Query(query, post_id)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get post by id - %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		comment := new(entity.Comment)
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Creator, &comment.Body, &comment.CreationDate, &comment.LikesCount, &comment.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get post by id - %w", err)
		}
		post.Comments = append(post.Comments, *comment)
	}

	return post, nil
}

func (r Repository) GetAllPosts() ([]entity.Post, error) {
	query := `SELECT post_id, username, title, post.creation_time, 
			COUNT(DISTINCT(comment.comment_id)), COUNT(DISTINCT(like.like_id)), 
			COUNT(DISTINCT(dislike.dislike_id)) 
			FROM post
			JOIN user USING (user_id) 
			LEFT JOIN comment USING (post_id) 
			LEFT JOIN like USING (post_id) 
			LEFT JOIN dislike USING (post_id) 
			GROUP BY post_id
			ORDER BY post_id desc;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get all posts - %w", err)
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		post := new(entity.Post)
		if err := rows.Scan(&post.ID, &post.Creator, &post.Title, &post.CreationDate, &post.CommentsCount, &post.LikesCount, &post.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get all posts - %w", err)
		}
		posts = append(posts, *post)
	}

	posts, err = r.addCategories(posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r Repository) GetPostsByCategoryName(category string) ([]entity.Post, error) {
	query := `SELECT post_id, username, title, post.creation_time, 
			COUNT(DISTINCT(comment.comment_id)), COUNT(DISTINCT(like.like_id)), 
			COUNT(DISTINCT(dislike.dislike_id)) 
			FROM post_category 
			JOIN category USING (category_id) 
			JOIN post using (post_id) 
			JOIN user USING (user_id) 
			LEFT JOIN comment USING (post_id) 
			LEFT JOIN like USING (post_id) 
			LEFT JOIN dislike USING (post_id) 
			WHERE category_name = $1 
			GROUP BY post_id 
			ORDER BY post_id desc;`

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get all posts - %w", err)
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		post := new(entity.Post)
		if err := rows.Scan(&post.ID, &post.Creator, &post.Title, &post.CreationDate, &post.CommentsCount, &post.LikesCount, &post.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get all posts - %w", err)
		}
		posts = append(posts, *post)
	}

	posts, err = r.addCategories(posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
