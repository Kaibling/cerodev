package service

import (
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
		return nil, err
	}

	newToken := &model.Token{
		UserID: userID,
		Token:  tokenKey,
	}

	return s.repo.Create(newToken)
}

func (s *TokenService) CreateUnsafe(token *model.Token) (*model.Token, error) {
	return s.repo.Create(token)
}

func (s *TokenService) Delete(token string) error {
	return s.repo.Delete(token)
}

func (s *TokenService) Validate(token string) (bool, error) {
	t, err := s.GetByTokenKey(token)
	if err != nil {
		return false, err
	}

	return t != nil, nil
}

func (s *TokenService) GetByTokenKey(token string) (*model.Token, error) {
	return s.repo.Get(token)
}
