package service

import (
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

func (s *UserService) GetByID(id string) (model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (s *UserService) GetAll() ([]model.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	user.Password = hashedPassword

	return s.userRepo.Create(user)
}

func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) Update(user model.User) error {
	return s.userRepo.Update(&user)
}

func (s *UserService) CheckToken(token string) (*model.User, error) {
	return s.userRepo.GetUserByToken(token)
}

func (s *UserService) GetUnsafeByUsername(username string) (*model.User, error) {
	return s.userRepo.GetUnsafeByUsername(username)
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
		return nil, err
	}

	return &model.LoginResponse{
		Username: user.Username,
		UserID:   user.ID,
		Token:    newToken.Token,
	}, nil
}
