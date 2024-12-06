package repository

import (
	"context"
	"gacha-master/helper"
	"gacha-master/model/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type RarityRepository interface {
	Save(ctx context.Context, rarity *domain.Rarity)
	FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *domain.Rarity
	FindByNameAndGachaSystemId(ctx context.Context, name string, gachaSystemId int) *domain.Rarity
	FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []domain.Rarity
	Update(ctx context.Context, rarity *domain.Rarity)
	Delete(ctx context.Context, id int, gachaSystemId int)
}

type RarityRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewRarityRepository(dbpool *pgxpool.Pool) RarityRepository {
	return &RarityRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *RarityRepositoryImpl) Save(ctx context.Context, rarity *domain.Rarity) {
	query := `INSERT INTO rarity (gacha_system_id, name, chance) 
				VALUES ($1, $2, $3) RETURNING id`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	var id int
	err = tx.QueryRow(ctx, query, rarity.GachaSystemId, rarity.Name, rarity.Chance).Scan(&id)
	helper.PanicIfError(err, helper.ErrUserNotFound)

	rarity.Id = id
}

func (repository *RarityRepositoryImpl) FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *domain.Rarity {
	query := `SELECT id, name, chance, gacha_system_id FROM rarity WHERE id = $1 AND gacha_system_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, id, gachaSystemId)

	var rarity domain.Rarity
	err = row.Scan(&rarity.Id, &rarity.Name, &rarity.Chance, &rarity.GachaSystemId)
	if err != nil {
		return nil
	}

	return &rarity
}

func (repository *RarityRepositoryImpl) FindByNameAndGachaSystemId(ctx context.Context, name string, gachaSystemId int) *domain.Rarity {
	query := `SELECT id, name, chance, gacha_system_id FROM rarity WHERE LOWER(name) = LOWER($1) AND gacha_system_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, name, gachaSystemId)

	var rarity domain.Rarity
	err = row.Scan(&rarity.Id, &rarity.Name, &rarity.Chance, &rarity.GachaSystemId)
	if err != nil {
		return nil
	}

	return &rarity
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

func (repository *RarityRepositoryImpl) Update(ctx context.Context, rarity *domain.Rarity) {
	query := `UPDATE rarity 
	          SET name = $1, chance = $2 
	          WHERE id = $3 AND gacha_system_id = $4`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	_, err = tx.Exec(ctx, query, rarity.Name, rarity.Chance, rarity.Id, rarity.GachaSystemId)
	helper.PanicIfError(err, "Failed to update rarity")
}

func (repository *RarityRepositoryImpl) Delete(ctx context.Context, id int, gachaSystemId int) {
	query := `DELETE FROM rarity WHERE id = $1 AND gacha_system_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	_, err = tx.Exec(ctx, query, id, gachaSystemId)
	helper.PanicIfError(err, "Failed to delete rarity")
}
