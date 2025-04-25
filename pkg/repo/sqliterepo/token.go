package sqliterepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
)

type Token struct {
	UserID string
	Token  string
}
type TokenRepo struct {
	db *sql.DB
	l  log.Writer
}

func NewTokenRepo(db *sql.DB, l log.Writer) *TokenRepo {
	return &TokenRepo{db: db, l: l.Named("repo_token")}
}

func (t *TokenRepo) Create(token *model.Token) (*model.Token, error) {
	// Insert token for a user
	_, err := t.db.Exec("INSERT INTO tokens (user_id, token) VALUES (?, ?)", token.UserID, token.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return t.Get(token.Token)
}

func (t *TokenRepo) Delete(token string) error {
	// Delete token by user ID
	_, err := t.db.Exec("DELETE FROM tokens WHERE token = ?", token)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func (t *TokenRepo) Get(tokenKey string) (*model.Token, error) {
	// Get token by user ID
	var token Token

	row := t.db.QueryRow("SELECT user_id, token FROM tokens WHERE token = ?", tokenKey)

	if row.Err() != nil {
		return nil, fmt.Errorf("failed to query token: %w", row.Err())
	}

	err := row.Scan(&token.UserID, &token.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no token found for user %s: %w", tokenKey, err)
		}

		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return marshalToken(token), nil
}

func (t *TokenRepo) GetByUserID(userID int) ([]*model.Token, error) {
	// Get token by user ID
	rows, err := t.db.Query("SELECT user_id, token FROM tokens WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tokens: %w", err)
	}

	defer rows.Close()

	tokens := []*model.Token{}

	for rows.Next() {
		var token Token
		if err := rows.Scan(&token.UserID, &token.Token); err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}

		tokens = append(tokens, marshalToken(token))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tokens, nil
}

func marshalToken(token Token) *model.Token {
	return &model.Token{
		UserID: token.UserID,
		Token:  token.Token,
	}
}
