package users

import (
	"context"
	"encoding/json"

	// "fmt"
	"log"
	"net/http"
	"strconv"

	// "github.com/CVWO/sample-go-app/internal/api"
	// users "github.com/CVWO/sample-go-app/internal/dataaccess"
	// "github.com/CVWO/sample-go-app/internal/database"
	"github.com/CVWO/sample-go-app/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/pkg/errors"
)

const (
	// ListUsers = "users.HandleList"

	SuccessfulListUsersMessage = "Successfully listed users"
	ErrRetrieveDatabase        = "Failed to retrieve database in %s"
	ErrRetrieveUsers           = "Failed to retrieve users in %s"
	ErrEncodeView              = "Failed to retrieve users in %s"
)

// func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
// 	db, err := database.GetDB()

// 	if err != nil {
// 		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveDatabase, ListUsers))
// 	}

// 	users, err := users.List(db)
// 	if err != nil {
// 		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveUsers, ListUsers))
// 	}

// 	data, err := json.Marshal(users)
// 	if err != nil {
// 		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListUsers))
// 	}

// 	return &api.Response{
// 		Payload: api.Payload{
// 			Data: data,
// 		},
// 		Messages: []string{SuccessfulListUsersMessage},
// 	}, nil
// }

func ListUsersHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//Lists all users
	rows, err := pool.Query(context.Background(), "SELECT user_id, username, email, hashed_password, created_at from users")
	if err != nil {
		log.Println("Error querying users:", err)
		http.Error(w, "Error querying users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		users = append(users, user)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println("JSON encode error in ListUsersHandler:", err)
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//Get a user by its ID
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	var user models.User
	err = pool.QueryRow(context.Background(), "Select user_id, username, email, hashed_password, created_at from users WHERE user_id = $1", userID).Scan(&user.UserID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error querying user in GetUserHandler: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	log.Println("successfully retrieved user")
}

func ListUserPostsHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//Lists all posts made by a user
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error querying posts from database: ", err)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	rows, err := pool.Query(context.Background(), "SELECT post_id, topic_id, user_id, text_title, text_body, created_at, updated_at from posts WHERE user_id = $1 ORDER BY created_at ASC", userID)
	if err != nil {
		http.Error(w, "Error querying posts from db", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err = rows.Scan(&post.PostID, &post.TopicID, &post.UserID, &post.TextTitle, &post.TextBody, &post.CreatedAt, &post.UpdatedAt); err != nil {
			log.Println("Error scanning row in ListUserPostsHandler: ", err)
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
		}
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Println("JSON encode error in ListUserPostsHandler: ", err)
		http.Error(w, "Error encoding to json", http.StatusInternalServerError)
		return
	}
}

func ListUserCommentsHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//Lists all comments by a user
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	rows, err := pool.Query(context.Background(), "SELECT comment_id, post_id, parent_id, user_id, text_body, created_at, updated_at from comments WHERE user_id=$1 AND is_deleted=false", userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.ParentID, &comment.UserID, &comment.TextBody, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		log.Println("JSON encode error in ListPostsHandler: ", err)
		return
	}
}
