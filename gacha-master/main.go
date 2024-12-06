package main

import (
	"gacha-master/app"
	"gacha-master/controller"
	"gacha-master/helper"
	"gacha-master/repository"
	"gacha-master/service"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
)

func main() {
	dbpool := app.NewDB()
	defer dbpool.Close()

	validate := validator.New()

	gachaSystemRepository := repository.NewGachaSystemRepository(dbpool)
	rarityRepository := repository.NewRarityRepository(dbpool)
	characterRepository := repository.NewCharacterRepository(dbpool)

	gachaSystemService := service.NewGachaSystemService(gachaSystemRepository, rarityRepository, characterRepository, validate)
	rarityService := service.NewRarityService(rarityRepository, gachaSystemRepository, validate)
	characterService := service.NewCharacterService(characterRepository, rarityRepository, gachaSystemRepository, validate)
	uploaderService := service.NewUploaderServiceImpl()
	defer uploaderService.Close()

	gachaSystemController := controller.NewGachaSystemController(gachaSystemService, rarityService, characterService, uploaderService)
	rarityController := controller.NewRarityController(rarityService)
	characterController := controller.NewCharacterController(characterService, uploaderService)

	router := app.NewRouter(gachaSystemController, rarityController, characterController)

	server := http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err, "Failed to start server")
}
