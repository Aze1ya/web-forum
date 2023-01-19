package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (r Repository) CreateUser(user *entity.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: create user - %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO user(email, username, password, creation_time) VALUES ($1, $2, $3, $4);`
	_, err = tx.Exec(query, user.Email, user.Username, user.Password, user.CreationDate)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: user.username") {
			return fmt.Errorf("repo sqlite: create user - %w", utils.ErrUsernameNotUnique)
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed: user.email") {
			return fmt.Errorf("repo sqlite: create user - %w", utils.ErrEmailNotUnique)
		}
		return fmt.Errorf("repo sqlite: create user - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: create user - %w", err)
	}

	return nil
}

func (r Repository) GetUserByUsername(username string) (*entity.User, error) {
	query := `SELECT user_id, email, username, password, creation_time FROM user WHERE username=$1;`
	row := r.db.QueryRow(query, username)
	user := new(entity.User)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("repo sqlite: get user by username - %w", utils.ErrSqlNotFound)
		}
		return nil, fmt.Errorf("repo sqlite: get user by username - %w", err)
	}
	return user, nil
}

func (r Repository) GetSessionByToken(token string) (*entity.Session, error) {
	query := `SELECT session_id, username, session_token, session_exp_date 
				FROM session 
				JOIN user using (user_id)
				WHERE session_token=$1;`
	row := r.db.QueryRow(query, token)
	session := new(entity.Session)
	err := row.Scan(&session.ID, &session.Username, &session.Token, &session.TokenExpDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("repo sqlite: get user by token - %w", utils.ErrSqlNotFound)
		}
		return nil, fmt.Errorf("repo sqlite: get user by token - %w", err)
	}

	return session, nil
}

func (r Repository) CreateSession(session *entity.Session) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: create session - %w", err)
	}
	defer tx.Rollback()

	query := `INSERT OR REPLACE INTO session(user_id, session_token, session_exp_date) 
				VALUES ((select user_id from user where username=$1), $2, $3);`
	_, err = tx.Exec(query, session.Username, session.Token, session.TokenExpDate)
	if err != nil {
		return fmt.Errorf("repo sqlite: create session - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: create session - %w", err)
	}
	return nil
}

func (r Repository) DeleteSession(token string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("repo sqlite: delete session - %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM session WHERE session_token=$1;`
	_, err = tx.Exec(query, token)
	if err != nil {
		return fmt.Errorf("repo sqlite: delete session - %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo sqlite: delete session - %w", err)
	}

	return nil
}
