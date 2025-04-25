package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/repo/sqlcrepo"
)

type TokenRepo struct {
	ctx      context.Context
	sqlcRepo *sqlcrepo.Queries
	l        log.Writer
}

func NewTokenRepo(ctx context.Context, db *sql.DB, l log.Writer) *TokenRepo {
	return &TokenRepo{ctx: ctx, sqlcRepo: sqlcrepo.New(db), l: l.Named("repo_token")}
}

func (r *TokenRepo) Create(token *model.Token) (*model.Token, error) {
	err := r.sqlcRepo.CreateToken(r.ctx, sqlcrepo.CreateTokenParams{
		Token:  token.Token,
		UserID: token.UserID,
	})
	if err != nil {
		r.l.Error("failed to create token", err)

		return nil, err
	}

	return r.Get(token.Token)
}

func (r *TokenRepo) Get(tokenKey string) (*model.Token, error) {
	token, err := r.sqlcRepo.GetToken(r.ctx, tokenKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.l.Info("no token found")

			return nil, fmt.Errorf("no token found: %w", err)
		}

		r.l.Error("failed to get token", err)

		return nil, err
	}

	return &model.Token{
		Token:  token.Token,
		UserID: token.UserID,
	}, nil
}

func (r *TokenRepo) GetByUserID(userID string) ([]*model.Token, error) {
	tokens, err := r.sqlcRepo.GetTokenByUserID(r.ctx, userID)
	if err != nil {
		r.l.Error("failed to get tokens by user ID", err)

		return nil, err
	}

	result := []*model.Token{}
	for _, token := range tokens {
		result = append(result, &model.Token{
			Token:  token.Token,
			UserID: token.UserID,
		})
	}

	return result, nil
}

func (r *TokenRepo) Delete(token string) error {
	err := r.sqlcRepo.DeleteToken(r.ctx, token)
	if err != nil {
		r.l.Error("failed to delete token", err)

		return err
	}

	return nil
}
