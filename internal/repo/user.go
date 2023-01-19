package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (r *Repository) GetUserProfile(username string) (*entity.User, error) {
	query := `SELECT user_id, email, username, creation_time FROM user WHERE username = $1`
	row := r.db.QueryRow(query, username)
	user := new(entity.User)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("repo sqlite: get user by username - %w", utils.ErrUserNotFound)
		}
		return nil, fmt.Errorf("repo sqlite: get user by username - %w", err)
	}
	return user, nil
}

func (r *Repository) GetUserPosts(username string) ([]entity.Post, error) {
	query := `SELECT post_id, username, title, post.creation_time, 
			COUNT(DISTINCT(comment.comment_id)), COUNT(DISTINCT(like.like_id)), 
			COUNT(DISTINCT(dislike.dislike_id)) 
			FROM user
			JOIN post USING (user_id) 
			LEFT JOIN comment USING (post_id) 
			LEFT JOIN like USING (post_id) 
			LEFT JOIN dislike USING (post_id) 
			WHERE user.username = $1
			GROUP BY post_id
			ORDER BY post_id desc;`

	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get user posts - %w", err)
	}
	defer rows.Close()

	var posts []entity.Post

	for rows.Next() {
		post := new(entity.Post)
		if err := rows.Scan(&post.ID, &post.Creator, &post.Title, &post.CreationDate, &post.CommentsCount, &post.LikesCount, &post.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get user posts - %w", err)
		}
		posts = append(posts, *post)
	}

	return r.addCategories(posts)
}

func (r *Repository) GetUsersCommentedPosts(username string) ([]entity.Post, error) {
	query := `SELECT post_id, username, title, post.creation_time, 
			COUNT(DISTINCT(comment.comment_id)), COUNT(DISTINCT(like.like_id)), 
			COUNT(DISTINCT(dislike.dislike_id)) 
			FROM comment c1
			JOIN user using (user_id)
			JOIN post USING (post_id) 
			LEFT JOIN comment USING (post_id) 
			LEFT JOIN like USING (post_id) 
			LEFT JOIN dislike USING (post_id) 
			WHERE c1.user_id = (SELECT user_id FROM user WHERE username = $1)
			GROUP BY post_id
			ORDER BY post_id desc;`
	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get user commented posts - %w", err)
	}
	defer rows.Close()

	var posts []entity.Post

	for rows.Next() {
		post := new(entity.Post)
		if err := rows.Scan(&post.ID, &post.Creator, &post.Title, &post.CreationDate, &post.CommentsCount, &post.LikesCount, &post.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get user commented posts - %w", err)
		}
		posts = append(posts, *post)
	}

	return r.addCategories(posts)
}

func (r *Repository) GetUsersLikedPosts(username string) ([]entity.Post, error) {
	query := `SELECT post_id, username, title, post.creation_time, 
			COUNT(DISTINCT(comment.comment_id)), COUNT(DISTINCT(like.like_id)), 
			COUNT(DISTINCT(dislike.dislike_id)) 
			FROM like l1 
			JOIN user USING (user_id) 
			JOIN post USING (post_id) 
			LEFT JOIN comment USING (post_id) 
			LEFT JOIN like USING (post_id) 
			LEFT JOIN dislike USING (post_id) 
			WHERE l1.user_id = (SELECT user_id FROM user WHERE username = $1)
			GROUP BY post_id
			ORDER BY post_id desc;`
	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get user liked posts - %w", err)
	}
	defer rows.Close()

	var posts []entity.Post

	for rows.Next() {
		post := new(entity.Post)
		if err := rows.Scan(&post.ID, &post.Creator, &post.Title, &post.CreationDate, &post.CommentsCount, &post.LikesCount, &post.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get user liked posts - %w", err)
		}
		posts = append(posts, *post)
	}

	return r.addCategories(posts)
}

func (r *Repository) GetUsersDislikedPosts(username string) ([]entity.Post, error) {
	query := `SELECT post_id, username, title, post.creation_time, 
			COUNT(DISTINCT(comment.comment_id)), COUNT(DISTINCT(like.like_id)), 
			COUNT(DISTINCT(dislike.dislike_id)) 
			FROM dislike dl1 
			JOIN user USING (user_id) 
			JOIN post USING (post_id) 
			LEFT JOIN comment USING (post_id) 
			LEFT JOIN like USING (post_id) 
			LEFT JOIN dislike USING (post_id) 
			WHERE dl1.user_id = (SELECT user_id FROM user WHERE username = $1)
			GROUP BY post_id
			ORDER BY post_id desc;`
	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: get user disliked posts - %w", err)
	}
	defer rows.Close()

	var posts []entity.Post

	for rows.Next() {
		post := new(entity.Post)
		if err := rows.Scan(&post.ID, &post.Creator, &post.Title, &post.CreationDate, &post.CommentsCount, &post.LikesCount, &post.DislikesCount); err != nil {
			return nil, fmt.Errorf("repo sqlite: get user disliked posts - %w", err)
		}
		posts = append(posts, *post)
	}

	return r.addCategories(posts)
}

func (r *Repository) addCategories(posts []entity.Post) ([]entity.Post, error) {
	query := `SELECT category_name 
					FROM post 
					JOIN user USING (user_id) 
					JOIN post_category USING (post_id) 
					JOIN category USING (category_id) 
					WHERE post_id=$1;`

	stmtCategories, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("repo sqlite: add posts categories - %w", err)
	}

	for i := 0; i < len(posts); i++ {
		rows, err := stmtCategories.Query(posts[i].ID)
		if err != nil {
			return nil, fmt.Errorf("repo sqlite: add posts categories - %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var tempCat string
			if err := rows.Scan(&tempCat); err != nil {
				return nil, fmt.Errorf("repo sqlite: add posts categories - %w", err)
			}
			posts[i].Categories = append(posts[i].Categories, tempCat)
		}
	}
	return posts, nil
}
