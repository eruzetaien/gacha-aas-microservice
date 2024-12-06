package helper

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CommitOrRollback(tx pgx.Tx, ctx context.Context) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback(ctx)
		PanicIfError(errorRollback, "Failed to rollback database transaction")
		panic(err)
	} else {
		errorCommit := tx.Commit(ctx)
		PanicIfError(errorCommit, "Failed to commit database transaction")
	}
}
