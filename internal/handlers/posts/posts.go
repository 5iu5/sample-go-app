package posts

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/CVWO/sample-go-app/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListPostsHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//List all main/thread posts
	rows, err := pool.Query(context.Background(), "SELECT posts.post_id, posts.topic_id, posts.user_id, posts.text_title, posts.text_body, posts.created_at, posts.updated_at, posts.is_deleted, users.username from posts INNER JOIN users ON posts.user_id = users.user_id ORDER BY posts.created_at ASC")
	if err != nil {
		http.Error(w, "Error querying posts from db", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err = rows.Scan(&post.PostID, &post.TopicID, &post.UserID, &post.TextTitle, &post.TextBody, &post.CreatedAt, &post.UpdatedAt, &post.IsDeleted, &post.Username); err != nil {
			http.Error(w, "Error scanning rows in ListPostsHandler", http.StatusInternalServerError)
			log.Println("Error scanning rows in ListPostsHandler")
			return
		}
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		log.Println("JSON encode error in ListPostsHandler: ", err)
		return
	}

}

func GetPostHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//Get a post by its ID
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Missing post id", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}
	var post models.Post
	err = pool.QueryRow(context.Background(), "SELECT posts.post_id, posts.topic_id, posts.user_id, posts.text_title, posts.text_body, posts.created_at, posts.updated_at, posts.is_deleted, users.username FROM posts INNER JOIN users ON users.user_id=posts.user_id WHERE posts.post_id = $1", postID).Scan(&post.PostID, &post.TopicID, &post.UserID, &post.TextTitle, &post.TextBody, &post.CreatedAt, &post.UpdatedAt, &post.IsDeleted, &post.Username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error querying user in GetUserHandler: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}

func AddPostHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {

	var form models.AddPostForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(w, "Error decoding post from json to GO struct", http.StatusInternalServerError)
		log.Println("Error converting add post input from json to GO struct")
		return
	}

	var postID int
	err = pool.QueryRow(context.Background(), "INSERT INTO posts (topic_id, user_id, text_title, text_body) VALUES ($1, $2, $3, $4) RETURNING post_id", form.TopicID, form.UserID, form.TextTitle, form.TextBody).Scan(&postID)
	if err != nil {
		http.Error(w, "Error inserting new post into database", http.StatusInternalServerError)
		log.Println("Error inserting new post into database", err)
		return
	}

	response := map[string]interface{}{
		"post_id": postID,
		"message": "Post created successfully",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created a new post with postID: %d", postID)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	type EditForm struct {
		TextBody string `json:"text_body"`
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Missing post id", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}
	var form EditForm
	err = json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		log.Println("Error decoding JSON")
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	_, err = pool.Exec(context.Background(), "UPDATE posts SET text_body = $1, updated_at = NOW() WHERE post_id = $2", form.TextBody, postID)
	if err != nil {
		log.Println("Error updating post in db")
		http.Error(w, "Error updating a post in table posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"post_id": postID,
		"message": "Post updated successfully",
	}
	json.NewEncoder(w).Encode(response)
	log.Printf("Successfully edit a post with postID: %d", postID)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Missing post id", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	_, err = pool.Exec(context.Background(), `UPDATE posts SET is_deleted=true, text_body='[deleted]' WHERE post_id=$1`, postID)
	if err != nil {
		http.Error(w, "Error soft-deleting a post", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"post_id": postID,
		"message": "Post soft-deleted",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response for deleting post", http.StatusInternalServerError)
		return
	}

}
