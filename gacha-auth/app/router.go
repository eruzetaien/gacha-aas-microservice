package app

import (
	"gacha-auth/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	"net/http"
)

func NewRouter(tokenAuth *jwtauth.JWTAuth, controller controller.UserController) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(RecoverMiddleware)

	// Protected routes
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator(tokenAuth))
	})

	// Public routes
	router.Group(func(router chi.Router) {
		router.Route("/api/v1", func(subRouter chi.Router) {
			subRouter.Post("/register", controller.Register)
			subRouter.Post("/login", controller.Login)
		})
	})

	return router
}
