package up

import (
	"fmt"
	"log"

	"01.alem.school/git/Taimas/forum/config"
	"01.alem.school/git/Taimas/forum/pkg/sqlite"
)

func DbSqliteUp(cfg *config.Config) error {
	db, err := sqlite.New(cfg.DbFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("can't close db err: %v\n", err)
		} else {
			log.Printf("db closed")
		}
	}()
	if err := db.Ping(); err != nil {
		return err
	}

	tables := []string{foreignKeysOn, userTable, sessionTable, postTable, postCategoryTable, categoryTable, commentTable, likeTable, dislikeTable}
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return fmt.Errorf("storage: create tables: %w", err)
		}
	}
	return nil
}

const foreignKeysOn = `PRAGMA foreign_keys=on;`

const userTable = `CREATE TABLE IF NOT EXISTS user (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	creation_time DATE NOT NULL
);`

const sessionTable = `CREATE TABLE IF NOT EXISTS session (
	session_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL UNIQUE,
	session_token TEXT NOT NULL,
	session_exp_date DATE NOT NULL,
	FOREIGN KEY (user_id) REFERENCES user (user_id)
	);`

const postTable = `CREATE TABLE IF NOT EXISTS post (
	post_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	creation_time DATE NOT NULL,
	FOREIGN KEY (user_id) REFERENCES user (user_id)
);`

const postCategoryTable = `CREATE TABLE IF NOT EXISTS post_category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	category_id INTEGER,
	FOREIGN KEY (category_id) REFERENCES category (category_id),
	FOREIGN KEY (post_id) REFERENCES post (post_id)
);`

const categoryTable = `CREATE TABLE IF NOT EXISTS category (
	category_id INTEGER PRIMARY KEY AUTOINCREMENT,
	category_name TEXT NOT NULL UNIQUE
);`

const commentTable = `CREATE TABLE IF NOT EXISTS comment (
	comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	body TEXT NOT NULL,
	creation_time DATE NOT NULL,
	FOREIGN KEY (user_id) REFERENCES user (user_id)
);`

const likeTable = `CREATE TABLE IF NOT EXISTS like (
	like_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	post_id INTEGER DEFAULT NULL,
	comment_id INTEGER DEFAULT NULL,
	FOREIGN KEY (post_id) REFERENCES post(post_id),
	FOREIGN KEY (user_id) REFERENCES user (user_id),
	FOREIGN KEY (comment_id) REFERENCES comment (comment_id)
);`

const dislikeTable = `CREATE TABLE IF NOT EXISTS dislike (
	dislike_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	post_id INTEGER DEFAULT NULL,
	comment_id INTEGER DEFAULT NULL,
	FOREIGN KEY (post_id) REFERENCES post(post_id),
	FOREIGN KEY (user_id) REFERENCES user (user_id),
	FOREIGN KEY (comment_id) REFERENCES comment (comment_id)
);`
