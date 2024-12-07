package repository

import (
	"context"
	"gacha-pull/helper"
	"gacha-pull/model/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GachaSystemRepository interface {
	FindByEndpointId(ctx context.Context, endpointId string) *domain.GachaSystem
}

type GachaSystemRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewGachaSystemRepository(dbpool *pgxpool.Pool) GachaSystemRepository {
	return &GachaSystemRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *GachaSystemRepositoryImpl) FindByEndpointId(ctx context.Context, endpointId string) *domain.GachaSystem {
	query := `SELECT id, name, endpoint_id
			FROM gacha_system
			WHERE endpoint_id = $1`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, endpointId)

	var gachaSystem domain.GachaSystem
	err = row.Scan(&gachaSystem.Id, &gachaSystem.Name, &gachaSystem.EndpointId)
	if err != nil {
		return nil
	}

	return &gachaSystem
}
