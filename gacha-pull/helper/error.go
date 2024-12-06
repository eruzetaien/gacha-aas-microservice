package helper

import "log"

func PanicIfError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
		panic(err)
	}
}

const (
	ErrBeginTransaction = "Failed to begin database transaction"
)
