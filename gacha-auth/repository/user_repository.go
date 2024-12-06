package repository

import (
	"context"
	"gacha-auth/helper"
	"gacha-auth/model/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User)
	FindByName(ctx context.Context, name string) *domain.User
	FindByUsername(ctx context.Context, username string) *domain.User
}

type UserRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewUserRepository(dbpool *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, user *domain.User) {
	query := `INSERT INTO users (name, username, password) 
              VALUES ($1, $2, $3) RETURNING id`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	var id int
	err = tx.QueryRow(ctx, query, user.Name, user.Username, user.Password).Scan(&id)
	helper.PanicIfError(err, helper.ErrUserNotFound)

	user.Id = id
}

func (repository *UserRepositoryImpl) FindByName(ctx context.Context, name string) *domain.User {
	query := `SELECT id, name, username, password FROM users WHERE name = $1`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, name)

	var user domain.User
	err = row.Scan(&user.Id, &user.Name, &user.Username, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}

func (repository *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) *domain.User {
	query := `SELECT id, name, username, password FROM users WHERE username = $1`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, username)

	var user domain.User
	err = row.Scan(&user.Id, &user.Name, &user.Username, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}
