package repository

import (
	"context"
	"gacha-pull/helper"
	"gacha-pull/model/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type RarityRepository interface {
	FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []domain.Rarity
}

type RarityRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewRarityRepository(dbpool *pgxpool.Pool) RarityRepository {
	return &RarityRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *RarityRepositoryImpl) FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []domain.Rarity {
	query := `SELECT id, name, chance, gacha_system_id FROM rarity WHERE gacha_system_id = $1`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	rows, err := tx.Query(ctx, query, gachaSystemId)
	if err != nil {
		log.Printf("Error while querying rarity: %v", err)
	}
	defer rows.Close()

	var rarities []domain.Rarity
	for rows.Next() {
		var rarity domain.Rarity
		err = rows.Scan(&rarity.Id, &rarity.Name, &rarity.Chance, &rarity.GachaSystemId)

		rarities = append(rarities, rarity)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error while scanning rarities: %v", err)
	}

	return rarities
}
