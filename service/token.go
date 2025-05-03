package service

import (
	"fmt"

	"github.com/kaibling/cerodev/config"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/crypto"
)

type tokenrepo interface {
	Create(token *model.Token) (*model.Token, error)
	Get(tokenKey string) (*model.Token, error)
	GetByUserID(userID string) ([]*model.Token, error)
	Delete(token string) error
}

type TokenService struct {
	repo tokenrepo
	cfg  config.Configuration
}

func NewTokenService(repo tokenrepo, cfg config.Configuration) *TokenService {
	return &TokenService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *TokenService) CreateForUser(userID string) (*model.Token, error) {
	tokenKey, err := crypto.GenerateToken(s.cfg.TokenLength)
	if err != nil {
		return nil, fmt.Errorf("failed to GenerateToken: %w", err)
	}

	newToken := &model.Token{
		UserID: userID,
		Token:  tokenKey,
	}

	val, err := s.repo.Create(newToken)

	return HandleError[*model.Token](val, err, "failed to db Create")
}

func (s *TokenService) CreateUnsafe(token *model.Token) (*model.Token, error) {
	val, err := s.repo.Create(token)

	return HandleError[*model.Token](val, err, "failed to db Create")
}

func (s *TokenService) Delete(token string) error {
	if err := s.repo.Delete(token); err != nil {
		return fmt.Errorf("failed to db Delete: %w", err)
	}

	return nil
}

func (s *TokenService) Validate(token string) (bool, error) {
	t, err := s.GetByTokenKey(token)
	if err != nil {
		return false, fmt.Errorf("failed to GetByTokenKey: %w", err)
	}

	return t != nil, nil
}

func (s *TokenService) GetByTokenKey(token string) (*model.Token, error) {
	val, err := s.repo.Get(token)

	return HandleError[*model.Token](val, err, "failed to db Get")
}
