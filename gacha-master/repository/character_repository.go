package repository

import (
	"context"
	"database/sql"
	"gacha-master/helper"
	"gacha-master/model/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type CharacterRepository interface {
	Save(ctx context.Context, character *domain.Character)
	FindByNameAndGachaSystemId(ctx context.Context, name string, gachaSystemId int) *domain.Character
	FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *domain.Character
	FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []domain.Character
	Update(ctx context.Context, character *domain.Character)
	InsertImageUrl(ctx context.Context, character *domain.Character)
	Delete(ctx context.Context, id int, gachaSystemId int)
}

type CharacterRepositoryImpl struct {
	Dbpool *pgxpool.Pool
}

func NewCharacterRepository(dbpool *pgxpool.Pool) CharacterRepository {
	return &CharacterRepositoryImpl{
		Dbpool: dbpool,
	}
}

func (repository *CharacterRepositoryImpl) Save(ctx context.Context, character *domain.Character) {
	query := `INSERT INTO character (name, rarity_id, gacha_system_id) 
				VALUES ($1, $2, $3) RETURNING id`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	var id int
	err = tx.QueryRow(ctx, query, character.Name, character.RarityId, character.GachaSystemId).Scan(&id)
	helper.PanicIfError(err, helper.ErrUserNotFound)

	character.Id = id
}

func (repository *CharacterRepositoryImpl) FindByNameAndGachaSystemId(ctx context.Context, name string, gachaSystemId int) *domain.Character {
	query := `SELECT id, name, image_url, rarity_id, gacha_system_id FROM character WHERE LOWER(name) = LOWER($1) AND gacha_system_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	row := tx.QueryRow(ctx, query, name, gachaSystemId)

	return getCharacterFromRow(row)
}

func (repository *CharacterRepositoryImpl) FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *domain.Character {
	query := `SELECT id, name, image_url, rarity_id, gacha_system_id FROM character WHERE id = $1 AND gacha_system_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)
	row := tx.QueryRow(ctx, query, id, gachaSystemId)

	return getCharacterFromRow(row)
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

func (repository *CharacterRepositoryImpl) Update(ctx context.Context, character *domain.Character) {
	query := `UPDATE character 
	          SET name = $1, rarity_id = $2, image_url = $3
	          WHERE id = $4`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	_, err = tx.Exec(ctx, query, character.Name, character.RarityId, character.ImageUrl, character.Id)
	helper.PanicIfError(err, "Failed to update character")
}

func (repository *CharacterRepositoryImpl) InsertImageUrl(ctx context.Context, character *domain.Character) {
	query := `UPDATE character 
	          SET image_url = $1 
	          WHERE id = $2 AND gacha_system_id = $3`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	_, err = tx.Exec(ctx, query, character.ImageUrl, character.Id, character.GachaSystemId)
	helper.PanicIfError(err, "Failed to update image url character")
}

func (repository *CharacterRepositoryImpl) Delete(ctx context.Context, id int, gachaSystemId int) {
	query := `DELETE FROM character WHERE id = $1 AND gacha_system_id = $2`

	tx, err := repository.Dbpool.Begin(ctx)
	helper.PanicIfError(err, helper.ErrBeginTransaction)

	defer helper.CommitOrRollback(tx, ctx)

	_, err = tx.Exec(ctx, query, id, gachaSystemId)
	helper.PanicIfError(err, "Failed to delete character")
}

func getCharacterFromRow(row pgx.Row) *domain.Character {
	var character domain.Character
	var imageUrl sql.NullString

	err := row.Scan(&character.Id, &character.Name, &imageUrl, &character.RarityId, &character.GachaSystemId)
	if err != nil {
		log.Printf("Error scanning row: %v", err)
		return nil
	}

	if imageUrl.Valid {
		character.ImageUrl = imageUrl.String
	} else {
		character.ImageUrl = ""
	}

	return &character
}
