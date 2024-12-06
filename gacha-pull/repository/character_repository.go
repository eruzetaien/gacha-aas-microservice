package repository

import (
	"context"
	"database/sql"
	"gacha-pull/helper"
	"gacha-pull/model/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type CharacterRepository interface {
	FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []domain.Character
}

type CharacterRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewCharacterRepository(dbpool *pgxpool.Pool) CharacterRepository {
	return &CharacterRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *CharacterRepositoryImpl) FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []domain.Character {
	query := `SELECT id, name, image_url, rarity_id, gacha_system_id FROM character WHERE gacha_system_id = $1`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	rows, err := tx.Query(ctx, query, gachaSystemId)
	if err != nil {
		log.Printf("Error while querying characters: %v", err)
	}
	defer rows.Close()

	var characters []domain.Character
	for rows.Next() {
		var character domain.Character
		var imageUrl sql.NullString

		err = rows.Scan(&character.Id, &character.Name, &imageUrl, &character.RarityId, &character.GachaSystemId)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if imageUrl.Valid {
			character.ImageUrl = imageUrl.String
		} else {
			character.ImageUrl = ""
		}
		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error while scanning characters: %v", err)
	}

	return characters
}
