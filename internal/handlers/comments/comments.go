package comments

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

func ListPostCommentsHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//List all comments for a post
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
	rows, err := pool.Query(context.Background(), `
		SELECT comments.comment_id, comments.post_id, comments.parent_id, comments.user_id, comments.text_body, comments.created_at, comments.updated_at, comments.is_deleted, users.username FROM comments INNER JOIN users ON comments.user_id = users.user_id
		WHERE comments.post_id = $1 ORDER BY comments.created_at ASC
		`, postID)
	if err != nil {
		log.Println("Error querying comments from database: ", err)
		http.Error(w, "Error querying comments for a post", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.CommentID, &comment.PostID, &comment.ParentID, &comment.UserID, &comment.TextBody, &comment.CreatedAt, &comment.UpdatedAt, &comment.IsDeleted, &comment.Username)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Println("Error scanning row in ListPostCommentsHandler", err)
			return
		}
		comments = append(comments, comment)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		log.Println("JSON encode error in ListUserPostsHandler: ", err)
		http.Error(w, "Error encoding to json", http.StatusInternalServerError)
		return
	}
	log.Println("successfully retrieved comments from database")
}
func GetCommentCount(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
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
	var count int
	err = pool.QueryRow(context.Background(), "SELECT COUNT (comment_id) FROM comments WHERE post_id = $1", postID).Scan(&count)
	if err != nil {
		http.Error(w, "Error counting number of comments in db", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(count)
	if err != nil {
		http.Error(w, "Error encoding to json", http.StatusInternalServerError)
		return
	}
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	var form models.AddCommentForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(w, "Error decoding json when creating post", http.StatusInternalServerError)
		log.Println("Error converting add comment input from json to GO struct")
		return
	}
	var commentID int
	err = pool.QueryRow(context.Background(), "INSERT INTO comments (post_id, parent_id, user_id, text_body) VALUES ($1, $2, $3, $4) RETURNING comment_id", form.PostID, form.ParentID, form.UserID, form.TextBody).Scan(&commentID)
	if err != nil {
		http.Error(w, "Error inserting to database comments table", http.StatusInternalServerError)
		log.Println("Error inserting new comment into database", err)
		return
	}

	response := map[string]interface{}{
		"comment_id": commentID,
		"message":    "comment created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created a new comment with commentID: %d", commentID)

}

func EditCommentHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	type EditForm struct {
		TextBody string `json:"text_body"`
	}
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		http.Error(w, "Comment id is missing", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}
	var form EditForm
	err = json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(w, "Error decoding json when editing comment", http.StatusInternalServerError)
		return
	}
	_, err = pool.Exec(context.Background(), "UPDATE comments SET text_body=$1 WHERE comment_id = $2", form.TextBody, commentID)
	if err != nil {
		http.Error(w, "Error updating comments table in database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"comment_id": commentID,
		"message":    "comment edited successfully",
	}
	json.NewEncoder(w).Encode(response)
	log.Printf("Successfully edit a comment with commentID: %d", commentID)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Missing comment id", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
		return
	}
	var replyCount int
	err = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM comments WHERE parent_id = $1", commentID).Scan(&replyCount)
	if err != nil {
		http.Error(w, "Error checking for replies", http.StatusInternalServerError)
		log.Println("Error checking replies:", err)
		return
	}

	var response map[string]interface{}

	if replyCount > 0 {
		_, err := pool.Exec(context.Background(), `
            UPDATE comments 
            SET is_deleted = true, text_body = '[deleted]'
            WHERE comment_id = $1
        `, commentID)
		if err != nil {
			http.Error(w, "Error performing soft delete", http.StatusInternalServerError)
			log.Println("Error soft deleting comment:", err)
			return
		}
		response = map[string]interface{}{
			"comment_id":  commentID,
			"message":     "Comment soft-deleted (had replies)",
			"soft_delete": true,
		}
	} else {
		result, err := pool.Exec(context.Background(), "DELETE FROM comments WHERE comment_id = $1", commentID)
		if err != nil {
			http.Error(w, "Error deleting comment", http.StatusInternalServerError)
			log.Println("Error hard deleting comment:", err)
			return
		}
		if result.RowsAffected() == 0 {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		response = map[string]interface{}{
			"post_id":     commentID,
			"message":     "Comment deleted successfully",
			"soft_delete": false,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding JSON in DeleteCommentHandler:", err)
	}
	log.Printf("DeleteCommentHandler successfully executed: commentID=%d, soft_delete=%v\n", commentID, replyCount > 0)
}
