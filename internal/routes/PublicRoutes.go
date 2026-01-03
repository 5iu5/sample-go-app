package routes

import (
	// "encoding/json"
	"encoding/json"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/handlers/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//routes that do not require auth

func GetPublicRoutes(pool *pgxpool.Pool) func(r chi.Router) {

	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(`Running / route`)

		})

		r.Post("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
			auth.SignupHandler(w, r, pool)
		})

		r.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
			auth.LoginHandler(w, r, pool) //Login
		})

	}
}
