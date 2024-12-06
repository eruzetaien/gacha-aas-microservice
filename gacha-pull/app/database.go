package app

import (
	"context"
	"gacha-pull/helper"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func NewDB() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	helper.PanicIfError(err, "Failed to connect to database")

	return dbpool
}
