package helper

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"strconv"
)

func ExtractUserID(ctx context.Context) int {
	_, claims, err := jwtauth.FromContext(ctx)
	PanicIfError(err, "Failed to parse JWT from context")

	if claims == nil {
		panic("JWT token is missing in the context")
	}

	userIdVal, ok := claims["userId"]
	if !ok || userIdVal == nil {
		panic("userId is missing or invalid in JWT claims")
	}

	var userId int
	switch v := userIdVal.(type) {
	case string:
		userId, err = strconv.Atoi(v)
		PanicIfError(err, "Failed to convert userId from string to int")

	case float64:
		userId = int(v)
	}

	return userId
}
