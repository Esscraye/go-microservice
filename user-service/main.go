package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	db *gorm.DB
)

// Define a custom type for the context key
type contextKey string

const userIDKey contextKey = "user_id"

func initDB() {
	dsn := "host=db user=user password=password dbname=microservices port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
}

func migrate() {
	db.AutoMigrate(&User{})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUsers(w)
	case "POST":
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case "GET":
		getUser(w, id)
	case "PUT":
		updateUser(w, r, id)
	case "DELETE":
		deleteUser(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifiez si l'utilisateur avec le même email existe déjà
	var existingUser User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email already in use", http.StatusConflict) // Renvoie un code 409
		return
	}

	// Créez l'utilisateur si l'email n'existe pas
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, id string) {
	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request, id string) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Model(&User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, id string) {
	if err := db.Delete(&User{}, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := http.Post("http://auth-service:8080/verify-token", "application/json", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var result struct {
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, result.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	initDB()
	migrate()

	mux := http.NewServeMux()
	mux.Handle("/users", jwtMiddleware(http.HandlerFunc(usersHandler)))
	mux.Handle("/users/", jwtMiddleware(http.HandlerFunc(userHandler)))
	mux.HandleFunc("/health", healthHandler)

	log.Println("Starting User Service on :8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
