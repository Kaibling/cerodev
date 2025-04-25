package dbrepo

import (
	"context"
	"database/sql"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/repo/sqlcrepo"
)

type UserRepo struct {
	ctx      context.Context
	sqlcRepo *sqlcrepo.Queries
	l        log.Writer
}

func NewUserRepo(ctx context.Context, db *sql.DB, l log.Writer) *UserRepo {
	return &UserRepo{ctx: ctx, sqlcRepo: sqlcrepo.New(db), l: l.Named("repo_user")}
}

func (r *UserRepo) GetByID(id string) (*model.User, error) {
	rows, err := r.sqlcRepo.GetUserByID(r.ctx, id)
	if err != nil {
		r.l.Error("failed to get user by username", err)

		return nil, err
	}

	user := model.User{} //nolint:exhaustruct
	tokens := []string{}

	for _, row := range rows {
		user.ID = row.ID
		user.Username = row.Username

		if row.Token.Valid {
			tokens = append(tokens, row.Token.String)
		}
	}

	if user.ID == "" {
		r.l.Warn("no user found with username")

		return nil, sql.ErrNoRows
	}

	user.Tokens = tokens

	return &user, nil
}

func (r *UserRepo) GetUnsafeByUsername(username string) (*model.User, error) {
	rows, err := r.sqlcRepo.GetUnsafeUserByUsername(r.ctx, username)
	if err != nil {
		r.l.Error("failed to get user by username", err)

		return nil, err
	}

	user := model.User{} //nolint:exhaustruct
	tokens := []string{}

	for _, row := range rows {
		user.ID = row.ID
		user.Username = row.Username
		user.Password = row.Password

		if row.Token.Valid {
			tokens = append(tokens, row.Token.String)
		}
	}

	if user.ID == "" {
		r.l.Warn("no user found with username")

		return nil, sql.ErrNoRows
	}

	user.Tokens = tokens

	return &user, nil
}

func (r *UserRepo) GetAll() ([]*model.User, error) {
	rows, err := r.sqlcRepo.GetAllUsers(r.ctx)
	if err != nil {
		r.l.Error("failed to get all users", err)

		return nil, err
	}

	users := []*model.User{}
	currentTokens := []string{}
	currentUser := model.User{} //nolint:exhaustruct

	for _, row := range rows {
		if currentUser.ID != row.ID && currentUser.ID != "" {
			currentUser := model.User{}
			currentUser.ID = row.ID
			currentUser.Username = row.Username
			currentUser.Tokens = currentTokens
			users = append(users, &currentUser)
			currentTokens = []string{}
		}

		if row.Token.Valid {
			currentTokens = append(currentTokens, row.Token.String)
		}
	}

	if len(users) == 0 {
		r.l.Warn("no users found")

		return nil, sql.ErrNoRows
	}

	return users, nil
}

func (r *UserRepo) Create(user *model.User) (*model.User, error) {
	userID, err := r.sqlcRepo.CreateUser(r.ctx, sqlcrepo.CreateUserParams{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		r.l.Error("failed to create user", err)

		return nil, err
	}

	return r.GetByID(userID)
}

func (r *UserRepo) Delete(id string) error {
	err := r.sqlcRepo.DeleteUser(r.ctx, id)
	if err != nil {
		r.l.Error("failed to delete user", err)

		return err
	}

	return nil
}

func (r *UserRepo) Update(user *model.User) error {
	err := r.sqlcRepo.UpdateUser(r.ctx, sqlcrepo.UpdateUserParams{
		Username: user.Username,
		Password: user.Password,
		ID:       user.ID,
	})
	if err != nil {
		r.l.Error("failed to update user", err)

		return err
	}

	return nil
}

func (r *UserRepo) GetUserByToken(token string) (*model.User, error) {
	rows, err := r.sqlcRepo.GetUserByToken(r.ctx, token)
	if err != nil {
		r.l.Error("failed to get user by token", err)

		return nil, err
	}

	user := model.User{}
	tokens := []string{}

	for _, row := range rows {
		user.ID = row.ID
		user.Username = row.Username
		tokens = append(tokens, row.Token)
	}

	if user.ID == "" {
		r.l.Warn("no user found with username")

		return nil, sql.ErrNoRows
	}

	user.Tokens = tokens

	return &user, nil
}
