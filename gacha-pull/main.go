package main

import (
	"gacha-pull/app"
	"gacha-pull/controller"
	"gacha-pull/helper"
	"gacha-pull/repository"
	"gacha-pull/service"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
)

func main() {
	dbpool := app.NewDB()
	defer dbpool.Close()

	gachaSystemRepository := repository.NewGachaSystemRepository(dbpool)
	rarityRepository := repository.NewRarityRepository(dbpool)
	characterRepository := repository.NewCharacterRepository(dbpool)

	characterService := service.NewCharacterService(characterRepository, rarityRepository, gachaSystemRepository)
	characterController := controller.NewCharacterController(characterService)

	router := app.NewRouter(characterController)

	server := http.Server{
		Addr:    ":8002",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err, "Failed to start server")
}
