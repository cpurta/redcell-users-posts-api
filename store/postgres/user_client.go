package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"redcellpartners.com/users-posts-api/model"
	"redcellpartners.com/users-posts-api/store"
)

var _ store.UserStore = &PostgresUserClient{}

type PostgresUserClient struct {
	listUsersStmt  *sql.Stmt
	createUserStmt *sql.Stmt
	getUserStmt    *sql.Stmt
	updateUserStmt *sql.Stmt
	deleteUserStmt *sql.Stmt

	logger *zap.Logger
}

func NewPostgresUserClient(db *sql.DB, logger *zap.Logger) (*PostgresUserClient, error) {
	client := &PostgresUserClient{
		logger: logger,
	}

	var err error

	client.listUsersStmt, err = db.Prepare("SELECT id, first_name, last_name, email, created_at, updated_at FROM users LIMIT 100;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare list users statement: %s", err.Error())
	}

	client.createUserStmt, err = db.Prepare("INSERT INTO users (first_name, last_name, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare create user statement: %s", err.Error())
	}

	client.getUserStmt, err = db.Prepare("SELECT id, first_name, last_name, email, created_at, updated_at FROM users WHERE id = $1;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare get user statement: %s", err.Error())
	}

	client.updateUserStmt, err = db.Prepare("UPDATE users SET first_name = $2, last_name = $3, email = $4, updated_at = $5 WHERE id = $1 RETURNING id, first_name, last_name, email, created_at, updated_at;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare update user statement: %s", err.Error())
	}

	client.deleteUserStmt, err = db.Prepare("DELETE FROM users WHERE id = $1;")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare update user statement: %s", err.Error())
	}

	return client, nil
}

func (client *PostgresUserClient) ListUsers() ([]*model.User, error) {
	rows, err := client.listUsersStmt.Query()
	if err != nil {
		client.logger.Error("unable to list all users", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	users := make([]*model.User, 0)

	for rows.Next() {
		var (
			user        = &model.User{}
			timeUpdated sql.NullString
		)

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.TimeCreated,
			&timeUpdated,
		); err != nil {
			client.logger.Error("unable to scan user, skipping for now", zap.Error(err))
			continue
		}

		if timeUpdated.Valid {
			updatedAt, err := time.Parse(time.RFC3339, timeUpdated.String)
			if err != nil {
				client.logger.Warn("unable to parse updated_at value", zap.Int("user_id", user.ID), zap.Error(err))
			} else {
				user.TimeUpdated = updatedAt
			}
		}

		users = append(users, user)
	}

	return users, nil
}

func (client *PostgresUserClient) CreateUser(user *model.User) (*model.User, error) {
	row := client.createUserStmt.QueryRow(user.FirstName, user.LastName, user.Email, time.Now())

	var userID int64

	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("unable to scan created user id: %s", err.Error())
	}

	createdUser, err := client.GetUser(int(userID))
	if err != nil {
		return nil, fmt.Errorf("unable to get created user: %s", err.Error())
	}

	return createdUser, nil
}

func (client *PostgresUserClient) GetUser(id int) (*model.User, error) {
	var (
		user        = &model.User{}
		timeUpdated sql.NullString
		err         error
	)

	row := client.getUserStmt.QueryRow(id)

	if err = row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.TimeCreated,
		&timeUpdated,
	); err != nil && err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("unable to scan user [%d]: %s", id, err.Error())
	}

	if timeUpdated.Valid {
		updatedAt, err := time.Parse(time.RFC3339, timeUpdated.String)
		if err != nil {
			client.logger.Warn("unable to parse updated_at value", zap.Int("user_id", user.ID), zap.Error(err))
		} else {
			user.TimeUpdated = updatedAt
		}
	}

	return user, nil
}

func (client *PostgresUserClient) UpdateUser(userInput *model.User) (*model.User, error) {
	row := client.updateUserStmt.QueryRow(userInput.ID, userInput.FirstName, userInput.LastName, userInput.Email, time.Now())

	var (
		user        = &model.User{}
		timeUpdated sql.NullString
		err         error
	)

	if err = row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.TimeCreated,
		&timeUpdated,
	); err != nil {
		return nil, fmt.Errorf("unable to scan user [%d]: %s", userInput.ID, err.Error())
	}

	if timeUpdated.Valid {
		updatedAt, err := time.Parse(time.RFC3339, timeUpdated.String)
		if err != nil {
			client.logger.Warn("unable to parse time_update value", zap.Int("user_id", user.ID), zap.Error(err))
		} else {
			user.TimeUpdated = updatedAt
		}
	}
	return user, nil
}

func (client *PostgresUserClient) DeleteUser(id int) error {
	result, err := client.deleteUserStmt.Exec(id)
	if err != nil {
		return fmt.Errorf("unable to delete user [%d]: %s", id, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for user [%d]: %s", id, err.Error())
	}

	if rowsAffected != int64(1) {
		return fmt.Errorf("deleted 0 or more than one user requested")
	}

	return nil
}
