package sqliterepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kaibling/apiforge/apierror"
	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
	_ "modernc.org/sqlite" // import sqlite3
)

type UserRepo struct {
	db *sql.DB
	l  log.Writer
}

func NewUserRepo(db *sql.DB, l log.Writer) *UserRepo {
	return &UserRepo{db: db, l: l.Named("repo_user")}
}

func (r *UserRepo) Create(user *model.User) error {
	// Insert new user into the users table
	_, err := r.db.Exec("INSERT INTO users (id,username, password) VALUES (?,?, ?)", user.ID, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepo) Delete(id string) error {
	// Delete a user by ID
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepo) Update(user *model.User) error {
	// Update user's username and password
	_, err := r.db.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", user.Username, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByID(id string) (*model.User, error) {
	// Get user by ID
	row := r.db.QueryRow("SELECT id, username FROM users WHERE id = ?", id)

	var user model.User

	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user found with ID %s: %w", id, err)
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) SecureGetByUsername(username string) (*model.User, error) {
	// Get user by username
	row := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)

	var user model.User

	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user found with username %s :%w", username, err)
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetAll() ([]*model.User, error) {
	// Get all users
	rows, err := r.db.Query("SELECT id, username FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", rows.Err())
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.l.Error("failed to close rows", err)
		}
	}()

	var users []*model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepo) GetUserByToken(token string) (*model.User, error) {
	// Get user by token
	rows, err := r.db.Query("SELECT u.id,u.username, t.token FROM tokens t JOIN users u ON t.user_id = u.id WHERE u.id = (SELECT user_id FROM tokens WHERE token = ?)", token) //nolint:lll
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to fetch tokens: %w", rows.Err())
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.l.Error("failed to close rows", err)
		}
	}()

	user := model.User{} //nolint:exhaustruct
	tokens := []string{}
	found := false

	for rows.Next() {
		found = true

		var t string

		if err := rows.Scan(&user.ID, &user.Username, &t); err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}

		tokens = append(tokens, t)
	}

	if !found {
		return nil, fmt.Errorf("no user found: %w", apierror.ErrDataNotFound)
	}

	user.Tokens = tokens

	return &user, nil
}
