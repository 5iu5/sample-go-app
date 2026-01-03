package routes

import (
	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/handlers/auth"
	"github.com/CVWO/sample-go-app/internal/handlers/comments"
	"github.com/CVWO/sample-go-app/internal/handlers/posts"
	"github.com/CVWO/sample-go-app/internal/handlers/topics"
	"github.com/CVWO/sample-go-app/internal/handlers/users"
	"github.com/CVWO/sample-go-app/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetProtectedRoutes(pool *pgxpool.Pool) func(r chi.Router) {
	return func(r chi.Router) {

		r.Use(middlewares.MiddlewareAuth)

		r.Get("/validate", func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value("user_id")
			userID, ok := val.(int)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			fmt.Fprintf(w, "Hello user %d!", userID)
		})

		r.Post("/auth/logout", auth.LogoutHandler) //Logout

		r.Get("/auth/me", func(w http.ResponseWriter, r *http.Request) {
			auth.CurrentUserHandler(w, r, pool)
		}) //gets current user

		r.Delete("/auth/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
			auth.DeleteUserHandler(w, r, pool)
		})

		/* Users */
		r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
			users.ListUsersHandler(w, r, pool)
		}) //List all users

		r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			users.GetUserHandler(w, r, pool)
		}) //Get user by ID

		r.Get("/users/{id}/posts", func(w http.ResponseWriter, r *http.Request) {
			users.ListUserPostsHandler(w, r, pool)
		}) //List posts made by user

		r.Get("/users/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
			users.ListUserCommentsHandler(w, r, pool)
		}) //List comments made by user

		/* Posts */
		r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
			posts.ListPostsHandler(w, r, pool)
		}) //List all main/thread posts

		r.Get("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
			posts.GetPostHandler(w, r, pool)
		}) //Get a post by ID

		r.Post("/posts", func(w http.ResponseWriter, r *http.Request) {
			posts.AddPostHandler(w, r, pool)
		})

		r.Patch("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
			posts.EditPostHandler(w, r, pool)
		})

		r.Delete("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
			posts.DeletePostHandler(w, r, pool)
		})

		/* Comments */
		r.Get("/posts/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
			comments.ListPostCommentsHandler(w, r, pool)
		}) //List all comments for a post ID

		r.Get("/posts/{id}/comments/count", func(w http.ResponseWriter, r *http.Request) {
			comments.GetCommentCount(w, r, pool)
		})

		r.Post("/posts/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
			comments.AddCommentHandler(w, r, pool)
		}) //Add comment to a post

		r.Patch("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
			comments.EditCommentHandler(w, r, pool)
		}) //Edit comment to a post

		r.Delete("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
			comments.DeleteCommentHandler(w, r, pool)
		})

		/* Topics */
		r.Get("/topics", func(w http.ResponseWriter, r *http.Request) {
			topics.ListTopicsHandler(w, r, pool)
		}) //List all topics

	}
}
