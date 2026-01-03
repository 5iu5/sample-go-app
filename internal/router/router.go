package router

import (
	"github.com/CVWO/sample-go-app/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(pool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		// MaxAge:           300, // 5 minutes
	}))

	setUpRoutes(r, pool)
	return r
}

func setUpRoutes(r chi.Router, pool *pgxpool.Pool) {

	//routes that require Auth (middleware)
	r.Group(routes.GetProtectedRoutes(pool))

	//routes that do not require auth
	r.Group(routes.GetPublicRoutes(pool))

}
