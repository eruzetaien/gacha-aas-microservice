package repository

import (
	"context"
	"gacha-master/helper"
	"gacha-master/model/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GachaSystemRepository interface {
	Save(ctx context.Context, gachaSystem *domain.GachaSystem)
	FindByNameAndUserId(ctx context.Context, name string, userId int) *domain.GachaSystem
	FindByIdAndUserId(ctx context.Context, id int, userId int) *domain.GachaSystem
	FindAllByUserId(ctx context.Context, userId int) []domain.GachaSystem
	Delete(ctx context.Context, gachaSystemId int)
}

type GachaSystemRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewGachaSystemRepository(dbpool *pgxpool.Pool) GachaSystemRepository {
	return &GachaSystemRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *GachaSystemRepositoryImpl) Save(ctx context.Context, gachaSystem *domain.GachaSystem) {
	query := `INSERT INTO gacha_system (name, user_id, endpoint_id) 
				VALUES ($1, $2, $3) RETURNING id`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	var id int
	err = tx.QueryRow(ctx, query, gachaSystem.Name, gachaSystem.UserId, gachaSystem.EndpointId).Scan(&id)
	helper.PanicIfError(err, helper.ErrUserNotFound)

	gachaSystem.Id = id
}

func (repository *GachaSystemRepositoryImpl) FindByNameAndUserId(ctx context.Context, name string, userId int) *domain.GachaSystem {
	query := `SELECT id, name, endpoint_id
			FROM gacha_system
			WHERE LOWER(name) = LOWER($1) AND user_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, name, userId)

	var gachaSystem domain.GachaSystem
	err = row.Scan(&gachaSystem.Id, &gachaSystem.Name, &gachaSystem.EndpointId)
	if err != nil {
		return nil
	}

	return &gachaSystem
}

func (repository *GachaSystemRepositoryImpl) FindByIdAndUserId(ctx context.Context, id int, userId int) *domain.GachaSystem {
	query := `SELECT id, name, endpoint_id
			FROM gacha_system
			WHERE id = $1 AND user_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, id, userId)

	var gachaSystem domain.GachaSystem
	err = row.Scan(&gachaSystem.Id, &gachaSystem.Name, &gachaSystem.EndpointId)
	if err != nil {
		return nil
	}

	return &gachaSystem
}

func (repository *GachaSystemRepositoryImpl) FindAllByUserId(ctx context.Context, userId int) []domain.GachaSystem {
	query := `SELECT id, name, endpoint_id
              FROM gacha_system
              WHERE user_id = $1`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	rows, err := tx.Query(ctx, query, userId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var gachaSystems []domain.GachaSystem
	for rows.Next() {
		var gachaSystem domain.GachaSystem
		if err := rows.Scan(&gachaSystem.Id, &gachaSystem.Name, &gachaSystem.EndpointId); err != nil {
			return nil
		}
		gachaSystems = append(gachaSystems, gachaSystem)
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	return gachaSystems
}

func (repository *GachaSystemRepositoryImpl) Delete(ctx context.Context, gachaSystemId int) {
	query := `DELETE FROM gacha_system WHERE id = $1 `

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	_, err = tx.Exec(ctx, query, gachaSystemId)
	helper.PanicIfError(err, "Failed to delete gacha system")
}
