package repository

import (
	"context"
	"gacha-pull/helper"
	"gacha-pull/model/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GachaSystemRepository interface {
	FindByEndpoint(ctx context.Context, endpoint string) *domain.GachaSystem
}

type GachaSystemRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewGachaSystemRepository(dbpool *pgxpool.Pool) GachaSystemRepository {
	return &GachaSystemRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *GachaSystemRepositoryImpl) FindByEndpoint(ctx context.Context, endpoint string) *domain.GachaSystem {
	query := `SELECT id, name, endpoint
			FROM gacha_system
			WHERE endpoint = $1`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, endpoint)

	var gachaSystem domain.GachaSystem
	err = row.Scan(&gachaSystem.Id, &gachaSystem.Name, &gachaSystem.Endpoint)
	if err != nil {
		return nil
	}

	return &gachaSystem
}
