package main

import (
	"gacha-auth/app"
	"gacha-auth/controller"
	"gacha-auth/helper"
	"gacha-auth/repository"
	"gacha-auth/service"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
)

func main() {
	dbpool := app.NewDB()
	defer dbpool.Close()

	validate := validator.New()
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)

	userRepository := repository.NewUserRepository(dbpool)
	userService := service.NewUserService(userRepository, validate, tokenAuth)
	userController := controller.NewUserController(userService)

	router := app.NewRouter(tokenAuth, userController)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err, "Failed to start server")
}
