package app

import (
	"gacha-master/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"os"

	"net/http"
)

func NewRouter(
	gachaSystemController controller.GachaSystemController,
	rarityController controller.RarityController,
	characterController controller.CharacterController,
) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(RecoverMiddleware)

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)

	// Protected routes
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator(tokenAuth))

		router.Route("/api/v1/gacha", func(subRouter chi.Router) {
			subRouter.Post("/create", gachaSystemController.Create)
			subRouter.Delete("/id/{gachaSystemId}", gachaSystemController.Delete)
			subRouter.Get("/id/{gachaSystemId}", gachaSystemController.FindById)
			subRouter.Get("/all", gachaSystemController.FindAll)
			subRouter.Get("/id/all", gachaSystemController.FindAll)

			subRouter.Post("/character/create", characterController.Create)
			subRouter.Patch("/character/update", characterController.Update)
			subRouter.Delete("/id/{gachaSystemId}/character/{characterId}", characterController.Delete)
			subRouter.Get("/id/{gachaSystemId}/character/{characterId}", characterController.GetById)

			subRouter.Put("/rarity/update", rarityController.Update)
			subRouter.Post("/rarity/create", rarityController.Create)
			subRouter.Get("/id/{gachaSystemId}/rarity/all", rarityController.GetAll)
			subRouter.Delete("/id/{gachaSystemId}/rarity/{rarityId}", rarityController.Delete)

		})
	})

	// Public routes
	router.Group(func(router chi.Router) {
		router.Get("/api/v1/gacha/userId/{userId}", gachaSystemController.FindEndpointByNameAndUserId)
	})

	return router
}
