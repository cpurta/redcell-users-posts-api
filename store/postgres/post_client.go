package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"redcellpartners.com/users-posts-api/model"
	"redcellpartners.com/users-posts-api/store"
)

var _ store.PostStore = &PostgresPostClient{}

type PostgresPostClient struct {
	listPostsStmt  *sql.Stmt
	createPostStmt *sql.Stmt
	getPostStmt    *sql.Stmt
	updatePostStmt *sql.Stmt
	deletePostStmt *sql.Stmt

	logger *zap.Logger
}

func NewPostgresPostClient(db *sql.DB, logger *zap.Logger) (*PostgresPostClient, error) {
	client := &PostgresPostClient{
		logger: logger,
	}

	var err error

	client.listPostsStmt, err = db.Prepare("SELECT id, user_id, title, content, created_at, updated_at FROM posts LIMIT 100;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare list posts statement: %s", err.Error())
	}

	client.createPostStmt, err = db.Prepare("INSERT INTO posts (user_id, title, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare create post statement: %s", err.Error())
	}

	client.getPostStmt, err = db.Prepare("SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE id = $1;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare get post statement: %s", err.Error())
	}

	client.updatePostStmt, err = db.Prepare("UPDATE posts SET title = $2, content = $3, updated_at = $4 WHERE id = $1 RETURNING id, user_id, title, content, created_at, updated_at;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare update post statement: %s", err.Error())
	}

	client.deletePostStmt, err = db.Prepare("DELETE FROM posts WHERE id = $1;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare update post statement: %s", err.Error())
	}

	return client, nil
}

func (client *PostgresPostClient) ListPosts() ([]*model.Post, error) {
	rows, err := client.listPostsStmt.Query()
	if err != nil {
		client.logger.Error("unable to list all posts", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	posts := make([]*model.Post, 0)

	for rows.Next() {
		var (
			post        = &model.Post{}
			timeUpdated sql.NullString
		)

		if err := rows.Scan(
			&post.ID,
			&post.CreatedByUser,
			&post.Title,
			&post.Content,
			&post.CreatedTime,
			&timeUpdated,
		); err != nil {
			client.logger.Error("unable to scan post, skipping for now", zap.Error(err))
			continue
		}

		if timeUpdated.Valid {
			updatedAt, err := time.Parse(time.RFC3339, timeUpdated.String)
			if err != nil {
				client.logger.Warn("unable to parse updated_at value", zap.Int("post_id", post.ID), zap.Error(err))
			} else {
				post.UpdatedTime = updatedAt
			}
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (client *PostgresPostClient) CreatePost(post *model.Post) (*model.Post, error) {
	row := client.createPostStmt.QueryRow(post.CreatedByUser, post.Title, post.Content, time.Now())

	var postID int64

	err := row.Scan(&postID)
	if err != nil {
		return nil, fmt.Errorf("unable to scan created post id: %s", err.Error())
	}

	createdPost, err := client.GetPost(int(postID))
	if err != nil {
		return nil, fmt.Errorf("unable to get created post: %s", err.Error())
	}

	return createdPost, nil
}

func (client *PostgresPostClient) GetPost(id int) (*model.Post, error) {
	var (
		post        = &model.Post{}
		timeUpdated sql.NullString
		err         error
	)

	row := client.getPostStmt.QueryRow(id)

	if err = row.Scan(
		&post.ID,
		&post.CreatedByUser,
		&post.Title,
		&post.Content,
		&post.CreatedTime,
		&timeUpdated,
	); err != nil && err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("unable to scan post [%d]: %s", id, err.Error())
	}

	if timeUpdated.Valid {
		updatedAt, err := time.Parse(time.RFC3339, timeUpdated.String)
		if err != nil {
			client.logger.Warn("unable to parse updated_at value", zap.Int("post_id", post.ID), zap.Error(err))
		} else {
			post.UpdatedTime = updatedAt
		}
	}

	return post, nil
}

func (client *PostgresPostClient) UpdatePost(postInput *model.Post) (*model.Post, error) {
	row := client.updatePostStmt.QueryRow(postInput.ID, postInput.Title, postInput.Content, time.Now())

	var (
		post        = &model.Post{}
		timeUpdated sql.NullString
		err         error
	)

	if err = row.Scan(
		&post.ID,
		&post.CreatedByUser,
		&post.Title,
		&post.Content,
		&post.CreatedTime,
		&timeUpdated,
	); err != nil {
		return nil, fmt.Errorf("unable to scan post [%d]: %s", postInput.ID, err.Error())
	}

	if timeUpdated.Valid {
		updatedAt, err := time.Parse(time.RFC3339, timeUpdated.String)
		if err != nil {
			client.logger.Warn("unable to parse time_update value", zap.Int("post_id", post.ID), zap.Error(err))
		} else {
			post.UpdatedTime = updatedAt
		}
	}
	return post, nil
}

func (client *PostgresPostClient) DeletePost(id int) error {
	result, err := client.deletePostStmt.Exec(id)
	if err != nil {
		return fmt.Errorf("unable to delete post [%d]: %s", id, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for post [%d]: %s", id, err.Error())
	}

	if rowsAffected != int64(1) {
		return fmt.Errorf("deleted 0 or more than one post requested")
	}

	return nil
}
