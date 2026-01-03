package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CVWO/sample-go-app/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	//Connect to database
	
	//load env
	err := godotenv.Load(`C:\Users\aloys\Documents\Local\CVWO assignment\code\go folder\sample-go-app\.env`)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("dbURL"))
	if err != nil {
		log.Fatalf("unable to connect to psql database: %v", err)
	}
	fmt.Println("connected to database")
	defer pool.Close()

	// _, err = pool.Exec(context.Background(), "INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3)", "benny", "benny@gmail.com", "passwordh")
	// if err != nil {
	// 	log.Fatalf("error inserting into users: %v", err)
	// }
	r := router.Setup(pool)
	fmt.Println("Listening on port 8000 at http://localhost:8000!")

	log.Fatalln(http.ListenAndServe(":8000", r))
}
