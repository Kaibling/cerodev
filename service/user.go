package service

import (
	"fmt"

	"github.com/kaibling/cerodev/config"
	"github.com/kaibling/cerodev/errs"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/crypto"
	"github.com/kaibling/cerodev/pkg/utils"
)

type userrepo interface {
	GetByID(id string) (*model.User, error)
	GetUnsafeByUsername(username string) (*model.User, error)
	GetAll() ([]*model.User, error)
	Create(user *model.User) (*model.User, error)
	Delete(id string) error
	Update(user *model.User) error
	GetUserByToken(token string) (*model.User, error)
}

type UserService struct {
	userRepo     userrepo
	tokenService *TokenService
	cfg          config.Configuration
}

func NewUserService(repo userrepo, tokenService *TokenService, cfg config.Configuration) *UserService {
	return &UserService{
		userRepo:     repo,
		tokenService: tokenService,
		cfg:          cfg,
	}
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to user GetByID: %w", err)
	}

	return user, nil
}

func (s *UserService) GetAll() ([]model.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to user GetAll: %w", err)
	}

	result := make([]model.User, len(users))
	for i, u := range users {
		result[i] = *u
	}

	return result, nil
}

func (s *UserService) Create(user *model.User) (*model.User, error) {
	user.ID = utils.GenerateULID()

	hashedPassword, err := crypto.HashPassword(user.Password, s.cfg.PasswordCost)
	if err != nil {
		return nil, fmt.Errorf("failed to user HashPassword: %w", err)
	}

	user.Password = hashedPassword

	val, err := s.userRepo.Create(user)

	return HandleError[*model.User](val, err, "failed to db Get")
}

func (s *UserService) Delete(id string) error {
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to db Delete: %w", err)
	}

	return nil
}

func (s *UserService) Update(user model.User) error {
	if err := s.userRepo.Update(&user); err != nil {
		return fmt.Errorf("failed to db Delete: %w", err)
	}

	return nil
}

func (s *UserService) CheckToken(token string) (*model.User, error) {
	val, err := s.userRepo.GetUserByToken(token)

	return HandleError[*model.User](val, err, "failed to db GGetUserByTokent")
}

func (s *UserService) GetUnsafeByUsername(username string) (*model.User, error) {
	val, err := s.userRepo.GetUnsafeByUsername(username)

	return HandleError[*model.User](val, err, "failed to db GetUnsafeByUsername")
}

func (s *UserService) Login(loginRequest *model.LoginRequest) (*model.LoginResponse, error) {
	// check credentials
	user, err := s.userRepo.GetUnsafeByUsername(loginRequest.Username)
	if err != nil {
		return nil, errs.ErrWrongCredentials
	}

	ok, err := crypto.CheckPasswordHash(loginRequest.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errs.ErrWrongCredentials
	}

	// create new token
	newToken, err := s.tokenService.CreateForUser(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to CreateForUser: %w", err)
	}

	return &model.LoginResponse{
		Username: user.Username,
		UserID:   user.ID,
		Token:    newToken.Token,
	}, nil
}
