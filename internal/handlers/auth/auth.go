package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"os"

	"github.com/CVWO/sample-go-app/internal/middlewares"
	"github.com/CVWO/sample-go-app/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	var newUser User
	//get email & pw from r
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//hash the pw
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	//create the user
	var UserID int
	err = pool.QueryRow(context.Background(), "INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3) RETURNING user_id", newUser.Username, newUser.Email, string(hashedpassword)).Scan(&UserID)
	if err != nil {
		http.Error(w, "error inserting new user to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully registered new user")
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {

	//get username & pw from r
	var body struct {
		Username string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//look up hashed_password from db
	var hashed_password string
	var userID int
	err := pool.QueryRow(context.Background(), "SELECT user_id, hashed_password from USERS where username=$1", body.Username).Scan(&userID, &hashed_password)
	if err == pgx.ErrNoRows {
		response := map[string]any{
			"success": false,
			"userId":  nil,
			"message": "invalid username",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	if err != nil {
		http.Error(w, "Failed to query database for hashed password", http.StatusInternalServerError)
		return
	}

	//compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(body.Password))
	if err != nil {
		response := map[string]any{
			"success": false,
			"userId":  nil,
			"message": "invalid password",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	//user authenticated
	//Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		//30days expiration
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	fmt.Println("token: "+tokenString, err)

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		// MaxAge:   3600 * 24 * 30,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	response := map[string]any{
		"success": true,
		"userId":  userID,
		"message": "login success",
	}
	json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}

func CurrentUserHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	//Gets current user
	var key = middlewares.CtxUserKey("user_id") //key is of type CtxUserKey
	val := r.Context().Value(key)
	userID, ok := val.(int) //type assertion from any to int
	if !ok {
		http.Error(w, "Error retrieving user id from cookie", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := pool.QueryRow(context.Background(), "SELECT user_id, username, email, hashed_password, created_at FROM users WHERE user_id=$1", userID).Scan(&user.UserID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error querying current user", http.StatusInternalServerError)
		return
	}
	//write retrieved user to w
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error Encoding to json", http.StatusInternalServerError)
		return
	}
	log.Println("success in retrieval of current user")

}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var username string
	err = pool.QueryRow(context.Background(), "DELETE FROM users WHERE user_id=$1 RETURNING username", userID).Scan(&username)
	if err != nil {
		http.Error(w, "Error deleting from users table in db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"deleted username": username,
		"message":          "User deleted successfully",
	})
}
