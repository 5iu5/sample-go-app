package topics

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListTopicsHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	log.Println("running")
	rows, err := pool.Query(context.Background(), "SELECT topic_id, name, description, created_at, is_deleted FROM topics WHERE is_deleted=false")
	if err != nil {
		http.Error(w, "Error querying topics table in db", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var topic models.Topic
		if err = rows.Scan(&topic.TopicID, &topic.TopicName, &topic.Description, &topic.CreatedAt, &topic.IsDeleted); err != nil {
			http.Error(w, "Error scanning rows in ListTopicsHandler", http.StatusInternalServerError)
			return
		}
		topics = append(topics, topic)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(topics)
	if err != nil {
		http.Error(w, "Error encoding JSON in ListTopicsHandler", http.StatusInternalServerError)
		return
	}
	log.Println("Successfully retrieved list of topics from db")

}
